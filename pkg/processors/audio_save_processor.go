package processors

import (
	"achatbot/pkg/utils"
	"fmt"
	"os"
	"time"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/logger"
	"github.com/weedge/pipeline-go/pkg/processors"

	achatbot_frames "achatbot/pkg/types/frames"
)

// AudioSaveProcessor saves AudioRawFrame to file
type AudioSaveProcessor struct {
	*processors.FrameProcessor
	saveDir      string
	prefixName   string
	passRawAudio bool
}

// NewAudioSaveProcessor creates a new AudioSaveProcessor
func NewAudioSaveProcessor(prefixName string, saveDir string, passRawAudio bool) *AudioSaveProcessor {
	// Ensure the directory exists
	err := os.MkdirAll(saveDir, os.ModePerm)
	if err != nil {
		logger.Info(fmt.Sprintf("Failed to create directory %s: %v", saveDir, err))
	}

	processor := &AudioSaveProcessor{
		FrameProcessor: processors.NewFrameProcessor("AudioSaveProcessor"),
		saveDir:        saveDir,
		prefixName:     prefixName,
		passRawAudio:   passRawAudio,
	}

	return processor
}

// ProcessFrame implements the IFrameProcessor interface
func (p *AudioSaveProcessor) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	// Call parent implementation
	p.FrameProcessor.ProcessFrame(frame, direction)
	audioFrame := utils.GetAudioRawFrame(frame)
	if audioFrame == nil {
		// For non-audio frames, just pass them through
		p.PushFrame(frame, direction)
		return
	}

	filePath, err := p.save(audioFrame)
	if err != nil {
		logger.Errorf("Failed to save audio frame: %v", err)
		// Still push the frame downstream even if saving failed
		p.PushFrame(frame, direction)
		return
	}

	if p.passRawAudio {
		p.PushFrame(frame, direction)
	} else {
		// Create a new frame with the file path
		pathFrame := achatbot_frames.NewPathAudioRawFrame(
			audioFrame.Audio, audioFrame.SampleRate, audioFrame.NumChannels, audioFrame.SampleWidth, filePath)
		p.PushFrame(pathFrame, direction)
	}
}

// save saves the audio frame to a file and returns the file path
func (p *AudioSaveProcessor) save(frame *frames.AudioRawFrame) (string, error) {
	now := time.Now()
	formattedTime := now.Format("2006-01-02_15-04-05.000000")
	if len(formattedTime) > 3 {
		formattedTime = formattedTime[:len(formattedTime)-3]
	}

	fileName := fmt.Sprintf("%s_%s.wav", p.prefixName, formattedTime)

	// Use utils.SaveAudioToFile to save the audio data as a proper WAV file
	filePath, err := utils.SaveAudioToFile(
		frame.Audio,
		fileName,
		utils.WithAudioDir(p.saveDir),
		utils.WithChannels(frame.NumChannels),
		utils.WithSampleWidth(frame.SampleWidth),
		utils.WithSampleRate(frame.SampleRate),
	)
	if err != nil {
		return "", fmt.Errorf("failed to save audio file: %v", err)
	}

	logger.Info(fmt.Sprintf("Saved frame %s to path: %s", frame.String(), filePath))
	return filePath, nil
}

// SaveAllAudioProcessor saves all audio to a single file
type SaveAllAudioProcessor struct {
	*processors.FrameProcessor
	saveDir         string
	prefixName      string
	sampleRate      int
	channels        int
	sampleWidth     int
	intervalSeconds int
	currTime        int64
	audioBytes      []byte
}

// NewSaveAllAudioProcessor creates a new SaveAllAudioProcessor
func NewSaveAllAudioProcessor(
	prefixName string,
	saveDir string,
	sampleRate int,
	channels int,
	sampleWidth int,
	intervalSeconds int,
) *SaveAllAudioProcessor {
	// Ensure the directory exists
	err := os.MkdirAll(saveDir, os.ModePerm)
	if err != nil {
		logger.Info(fmt.Sprintf("Failed to create directory %s: %v", saveDir, err))
	}

	processor := &SaveAllAudioProcessor{
		FrameProcessor:  processors.NewFrameProcessor("SaveAllAudioProcessor"),
		saveDir:         saveDir,
		prefixName:      prefixName,
		sampleRate:      sampleRate,
		channels:        channels,
		sampleWidth:     sampleWidth,
		intervalSeconds: intervalSeconds,
		audioBytes:      make([]byte, 0),
	}

	return processor
}

// ProcessFrame implements the IFrameProcessor interface
func (p *SaveAllAudioProcessor) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	// Call parent implementation
	p.FrameProcessor.ProcessFrame(frame, direction)

	if _, ok := frame.(*frames.StartFrame); ok {
		logger.Info(fmt.Sprintf("%s started", p.Name()))
		p.audioBytes = make([]byte, 0)
		p.currTime = time.Now().Unix()
	}

	audioFrame := utils.GetAudioRawFrame(frame)
	if audioFrame != nil {
		p.audioBytes = append(p.audioBytes, audioFrame.Audio...)
	}

	if p.intervalSeconds > 0 {
		currentTime := time.Now().Unix()
		if int(currentTime-p.currTime) > p.intervalSeconds {
			p.save()
			p.currTime = currentTime
		}
	}

	// Push the frame downstream
	p.PushFrame(frame, direction)

	if _, ok := frame.(*frames.EndFrame); ok {
		logger.Info(fmt.Sprintf("%s end - calling save", p.Name()))
		p.save()
		p.audioBytes = make([]byte, 0)
	}

	if _, ok := frame.(*frames.CancelFrame); ok {
		logger.Info(fmt.Sprintf("%s cancelled", p.Name()))
		p.save()
		p.audioBytes = make([]byte, 0)
	}
}

// save saves accumulated audio bytes to a file
func (p *SaveAllAudioProcessor) save() error {
	// todo: taskpool to do i/o tasks
	if len(p.audioBytes) == 0 {
		logger.Info("No audio bytes to save")
		return nil // Nothing to save
	}

	now := time.Now()
	formattedTime := now.Format("2006-01-02_15-04-05.000000")
	if len(formattedTime) > 3 {
		formattedTime = formattedTime[:len(formattedTime)-3]
	}

	fileName := fmt.Sprintf("%s_%s.wav", p.prefixName, formattedTime)

	// Use utils.SaveAudioToFile to save the audio data as a proper WAV file
	logger.Info(fmt.Sprintf("Saving %d bytes to file %s with sample rate %d, channels %d, sample width %d",
		len(p.audioBytes), fileName, p.sampleRate, p.channels, p.sampleWidth))

	filePath, err := utils.SaveAudioToFile(
		p.audioBytes,
		fileName,
		utils.WithAudioDir(p.saveDir),
		utils.WithChannels(p.channels),
		utils.WithSampleWidth(p.sampleWidth),
		utils.WithSampleRate(p.sampleRate),
	)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to save audio file: %v", err))
		return fmt.Errorf("failed to save audio file: %v", err)
	}

	logger.Info(fmt.Sprintf("Saved %d bytes to path: %s", len(p.audioBytes), filePath))
	return nil
}
