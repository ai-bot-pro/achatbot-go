package asr

import (
	"achatbot/pkg/consts"
	"path/filepath"

	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"
	"github.com/weedge/pipeline-go/pkg/logger"

	"achatbot/pkg/utils"
)

type SherpaOnnxProvider struct {
	config     sherpa.OfflineRecognizerConfig
	recognizer *sherpa.OfflineRecognizer
	name       string
	sampleRate int
}

// https://github.com/k2-fsa/sherpa-onnx/blob/v1.12.14/sherpa-onnx/csrc/offline-paraformer-model-config.h
// https://github.com/modelscope/FunASR
// https://arxiv.org/abs/2206.08317
func NewDefaultSherpaOnnxOfflineParaformerModelConfig() (sherpa.OfflineParaformerModelConfig, string) {
	return sherpa.OfflineParaformerModelConfig{
		Model: filepath.Join(consts.MODELS_DIR, "csukuangfj/sherpa-onnx-paraformer-zh-2023-09-14/model.int8.onnx"),
	}, filepath.Join(consts.MODELS_DIR, "csukuangfj/sherpa-onnx-paraformer-zh-2023-09-14/tokens.txt")
}

// https://github.com/k2-fsa/sherpa-onnx/blob/v1.12.14/sherpa-onnx/csrc/offline-whisper-model-config.h
// https://github.com/openai/whisper
// https://arxiv.org/abs/2212.04356
func NewDefaultSherpaOnnxOfflineWhisperModelConfig() (sherpa.OfflineWhisperModelConfig, string) {
	return sherpa.OfflineWhisperModelConfig{
		Encoder: filepath.Join(consts.MODELS_DIR, "csukuangfj/sherpa-onnx-whisper-tiny.en/tiny-encoder.onnx"),
		Decoder: filepath.Join(consts.MODELS_DIR, "csukuangfj/sherpa-onnx-whisper-tiny.en/tiny-decoder.onnx"),
		// Available languages can be found at
		// https://github.com/openai/whisper/blob/main/whisper/tokenizer.py#L10
		//
		// Note: For non-multilingual models, it supports only "en"
		//
		// If empty, we will infer it from the input audio file when
		// the model is multilingual.
		Language: "en",
		// Valid values are transcribe and translate
		//
		// Note: For non-multilingual models, it supports only "transcribe"
		Task: "transcribe",
		// Number of tail padding frames.
		//
		// Since we remove the 30-second constraint, we need to add some paddings
		// at the end.
		//
		// Recommended values:
		//   - 50 for English models
		//   - 300 for multilingual models
		TailPaddings: -1,
	}, filepath.Join(consts.MODELS_DIR, "csukuangfj/sherpa-onnx-whisper-tiny.en/tokens.txt")
}

// https://github.com/k2-fsa/sherpa-onnx/blob/v1.12.14/sherpa-onnx/csrc/offline-zipformer-ctc-model-config.h
// https://github.com/k2-fsa/icefall
// https://arxiv.org/abs/2310.11230
func NewDefaultSherpaOnnxOfflineZipformerCtcModelConfig() (sherpa.OfflineZipformerCtcModelConfig, string) {
	return sherpa.OfflineZipformerCtcModelConfig{
		Model: filepath.Join(consts.MODELS_DIR, "csukuangfj/sherpa-onnx-zipformer-ctc-zh-int8-2025-07-03/model.int8.onnx"),
	}, filepath.Join(consts.MODELS_DIR, "csukuangfj/sherpa-onnx-zipformer-ctc-zh-int8-2025-07-03/tokens.txt")
}

// https://github.com/k2-fsa/sherpa-onnx/blob/v1.12.14/sherpa-onnx/csrc/offline-sense-voice-model-config.h
// https://github.com/FunAudioLLM/SenseVoice
// https://arxiv.org/abs/2407.04051
func NewDefaultSherpaOnnxOfflineSenseVoiceModelConfig() (sherpa.OfflineSenseVoiceModelConfig, string) {
	return sherpa.OfflineSenseVoiceModelConfig{
		Model:                       filepath.Join(consts.MODELS_DIR, "csukuangfj/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/model.onnx"),
		Language:                    "", // If not empty, specify the Language for the input wave
		UseInverseTextNormalization: 1,  // 1 to use inverse text normalization
	}, filepath.Join(consts.MODELS_DIR, "csukuangfj/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/tokens.txt")

}

// https://github.com/k2-fsa/sherpa-onnx/blob/v1.12.14/sherpa-onnx/csrc/offline-moonshine-model-config.h
// https://github.com/moonshine-ai/moonshine
// https://arxiv.org/abs/2410.15608
// https://arxiv.org/abs/2509.02523
func NewDefaultSherpaOnnxOfflineMoonshineModelConfig() (sherpa.OfflineMoonshineModelConfig, string) {
	return sherpa.OfflineMoonshineModelConfig{
		Preprocessor:    filepath.Join(consts.MODELS_DIR, "csukuangfj/sherpa-onnx-moonshine-tiny-en-int8/preprocess.onnx"),
		Encoder:         filepath.Join(consts.MODELS_DIR, "csukuangfj/sherpa-onnx-moonshine-tiny-en-int8/encoder.int8.onnx"),
		UncachedDecoder: filepath.Join(consts.MODELS_DIR, "csukuangfj/sherpa-onnx-moonshine-tiny-en-int8/uncached_decoder.int8.onnx"),
		CachedDecoder:   filepath.Join(consts.MODELS_DIR, "csukuangfj/sherpa-onnx-moonshine-tiny-en-int8/cached_decoder.int8.onnx"),
	}, filepath.Join(consts.MODELS_DIR, "csukuangfj/sherpa-onnx-moonshine-tiny-en-int8/tokens.txt")
}

// https://github.com/k2-fsa/sherpa-onnx/blob/v1.12.14/sherpa-onnx/csrc/offline-fire-red-asr-model-config.h
// https://github.com/FireRedTeam/FireRedASR
// https://arxiv.org/abs/2501.14350
func NewDefaultSherpaOnnxOfflineFireRedAsrModelConfig() (sherpa.OfflineFireRedAsrModelConfig, string) {
	return sherpa.OfflineFireRedAsrModelConfig{
		Encoder: filepath.Join(consts.MODELS_DIR, "csukuangfj/sherpa-onnx-fire-red-asr-large-zh_en-2025-02-16/encoder.int8.onnx"),
		Decoder: filepath.Join(consts.MODELS_DIR, "csukuangfj/sherpa-onnx-fire-red-asr-large-zh_en-2025-02-16/decoder.int8.onnx"),
	}, filepath.Join(consts.MODELS_DIR, "csukuangfj/sherpa-onnx-fire-red-asr-large-zh_en-2025-02-16/tokens.txt")
}

// https://github.com/k2-fsa/sherpa-onnx/blob/v1.12.14/sherpa-onnx/csrc/offline-dolphin-model-config.h
// https://github.com/DataoceanAI/Dolphin
// https://arxiv.org/abs/2503.20212
func NewDefaultSherpaOnnxOfflineDolphinModelConfig() (sherpa.OfflineDolphinModelConfig, string) {
	return sherpa.OfflineDolphinModelConfig{
		Model: filepath.Join(consts.MODELS_DIR, "csukuangfj/sherpa-onnx-dolphin-base-ctc-multi-lang-int8-2025-04-02/model.int8.onnx"),
	}, filepath.Join(consts.MODELS_DIR, "csukuangfj/sherpa-onnx-dolphin-base-ctc-multi-lang-int8-2025-04-02/tokens.txt")
}

// https://github.com/k2-fsa/sherpa-onnx/blob/v1.12.14/sherpa-onnx/csrc/offline-nemo-enc-dec-ctc-model-config.h
// https://k2-fsa.github.io/sherpa/onnx/pretrained_models/offline-ctc/index.html
// https://catalog.ngc.nvidia.com/orgs/nvidia/collections/nemo_asr
// https://github.com/NVIDIA-NeMo/NeMo
func NewDefaultSherpaOnnxOfflineNemoEncDecCtcModelConfig() (sherpa.OfflineNemoEncDecCtcModelConfig, string) {
	return sherpa.OfflineNemoEncDecCtcModelConfig{
		Model: filepath.Join(consts.MODELS_DIR, "csukuangfj/sherpa-onnx-nemo-ctc-en-conformer-medium/model.onnx"),
	}, filepath.Join(consts.MODELS_DIR, "csukuangfj/sherpa-onnx-nemo-ctc-en-conformer-medium/tokens.txt")
}

func NewDefaultSherpaOnnxOfflineRecognizerConfig() sherpa.OfflineRecognizerConfig {
	asrConf, tokenPath := NewDefaultSherpaOnnxOfflineSenseVoiceModelConfig()
	conf := sherpa.OfflineRecognizerConfig{
		FeatConfig: sherpa.FeatureConfig{SampleRate: consts.DefaultRate, FeatureDim: 80},
		ModelConfig: sherpa.OfflineModelConfig{
			SenseVoice: asrConf,
			Tokens:     tokenPath,
			NumThreads: 1, Debug: 0, Provider: "cpu", ModelingUnit: "cjkchar",
		},
		DecodingMethod: "greedy_search", // greedy_search, modified_beam_search
		MaxActivePaths: 4,               // only valid when decoding_method is modified_beam_search
	}
	return conf
}

func NewSherpaOnnxProvider(config sherpa.OfflineRecognizerConfig) *SherpaOnnxProvider {
	provider := &SherpaOnnxProvider{
		config: config,
	}
	provider.recognizer = sherpa.NewOfflineRecognizer(&config)
	if provider.recognizer == nil {
		logger.Error("Fail to create ASR")
		return nil
	}

	provider.sampleRate = config.FeatConfig.SampleRate

	logger.Info("ASR NewSherpaOnnxProvider Done", "name", provider.name)

	return provider
}

func (p *SherpaOnnxProvider) Transcribe(data []byte) string {
	samples := utils.SamplesInt16ToFloat(data)
	stream := sherpa.NewOfflineStream(p.recognizer)
	stream.AcceptWaveform(p.sampleRate, samples)
	p.recognizer.Decode(stream)
	result := stream.GetResult()
	sherpa.DeleteOfflineStream(stream)
	return result.Text
}

func (p *SherpaOnnxProvider) Release() {
	sherpa.DeleteOfflineRecognizer(p.recognizer)
}

func (p *SherpaOnnxProvider) Name() string {
	return p.name
}
