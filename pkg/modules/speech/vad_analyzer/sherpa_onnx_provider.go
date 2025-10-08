package vad_analyzer

import (
	"fmt"
	"log/slog"

	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"

	"achatbot/pkg/utils"
)

type SherpaOnnxProvider struct {
	config sherpa.VadModelConfig
	//BufferSizeInSeconds float32
	vad        *sherpa.VoiceActivityDetector
	windowSize int
}

func NewSherpaOnnxProvider(config sherpa.VadModelConfig, bufferSizeInSeconds float32) *SherpaOnnxProvider {
	// Please download silero_vad.onnx from
	// https://github.com/k2-fsa/sherpa-onnx/releases/download/asr-models/silero_vad.onnx
	// or ten-vad.onnx from
	// https://github.com/k2-fsa/sherpa-onnx/releases/download/asr-models/ten-vad.onnx

	windowSize := 0
	if utils.FileExists(config.SileroVad.Model) {
		slog.Info("Use ten-vad")
		//config.SileroVad.Threshold = 0.5
		//config.SileroVad.MinSilenceDuration = 0.5
		//config.SileroVad.MinSpeechDuration = 0.25
		//config.SileroVad.MaxSpeechDuration = 10
		config.SileroVad.WindowSize = 512
		windowSize = config.SileroVad.WindowSize
	} else if utils.FileExists(config.TenVad.Model) {
		slog.Info("Use ten-vad")
		//config.TenVad.Threshold = 0.5
		//config.TenVad.MinSilenceDuration = 0.5
		//config.TenVad.MinSpeechDuration = 0.25
		//config.TenVad.MaxSpeechDuration = 10
		config.TenVad.WindowSize = 256
		windowSize = config.TenVad.WindowSize
	} else {
		slog.Error("Please download either ./models/silero_vad.onnx or ./models/ten-vad.onnx")
		return nil
	}

	config.SampleRate = 16000
	config.NumThreads = 1
	//config.Provider = "cpu"
	//config.Debug = 1

	vad := sherpa.NewVoiceActivityDetector(&config, float32(bufferSizeInSeconds))

	return &SherpaOnnxProvider{
		config:     config,
		vad:        vad,
		windowSize: windowSize,
	}
}

func (s *SherpaOnnxProvider) IsActiveSpeech(audio []byte) bool {
	samples := utils.SamplesInt16ToFloat(audio)
	slog.Debug(fmt.Sprintf("AcceptWaveform len %d", len(samples)))
	s.vad.AcceptWaveform(samples)
	return s.vad.IsSpeech()
}

func (s *SherpaOnnxProvider) Release() {
	sherpa.DeleteVoiceActivityDetector(s.vad)
}

func (s *SherpaOnnxProvider) GetSampleInfo() (int, int) {
	return s.config.SampleRate, s.windowSize
}
