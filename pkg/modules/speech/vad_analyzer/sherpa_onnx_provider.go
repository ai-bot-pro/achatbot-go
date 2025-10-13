package vad_analyzer

import (
	"path/filepath"

	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"
	"github.com/weedge/pipeline-go/pkg/logger"

	"achatbot/pkg/consts"
	"achatbot/pkg/utils"
)

type SherpaOnnxProvider struct {
	config sherpa.VadModelConfig
	//BufferSizeInSeconds float32
	vad        *sherpa.VoiceActivityDetector
	windowSize int
	name       string
}

// https://github.com/k2-fsa/sherpa-onnx/blob/v1.12.14/sherpa-onnx/csrc/silero-vad-model-config.h
func NewDefaultSherpaOnnxSileroVadModelConfig() sherpa.SileroVadModelConfig {
	return sherpa.SileroVadModelConfig{
		Model:              filepath.Join(consts.MODELS_DIR, "silero_vad.onnx"),
		Threshold:          0.5,
		MinSilenceDuration: 0.05,  //seconds
		MinSpeechDuration:  0.025, //seconds
		// If the current segment is larger than this value, then it increases
		// the threshold to 0.9 internally. After detecting this segment,
		// it resets the threshold to its original value.
		MaxSpeechDuration: 10,
	}
}

// https://github.com/k2-fsa/sherpa-onnx/blob/v1.12.14/sherpa-onnx/csrc/ten-vad-model-config.h
// https://github.com/k2-fsa/sherpa-onnx/pull/2377
func NewDefaultSherpaOnnxTenVadModelConfig() sherpa.TenVadModelConfig {
	return sherpa.TenVadModelConfig{
		Model:              filepath.Join(consts.MODELS_DIR, "ten-vad.onnx"),
		Threshold:          0.5,
		MinSilenceDuration: 0.05, //seconds
		MinSpeechDuration:  0.05, //seconds
		// If the current segment is larger than this value, then it increases
		// the threshold to 0.9 internally. After detecting this segment,
		// it resets the threshold to its original value.
		MaxSpeechDuration: 10,
	}
}
func NewDefaultSherpaOnnxVadModelConfig(name string) sherpa.VadModelConfig {
	conf := sherpa.VadModelConfig{
		TenVad:     NewDefaultSherpaOnnxTenVadModelConfig(), // small and quick than silero
		SampleRate: consts.DefaultRate,
		NumThreads: 1,
		Provider:   "cpu",
		Debug:      0,
	}
	if name == "silero" {
		conf.SileroVad = NewDefaultSherpaOnnxSileroVadModelConfig()
	}

	return conf
}

// https://github.com/k2-fsa/sherpa-onnx/blob/v1.12.14/sherpa-onnx/csrc/vad-model-config.h
// https://github.com/k2-fsa/sherpa-onnx/blob/v1.12.14/sherpa-onnx/csrc/voice-activity-detector.cc
func NewSherpaOnnxProvider(config sherpa.VadModelConfig, bufferSizeInSeconds float32) *SherpaOnnxProvider {
	// Please download silero_vad.onnx from
	// https://github.com/k2-fsa/sherpa-onnx/releases/download/asr-models/silero_vad.onnx
	// or ten-vad.onnx from
	// https://github.com/k2-fsa/sherpa-onnx/releases/download/asr-models/ten-vad.onnx

	windowSize := 0
	name := ""
	logger.Infof("%+v", config)
	if utils.FileExists(config.SileroVad.Model) {
		logger.Info("Use SileroVad")
		// 512, 1024, 1536 samples for 16000 Hz
		config.SileroVad.WindowSize = 512
		windowSize = config.SileroVad.WindowSize
		name = "SileroVad"
	} else if utils.FileExists(config.TenVad.Model) {
		logger.Info("Use TenVad")
		// 160 or 256
		config.TenVad.WindowSize = 256
		windowSize = config.TenVad.WindowSize
		name = "TenVad"
	} else {
		logger.Errorf("Please download either silero_vad.onnx or ten-vad.onnx")
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
		name:       name,
	}
}

func (s *SherpaOnnxProvider) IsActiveSpeech(audio []byte) bool {
	samples := utils.SamplesInt16ToFloat(audio)
	//logger.Infof("AcceptWaveform len %d", len(samples))
	s.vad.AcceptWaveform(samples)
	return s.vad.IsSpeech()
}

func (s *SherpaOnnxProvider) Warmup() {
}

func (s *SherpaOnnxProvider) Name() string {
	return s.name
}

func (s *SherpaOnnxProvider) Reset() {
	s.vad.Reset()
}

func (s *SherpaOnnxProvider) Release() {
	sherpa.DeleteVoiceActivityDetector(s.vad)
}

func (s *SherpaOnnxProvider) GetSampleInfo() (int, int) {
	return s.config.SampleRate, s.windowSize
}
