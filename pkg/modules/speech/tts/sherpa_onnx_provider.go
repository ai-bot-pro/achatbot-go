package tts

import (
	"achatbot/pkg/consts"
	"achatbot/pkg/utils"
	"math"
	"path/filepath"
	"strings"

	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"
	"github.com/weedge/pipeline-go/pkg/logger"
)

type SherpaOnnxProvider struct {
	config     sherpa.OfflineTtsConfig
	tts        *sherpa.OfflineTts
	name       string
	sid        int     // Speaker ID (multi-speaker models only)
	speed      float32 // Speech speed. larger->faster; smaller->slower
	sampleRate int
}

// https://github.com/k2-fsa/sherpa-onnx/blob/v1.12.14/sherpa-onnx/csrc/offline-tts-kokoro-model-config.h
// https://huggingface.co/hexgrad/Kokoro-82M
func NewDefaultSherpaOnnxOfflineTtsKokoroModelConfig() sherpa.OfflineTtsKokoroModelConfig {
	return sherpa.OfflineTtsKokoroModelConfig{
		Model:   filepath.Join(consts.MODELS_DIR, "csukuangfj/kokoro-multi-lang-v1_0/model.onnx"),
		Voices:  filepath.Join(consts.MODELS_DIR, "csukuangfj/kokoro-multi-lang-v1_0/voices.bin"),
		Tokens:  filepath.Join(consts.MODELS_DIR, "csukuangfj/kokoro-multi-lang-v1_0/tokens.txt"),
		DataDir: filepath.Join(consts.MODELS_DIR, "csukuangfj/kokoro-multi-lang-v1_0/espeak-ng-data"),
		DictDir: filepath.Join(consts.MODELS_DIR, "csukuangfj/kokoro-multi-lang-v1_0/dict"),
		Lexicon: strings.Join([]string{
			filepath.Join(consts.MODELS_DIR, "csukuangfj/kokoro-multi-lang-v1_0/lexicon-us-en.txt"),
			filepath.Join(consts.MODELS_DIR, "csukuangfj/kokoro-multi-lang-v1_0/lexicon-zh.txt"),
		}, ","),
	}
}

// https://k2-fsa.github.io/sherpa/onnx/tts/index.html
func NewDefaultSherpaOnnxOfflineTtsConfig() sherpa.OfflineTtsConfig {
	ttsConf := NewDefaultSherpaOnnxOfflineTtsKokoroModelConfig()
	conf := sherpa.OfflineTtsConfig{
		Model: sherpa.OfflineTtsModelConfig{
			Kokoro:     ttsConf,
			Provider:   "cpu",
			NumThreads: 1,
			Debug:      0,
		},
		RuleFsts:        "", //Path to rule.fst
		RuleFars:        "", //Path to rule.far
		MaxNumSentences: 1,  //Batch size (split long text to avoid OOM)
	}
	return conf
}

const (
	//kokoroTTS: 45->zf_xiaobei, 46->zf_xiaoni, 47->zf_xiaoxiao, 48->zf_xiaoyi, 49->zm_yunjian, 50->zm_yunxi, 51->zm_yunxia, 52->zm_yunyang
	KokoroTTS_Speaker_ZF_XiaoBei  = 45
	KokoroTTS_Speaker_ZF_XiaoNi   = 46
	KokoroTTS_Speaker_ZF_XiaoXiao = 47
	KokoroTTS_Speaker_ZF_XiaoYi   = 48
	KokoroTTS_Speaker_ZM_YunJian  = 49
	KokoroTTS_Speaker_ZM_YunXi    = 50
	KokoroTTS_Speaker_ZM_YunYang  = 52
)

func NewSherpaOnnxProvider(config sherpa.OfflineTtsConfig, sid int, speed float32, name string) *SherpaOnnxProvider {
	provider := &SherpaOnnxProvider{
		config: config,
	}
	provider.tts = sherpa.NewOfflineTts(&config)
	if provider.tts == nil {
		logger.Error("Failed to create TTS")
		return nil
	}

	if sid < 0 {
		sid = 0
	}
	provider.sid = sid
	provider.speed = speed
	provider.name = name

	provider.Warmup() // warm up and get sample rate
	logger.Info("TTS NewSherpaOnnxProvider Done", "name", name, "sid", sid, "speed", speed)

	return provider
}

func (p *SherpaOnnxProvider) Synthesize(text string) []byte {
	generateAudio := p.tts.Generate(text, p.sid, float32(math.Max(float64(p.speed), 1e-6)))
	p.sampleRate = generateAudio.SampleRate
	return utils.SamplesFloatToInt16(generateAudio.Samples)
}

func (p *SherpaOnnxProvider) Warmup() {
	generateAudio := p.tts.Generate("hi", p.sid, float32(math.Max(float64(p.speed), 1e-6)))
	p.sampleRate = generateAudio.SampleRate
}

func (p *SherpaOnnxProvider) GetSampleInfo() (int, int, int) {
	return p.sampleRate, consts.DefaultChannels, consts.DefaultSampleWidth
}

func (p *SherpaOnnxProvider) SetPromptAudio(string, []byte) error {
	return nil
}

func (p *SherpaOnnxProvider) Release() {
	sherpa.DeleteOfflineTts(p.tts)
}

func (p *SherpaOnnxProvider) Name() string {
	return p.name
}
