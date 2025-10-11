package processors

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/logger"
	"github.com/weedge/pipeline-go/pkg/processors"

	"achatbot/pkg/common"
	"achatbot/pkg/params"
	"achatbot/pkg/types"
	achatbot_frames "achatbot/pkg/types/frames"
	"achatbot/pkg/utils"
)

// AudioVADInputProcessor processes audio input with VAD
type AudioVADInputProcessor struct {
	*processors.AsyncFrameProcessor
	params        *params.AudioVADParams
	ctx           context.Context
	cancel        context.CancelFunc
	audioInQueue  chan *frames.AudioRawFrame
	audioTask     *sync.WaitGroup
	pushFrameTask *sync.WaitGroup
	vadAnalyzer   common.IVADAnalyzer
	vadRingBuffer *utils.RingBuffer // buffer audio byte
}

// NewAudioVADInputProcessor creates a new AudioVADInputProcessor
func NewAudioVADInputProcessor(name string, params *params.AudioVADParams) *AudioVADInputProcessor {
	ctx, cancel := context.WithCancel(context.Background())

	return &AudioVADInputProcessor{
		AsyncFrameProcessor: processors.NewAsyncFrameProcessor(name),
		params:              params,
		ctx:                 ctx,
		cancel:              cancel,
		audioInQueue:        make(chan *frames.AudioRawFrame, 1024), // buffer size
		audioTask:           &sync.WaitGroup{},
		pushFrameTask:       &sync.WaitGroup{},
		vadAnalyzer:         params.VADAnalyzer,
		vadRingBuffer: utils.NewRingBuffer(
			params.AudioInBufferSecs * params.AudioParams.AudioInSampleRate * params.AudioParams.AudioInSampleWidth,
		),
	}
}

// GetVADAnalyzer returns the VAD analyzer
func (p *AudioVADInputProcessor) GetVADAnalyzer() common.IVADAnalyzer {
	return p.vadAnalyzer
}

// SetVADAnalyzer sets the VAD analyzer
func (p *AudioVADInputProcessor) SetVADAnalyzer(analyzer common.IVADAnalyzer) {
	p.vadAnalyzer = analyzer
}

// Start starts the processor
func (p *AudioVADInputProcessor) Start(frame *frames.StartFrame) {
	if p.params.AudioInEnabled || p.params.VADEnabled {
		p.audioTask.Add(1)
		go p.audioTaskHandler()
	}

	logger.Info("AudioVADInputProcessor Start")
}

// Stop stops the processor
func (p *AudioVADInputProcessor) Stop() {
	logger.Info("AudioVADInputProcessor Stopping")
	if p.params.AudioInEnabled || p.params.VADEnabled {
		p.cancel()
		p.audioTask.Wait()
	}

	// Wait for push frame task to finish with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		p.pushFrameTask.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Task completed successfully
	case <-ctx.Done():
		// Timeout occurred
		logger.Warn("Timeout occurred while waiting for push frame task to finish")
	}
	logger.Info("AudioVADInputProcessor Stop Done")
}

// Cancel cancels all tasks
func (p *AudioVADInputProcessor) Cancel(frame *frames.CancelFrame) {
	logger.Info("AudioVADInputProcessor Cancelling")
	if p.params.AudioInEnabled || p.params.VADEnabled {
		p.cancel()
		p.audioTask.Wait()
	}

	// p.Cleanup() // pipeline.Cleanup() to call processor Cleanup
	logger.Info("AudioVADInputProcessor Cancel Done")
}

// PushAudioFrame pushes an audio frame to the queue
func (p *AudioVADInputProcessor) PushAudioFrame(frame *frames.AudioRawFrame) error {
	if p.params.AudioInEnabled || p.params.VADEnabled {
		select {
		case p.audioInQueue <- frame:
			return nil
		case <-p.ctx.Done():
			return p.ctx.Err()
		default:
			// Channel may be closed, return error instead of panicking
			return fmt.Errorf("audioInQueue is closed, cannot push audio frame")
		}
	}
	return nil
}

// audioTaskHandler handles audio processing
func (p *AudioVADInputProcessor) audioTaskHandler() {
	defer p.audioTask.Done()

	vadState := types.Quiet

	for {
		select {
		case audioFrame := <-p.audioInQueue:
			if audioFrame == nil {
				logger.Infof("%s audioFrame is nil, audioInQueue closed", p.Name())
				return
			}

			// Get audio frame to ringbuffer and get vad chunk (window_size*2 bytes) to process
			p.vadRingBuffer.PushBytes(audioFrame.Audio)
			bytesChunkSize := p.vadAnalyzer.GetWindowSize() * p.params.AudioInSampleWidth
			for p.vadRingBuffer.Size() >= bytesChunkSize {
				chunkBytes := p.vadRingBuffer.PopBytes(bytesChunkSize)

				var vadStateFrame *achatbot_frames.VADStateAudioRawFrame
				var userInterruptionFrame frames.Frame

				// Check VAD and push event if necessary
				if p.params.VADEnabled {
					vadStateFrame, userInterruptionFrame = p.handleVAD(chunkBytes, vadState)
					if vadStateFrame != nil {
						vadState = vadStateFrame.State
					}
				}

				// Handle user started speaking
				if _, ok := userInterruptionFrame.(*achatbot_frames.UserStartedSpeakingFrame); ok {
					p.handleInterruptions(userInterruptionFrame, true)
				}

				// Push audio downstream if passthrough
				if p.params.VADEnabled && p.params.VADAudioPassthrough {
					if vadStateFrame != nil && len(vadStateFrame.Audio) > 0 {
						p.PushDownstreamFrame(vadStateFrame)
					}
				} else {
					p.PushDownstreamFrame(audioFrame)
				}

				// Handle user stopped speaking
				if _, ok := userInterruptionFrame.(*achatbot_frames.UserStoppedSpeakingFrame); ok {
					p.handleInterruptions(userInterruptionFrame, true)
				}
			}
		case <-p.ctx.Done():
			logger.Info(fmt.Sprintf("%s audio_task_handler cancelled", p.Name()))
			return
		}
	}
}

// vadAnalyze analyzes audio using VAD
func (p *AudioVADInputProcessor) vadAnalyze(audioBytes []byte) *achatbot_frames.VADStateAudioRawFrame {
	//println(len(audioBytes), p.params.String())
	vadStateFrame := &achatbot_frames.VADStateAudioRawFrame{
		State:         types.Quiet,
		AudioRawFrame: frames.NewAudioRawFrame(audioBytes, p.params.AudioInSampleRate, p.params.AudioInChannels, p.params.AudioInSampleWidth),
	}

	if p.vadAnalyzer != nil {
		vadStateFrame = p.vadAnalyzer.AnalyzeAudio(audioBytes)
	}

	return vadStateFrame
}

// handleVAD handles VAD processing
func (p *AudioVADInputProcessor) handleVAD(audioBytes []byte, vadState types.VADState) (*achatbot_frames.VADStateAudioRawFrame, frames.Frame) {
	vadStateFrame := p.vadAnalyze(audioBytes)

	newVadState := vadStateFrame.State
	var userInterruptionFrame frames.Frame

	if newVadState != vadState &&
		newVadState != types.Starting &&
		newVadState != types.Stopping {
		switch newVadState {
		case types.Speaking:
			userInterruptionFrame = achatbot_frames.NewUserStartedSpeakingFrame()
		case types.Quiet:
			userInterruptionFrame = achatbot_frames.NewUserStoppedSpeakingFrame()
		}
	}

	return vadStateFrame, userInterruptionFrame
}

// startInterruption starts an interruption
func (p *AudioVADInputProcessor) startInterruption() {
	if !p.InterruptionsAllowed() {
		return
	}

	p.AsyncFrameProcessor.HandleInterruptions(frames.NewStartInterruptionFrame())
}

// stopInterruption stops an interruption
func (p *AudioVADInputProcessor) stopInterruption() {
	if !p.InterruptionsAllowed() {
		return
	}

	p.PushDownstreamFrame(frames.NewStopInterruptionFrame())
}

// handleInterruptions handles interruption frames
func (p *AudioVADInputProcessor) handleInterruptions(frame frames.Frame, pushFrame bool) {
	if p.InterruptionsAllowed() {
		switch frame.(type) {
		case *achatbot_frames.UserStartedSpeakingFrame:
			logger.Info("User started speaking")
			p.startInterruption()
		case *achatbot_frames.UserStoppedSpeakingFrame:
			logger.Info("User stopped speaking")
			p.stopInterruption()
		case *achatbot_frames.BotInterruptionFrame:
			logger.Info("Bot interruption")
			p.startInterruption()
		}
	}

	if pushFrame {
		p.PushDownstreamFrame(frame)
	}

}

// ProcessFrame processes a frame
func (p *AudioVADInputProcessor) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	// call frame processor to init star frame init
	p.AsyncFrameProcessor.WithPorcessFrameAllowPush(false).ProcessFrame(frame, direction)

	switch f := frame.(type) {
	case *frames.StartFrame:
		p.PushFrame(f, direction)
		p.Start(f)
	case *frames.EndFrame:
		p.PushFrame(f, direction)
		p.Stop()
	case *frames.CancelFrame:
		p.PushFrame(f, direction)
		p.Cancel(f)
	case *achatbot_frames.BotInterruptionFrame:
		p.handleInterruptions(f, false)
	case *frames.StartInterruptionFrame:
		p.startInterruption()
	case *frames.StopInterruptionFrame:
		p.stopInterruption()
	default:
		p.QueueFrame(f, direction)
	}

}

// Cleanup performs cleanup operations
func (p *AudioVADInputProcessor) Cleanup() {
	// Close the audio input queue
	if p.audioInQueue != nil {
		close(p.audioInQueue)
		p.audioInQueue = nil
	}
	p.AsyncFrameProcessor.Cleanup()
	logger.Info("AudioVADInputProcessor Cleanup Done")
}
