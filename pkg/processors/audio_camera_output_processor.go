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
	achatbot_frames "achatbot/pkg/types/frames"
	"achatbot/pkg/utils"
)

// AudioCameraOutputProcessor processes audio and camera output
type AudioCameraOutputProcessor struct {
	*processors.AsyncFrameProcessor
	params *params.AudioCameraParams

	ctx            context.Context
	cancel         context.CancelFunc
	audioOutQueue  chan []byte
	audioOutTask   *sync.WaitGroup
	cameraOutQueue chan frames.Frame
	cameraOutTask  *sync.WaitGroup

	// Audio accumulation buffer for 16-bit samples to write out stream device
	audioOutBuff []byte

	// Out audio chunk size in bytes
	audioChunkSize int

	// These are the images that we should send to the camera at our desired framerate
	cameraImages []*frames.ImageRawFrame

	// Indicates if the bot is currently speaking
	botSpeaking bool

	// Camera output timing
	cameraOutStartTime     time.Time
	cameraOutFrameIndex    int
	cameraOutFrameDuration time.Duration
	cameraOutFrameReset    time.Duration

	// Transport writer
	transportWriter common.ITransportWriter
}

// NewAudioCameraOutputProcessor creates a new AudioCameraOutputProcessor
func NewAudioCameraOutputProcessor(name string, params *params.AudioCameraParams) *AudioCameraOutputProcessor {
	ctx, cancel := context.WithCancel(context.Background())

	p := &AudioCameraOutputProcessor{
		AsyncFrameProcessor: processors.NewAsyncFrameProcessor(name),
		params:              params,
		ctx:                 ctx,
		cancel:              cancel,
		audioOutQueue:       make(chan []byte, 100), // buffer size
		audioOutTask:        &sync.WaitGroup{},
		cameraOutQueue:      make(chan frames.Frame, 100), // buffer size
		cameraOutTask:       &sync.WaitGroup{},
		audioOutBuff:        make([]byte, 0),
		cameraImages:        make([]*frames.ImageRawFrame, 0),
		botSpeaking:         false,
		transportWriter:     params.TransportWriter,
	}

	p.audioChunkSize = (p.params.AudioOutSampleRate / 100) * p.params.AudioOutChannels * 2 * p.params.AudioOut10msChunks
	if p.audioChunkSize == 0 {
		logger.Warn("audioChunkSize is 0, please check your audioOutSampleRate and audioOutChannels and audioOut10msChunks values in your config")
		return nil
	}

	return p
}

// Start starts the processor
func (p *AudioCameraOutputProcessor) Start(frame *frames.StartFrame) {
	// Create media threads queues and task
	if p.params.CameraOutEnabled {
		p.cameraOutTask.Add(1)
		go p.cameraOutTaskHandler()
	}
	if p.params.AudioOutEnabled {
		p.audioOutTask.Add(1)
		go p.audioOutTaskHandler()
	}
	logger.Info("AudioCameraOutputProcessor Start")
}

// Stop stops the processor
func (p *AudioCameraOutputProcessor) Stop(frame *frames.EndFrame) {
	logger.Info("AudioCameraOutputProcessor Stopping")

	// Cancel tasks
	p.cancel()

	// Wait for audio output task to finish
	if p.params.AudioOutEnabled {
		p.audioOutTask.Wait()
	}

	// Wait for camera output task to finish
	if p.params.CameraOutEnabled {
		p.cameraOutTask.Wait()
	}

	logger.Info("AudioCameraOutputProcessor Stop Done")
}

// Cancel cancels the processor
func (p *AudioCameraOutputProcessor) Cancel(frame *frames.CancelFrame) {
	logger.Info("AudioCameraOutputProcessor Cancelling")

	// Cancel tasks
	p.cancel()

	// Wait for audio output task to finish
	if p.params.AudioOutEnabled {
		p.audioOutTask.Wait()
	}

	// Wait for camera output task to finish
	if p.params.CameraOutEnabled {
		p.cameraOutTask.Wait()
	}

	logger.Info("AudioCameraOutputProcessor Cancel Done")
}

// ProcessFrame processes a frame
func (p *AudioCameraOutputProcessor) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	// call frame processor to init star frame init
	p.AsyncFrameProcessor.WithPorcessFrameAllowPush(false).ProcessFrame(frame, direction)

	switch f := frame.(type) {
	case *frames.StartFrame:
		p.PushFrame(f, direction)
		p.Start(f)
	case *frames.EndFrame:
		p.PushFrame(f, direction)
		p.Stop(f)
	case *frames.CancelFrame:
		p.PushFrame(f, direction)
		p.Cancel(f)
	case *frames.StartInterruptionFrame:
		p.PushFrame(f, direction)
		p.handleInterruptions(frame)
	case *frames.StopInterruptionFrame:
		p.PushFrame(f, direction)
	case *achatbot_frames.TransportMessageFrame, *frames.TextFrame, *achatbot_frames.AnimationAudioRawFrame:
		err := p.transportWriter.WriteFrame(f)
		if err != nil {
			logger.Error(fmt.Sprintf("Error Write %T", f), "error", err)
		}
	case *frames.AudioRawFrame:
		p.handleAudio(f)
	case *frames.ImageRawFrame:
		p.handleImage(f)
	case *achatbot_frames.SpriteFrame:
		p.handleSpriteImages(f)
	case *achatbot_frames.TTSStartedFrame:
		p.botStartedSpeaking()
		p.QueueFrame(frame, processors.FrameDirectionDownstream)
	case *achatbot_frames.TTSStoppedFrame:
		p.botStoppedSpeaking()
		p.QueueFrame(frame, processors.FrameDirectionDownstream)
	default:
		p.QueueFrame(f, direction)
	}
}

// handleInterruptions handles interruption frames
func (p *AudioCameraOutputProcessor) handleInterruptions(frame frames.Frame) {
	// Call parent implementation
	p.AsyncFrameProcessor.HandleInterruptions(frame)

	// Handle start interruption
	if _, ok := frame.(*frames.StartInterruptionFrame); ok {
		// Let's send a bot stopped speaking if we have to
		if p.botSpeaking {
			p.botStoppedSpeaking()
		}
	}
}

// botStartedSpeaking handles bot started speaking
func (p *AudioCameraOutputProcessor) botStartedSpeaking() {
	logger.Debug("Bot started speaking")
	p.botSpeaking = true
	p.PushFrame(achatbot_frames.NewBotStartedSpeakingFrame(), processors.FrameDirectionUpstream)
}

// botStoppedSpeaking handles bot stopped speaking
func (p *AudioCameraOutputProcessor) botStoppedSpeaking() {
	logger.Debug("Bot stopped speaking")
	p.botSpeaking = false
	p.PushFrame(achatbot_frames.NewBotStoppedSpeakingFrame(), processors.FrameDirectionUpstream)
}

// handleAudio handles audio frames
func (p *AudioCameraOutputProcessor) handleAudio(frame *frames.AudioRawFrame) {
	if !p.params.AudioOutEnabled {
		return
	}

	audio := frame.Audio
	for i := 0; i < len(audio); i += p.audioChunkSize {
		end := min(i+p.audioChunkSize, len(audio))
		chunk := audio[i:end]

		// Add chunk to queue
		select {
		case p.audioOutQueue <- chunk:
		case <-p.ctx.Done():
			return
		}

		// Push bot speaking frame upstream if bot is speaking,
		// quickly push, this can test upstream processor process frame speed :)
		p.PushFrame(achatbot_frames.NewBotSpeakingFrame(), processors.FrameDirectionUpstream)
	}
}

// audioOutTaskHandler handles audio output task
func (p *AudioCameraOutputProcessor) audioOutTaskHandler() {
	defer p.audioOutTask.Done()

	for {
		select {
		case chunk := <-p.audioOutQueue:
			err := p.transportWriter.WriteRawAudio(chunk)
			if err != nil {
				logger.Error(fmt.Sprintf("%s audio_out_task_handler error", p.Name()), "error", err)
			}
		case <-p.ctx.Done():
			logger.Info(fmt.Sprintf("%s audio_out_task_handler cancelled", p.Name()))
			return
		case <-time.After(1 * time.Second):
			// Timeout, continue the loop
			continue
		}
	}
}

// handleImage handles image frames
func (p *AudioCameraOutputProcessor) handleImage(frame *frames.ImageRawFrame) {
	if !p.params.CameraOutEnabled {
		return
	}

	if p.params.CameraOutIsLive {
		// NOTE: out processor is last start to init camera out queue,
		// if pipeline other processor is start slow,
		// and push frame before out processor init, out processor will lost frame
		select {
		case p.cameraOutQueue <- frame:
		case <-p.ctx.Done():
			return
		}
	} else {
		p.setCameraImages([]*frames.ImageRawFrame{frame})
	}
}

// handleSpriteImages handles sprite image frames
func (p *AudioCameraOutputProcessor) handleSpriteImages(frame *achatbot_frames.SpriteFrame) {
	if !p.params.CameraOutEnabled {
		return
	}
	p.setCameraImages(frame.Images)
}

// setCameraImages sets camera images
func (p *AudioCameraOutputProcessor) setCameraImages(images []*frames.ImageRawFrame) {
	p.cameraImages = images
}

// drawImage draws an image
func (p *AudioCameraOutputProcessor) drawImage(frame *frames.ImageRawFrame) error {
	// Resize if needed
	if frame.Size.Width != p.params.CameraOutWidth || frame.Size.Height != p.params.CameraOutHeight {
		// Convert byte data to image
		imageObj := utils.ImageFromBytes(frame.Image, frame.Size.Width, frame.Size.Height, frame.Mode)

		// Resize image
		dstImageObj := utils.ResizeImage(imageObj, p.params.CameraOutWidth, p.params.CameraOutHeight)

		// Change image frame
		imageInfo := utils.GetImageInfo(dstImageObj, frame.Format)
		frame.Image = imageInfo.Bytes
		frame.Mode = imageInfo.Mode
		frame.Size = frames.ImageSize{
			Width:  imageInfo.Width,
			Height: imageInfo.Height,
		}

		logger.Warn(fmt.Sprintf("%v does not have the expected width: %d and height: %d, resizing would be needed", frame, p.params.CameraOutWidth, p.params.CameraOutHeight))
	}

	err := p.transportWriter.WriteFrame(frame)
	return err
}

// cameraOutIsLiveHandler handles live camera output
func (p *AudioCameraOutputProcessor) cameraOutIsLiveHandler() error {
	// Get image from queue
	var imageFrame frames.Frame
	select {
	case imageFrame = <-p.cameraOutQueue:
	case <-p.ctx.Done():
		return p.ctx.Err()
	}

	// We get the start time as soon as we get the first image
	if p.cameraOutStartTime.IsZero() {
		p.cameraOutStartTime = time.Now()
		p.cameraOutFrameIndex = 0
	}

	// Calculate how much time we need to wait before rendering next image
	realElapsedTime := time.Since(p.cameraOutStartTime)
	realRenderTime := time.Duration(p.cameraOutFrameIndex) * p.cameraOutFrameDuration
	delayTime := p.cameraOutFrameDuration + realRenderTime - realElapsedTime

	if utils.Abs(delayTime) > p.cameraOutFrameReset {
		p.cameraOutStartTime = time.Now()
		p.cameraOutFrameIndex = 0
	} else if delayTime > 0 {
		time.Sleep(delayTime)
		p.cameraOutFrameIndex++
	}

	// Render image
	if imgFrame, ok := imageFrame.(*frames.ImageRawFrame); ok {
		err := p.drawImage(imgFrame)
		return err
	}

	return nil
}

// cameraOutTaskHandler handles camera output task
func (p *AudioCameraOutputProcessor) cameraOutTaskHandler() {
	defer p.cameraOutTask.Done()

	p.cameraOutFrameDuration = time.Second / time.Duration(p.params.CameraOutFramerate)
	p.cameraOutFrameReset = p.cameraOutFrameDuration * 5

	for {
		select {
		case <-p.ctx.Done():
			logger.Info(fmt.Sprintf("%s camera_out_task_handler cancelled", p.Name()))
			return
		default:
			var err error
			if p.params.CameraOutIsLive {
				err = p.cameraOutIsLiveHandler()
			} else if len(p.cameraImages) > 0 {
				// Cycle through images
				index := p.cameraOutFrameIndex % len(p.cameraImages)
				err = p.drawImage(p.cameraImages[index])
				p.cameraOutFrameIndex++
				time.Sleep(p.cameraOutFrameDuration)
			} else {
				time.Sleep(p.cameraOutFrameDuration)
			}

			if err != nil {
				logger.Error("Error writing to camera", "error", err)
			}
		}
	}
}

// Cleanup performs cleanup operations
func (p *AudioCameraOutputProcessor) Cleanup() {
	// Close the audio output queue
	if p.audioOutQueue != nil {
		close(p.audioOutQueue)
		p.audioOutQueue = nil
	}
	// Close the camera output queue
	if p.cameraOutQueue != nil {
		close(p.cameraOutQueue)
		p.cameraOutQueue = nil
	}
	p.AsyncFrameProcessor.Cleanup()
	logger.Info("AudioCameraOutputProcessor Cleanup Done")
}
