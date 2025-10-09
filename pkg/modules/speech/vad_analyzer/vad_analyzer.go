package vad_analyzer

import (
	"math"

	pipelineframes "github.com/weedge/pipeline-go/pkg/frames"

	"achatbot/pkg/common"
	"achatbot/pkg/types"
	localframes "achatbot/pkg/types/frames"
)

// VADAnalyzerArgs VAD分析器参数
type VADAnalyzerArgs struct {
	SampleRate  int     `json:"sample_rate"`
	NumChannels int     `json:"num_channels"`
	SampleWidth int     `json:"sample_width"`
	Confidence  float64 `json:"confidence"`
	MinVolume   float64 `json:"min_volume"`
	StartSecs   float64 `json:"start_secs"`
	StopSecs    float64 `json:"stop_secs"`
}

// VADAnalyzer 基础VAD分析器
type VADAnalyzer struct {
	args                     *VADAnalyzerArgs
	vadFrames                int
	vadFramesNumBytes        int
	sampleNumBytes           int
	vadStartFrames           int
	vadStopFrames            int
	vadStartingCount         int
	vadStoppingCount         int
	vadState                 types.VADState
	vadBuffer                []byte
	smoothingFactor          float64
	prevVolume               float64
	speechID                 int
	accumulateSpeechBytesLen int
	isFinal                  bool
	startAtS                 float64
	curAtS                   float64
	endAtS                   float64

	// 语音置信度提供者
	IVoiceConfidenceProvider common.IVoiceConfidenceProvider
}

// NewVADAnalyzer 创建新的VAD分析器
func NewVADAnalyzer(args *VADAnalyzerArgs, vcp common.IVoiceConfidenceProvider) *VADAnalyzer {
	analyzer := &VADAnalyzer{
		args:                     args,
		vadBuffer:                make([]byte, 0),
		smoothingFactor:          0.2,
		prevVolume:               0,
		vadState:                 types.Quiet,
		IVoiceConfidenceProvider: vcp,
	}

	analyzer.vadFrames = analyzer.numFramesRequired()
	analyzer.vadFramesNumBytes = analyzer.vadFrames * analyzer.args.NumChannels * analyzer.args.SampleWidth
	analyzer.sampleNumBytes = analyzer.args.SampleRate * analyzer.args.NumChannels * analyzer.args.SampleWidth

	vadFramesPerSec := float64(analyzer.vadFrames) / float64(analyzer.args.SampleRate)
	analyzer.vadStartFrames = int(math.Round(analyzer.args.StartSecs / vadFramesPerSec))
	analyzer.vadStopFrames = int(math.Round(analyzer.args.StopSecs / vadFramesPerSec))

	analyzer.Reset()
	return analyzer
}

// 创建VAD分析器的便捷函数
func NewSherpaOnnxVADAnalyzer(args *VADAnalyzerArgs) *VADAnalyzer {
	provider := &SherpaOnnxProvider{}
	return NewVADAnalyzer(args, provider)
}

// Reset 重置分析器状态
func (b *VADAnalyzer) Reset() {
	b.vadStartingCount = 0
	b.vadStoppingCount = 0
	b.vadState = types.Quiet
	b.speechID = 0
	b.accumulateSpeechBytesLen = 0
	b.isFinal = false
	b.startAtS = 0.0
	b.curAtS = 0.0
	b.endAtS = 0.0
	b.vadBuffer = make([]byte, 0)
}

func (b *VADAnalyzer) GetSampleRate() int {
	sr, _ := b.IVoiceConfidenceProvider.GetSampleInfo()
	return sr
}

func (b *VADAnalyzer) Release() {
	b.IVoiceConfidenceProvider.Release()
}

// numFramesRequired 计算需要的帧数
func (b *VADAnalyzer) numFramesRequired() int {
	_, windowSize := b.IVoiceConfidenceProvider.GetSampleInfo()
	return windowSize
}

// isActiveSpeech 判断是否是活跃的语音
func (b *VADAnalyzer) isActiveSpeech(audio []byte) bool {
	return b.IVoiceConfidenceProvider.IsActiveSpeech(audio)
}

// AnalyzeAudio 分析音频
func (b *VADAnalyzer) AnalyzeAudio(buffer []byte) *localframes.VADStateAudioRawFrame {
	// 追加缓冲区
	b.vadBuffer = append(b.vadBuffer, buffer...)
	b.curAtS = math.Round(float64(b.accumulateSpeechBytesLen)/float64(b.sampleNumBytes)*1000) / 1000
	b.accumulateSpeechBytesLen += len(buffer)

	numRequiredBytes := b.vadFramesNumBytes
	if len(b.vadBuffer) < numRequiredBytes {
		// 返回当前状态但不处理音频数据
		return &localframes.VADStateAudioRawFrame{
			AudioRawFrame: &pipelineframes.AudioRawFrame{
				DataFrame:   pipelineframes.NewDataFrameWithName("VADStateAudioRawFrame"),
				Audio:       make([]byte, 0),
				SampleRate:  b.args.SampleRate,
				NumChannels: b.args.NumChannels,
				SampleWidth: b.args.SampleWidth,
			},
			State:    b.vadState,
			SpeechID: b.speechID,
			IsFinal:  b.isFinal,
			StartAtS: b.startAtS,
			CurAtS:   b.curAtS,
			EndAtS:   b.endAtS,
		}
	}

	audioBytes := make([]byte, numRequiredBytes)
	copy(audioBytes, b.vadBuffer[:numRequiredBytes])
	b.vadBuffer = b.vadBuffer[numRequiredBytes:]

	speaking := b.isActiveSpeech(audioBytes)
	if speaking {
		switch b.vadState {
		case types.Quiet:
			b.vadState = types.Starting
			b.vadStartingCount = 1
		case types.Starting:
			b.vadStartingCount++
		case types.Stopping:
			b.vadState = types.Speaking
			b.vadStoppingCount = 0
		}
	} else {
		switch b.vadState {
		case types.Starting:
			b.vadState = types.Quiet
			b.vadStartingCount = 0
			b.Reset()
		case types.Speaking:
			b.vadState = types.Stopping
			b.vadStoppingCount = 1
		case types.Stopping:
			b.vadStoppingCount++
		}
	}

	// 检查是否开始说话
	if b.vadState == types.Starting && b.vadStartingCount >= b.vadStartFrames {
		b.vadState = types.Speaking
		b.vadStartingCount = 0

		b.speechID++
		b.isFinal = false
		b.startAtS = math.Round(float64(b.accumulateSpeechBytesLen-len(buffer))/float64(b.sampleNumBytes)*1000) / 1000
		b.endAtS = 0.0
	}

	// 检查是否结束说话
	if b.vadState == types.Stopping && b.vadStoppingCount >= b.vadStopFrames {
		b.vadState = types.Quiet
		b.vadStoppingCount = 0

		b.isFinal = true
		b.endAtS = math.Round(float64(b.accumulateSpeechBytesLen)/float64(b.sampleNumBytes)*1000) / 1000
	}

	return &localframes.VADStateAudioRawFrame{
		AudioRawFrame: &pipelineframes.AudioRawFrame{
			DataFrame:   pipelineframes.NewDataFrameWithName("VADStateAudioRawFrame"),
			Audio:       audioBytes,
			SampleRate:  b.args.SampleRate,
			NumChannels: b.args.NumChannels,
			SampleWidth: b.args.SampleWidth,
		},
		State:    b.vadState,
		SpeechID: b.speechID,
		IsFinal:  b.isFinal,
		StartAtS: b.startAtS,
		CurAtS:   b.curAtS,
		EndAtS:   b.endAtS,
	}
}
