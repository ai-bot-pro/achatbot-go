package vad_analyzer

import (
	"math"

	pipelineframes "github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/logger"

	"achatbot/pkg/common"
	"achatbot/pkg/params"
	"achatbot/pkg/types"
	localframes "achatbot/pkg/types/frames"
)

// VADAnalyzer 基础VAD分析器
type VADAnalyzer struct {
	args                     *params.VADAnalyzerArgs
	vadFrames                int
	vadFramesNumBytes        int
	sampleNumBytes           int
	vadStartFrames           int
	vadStopFrames            int
	vadStartingCount         int
	vadStoppingCount         int
	vadState                 types.VADState
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
func NewVADAnalyzer(args *params.VADAnalyzerArgs, vcp common.IVoiceConfidenceProvider) *VADAnalyzer {
	analyzer := &VADAnalyzer{
		args:                     args,
		vadState:                 types.Quiet,
		IVoiceConfidenceProvider: vcp,
	}

	analyzer.vadFrames = analyzer.GetWindowSize()
	analyzer.vadFramesNumBytes = analyzer.vadFrames * analyzer.args.NumChannels * analyzer.args.SampleWidth
	analyzer.sampleNumBytes = analyzer.args.SampleRate * analyzer.args.NumChannels * analyzer.args.SampleWidth

	vadFramesPerSec := float64(analyzer.vadFrames) / float64(analyzer.args.SampleRate)
	analyzer.vadStartFrames = int(math.Round(analyzer.args.StartSecs / vadFramesPerSec))
	analyzer.vadStopFrames = int(math.Round(analyzer.args.StopSecs / vadFramesPerSec))

	analyzer.Reset()
	return analyzer
}

// Reset 重置分析器状态
func (b *VADAnalyzer) Reset() error {
	b.vadStartingCount = 0
	b.vadStoppingCount = 0
	b.vadState = types.Quiet
	b.speechID = 0
	b.accumulateSpeechBytesLen = 0
	b.isFinal = false
	b.startAtS = 0.0
	b.curAtS = 0.0
	b.endAtS = 0.0

	return b.IVoiceConfidenceProvider.Reset()
}

func (b *VADAnalyzer) GetSampleRate() int {
	sr, _ := b.IVoiceConfidenceProvider.GetSampleInfo()
	return sr
}

func (b *VADAnalyzer) Release() error {
	return b.IVoiceConfidenceProvider.Release()
}

// GetWindowSize 计算需要的帧数
func (b *VADAnalyzer) GetWindowSize() int {
	_, windowSize := b.IVoiceConfidenceProvider.GetSampleInfo()
	return windowSize
}

// isActiveSpeech 判断是否是活跃的语音
func (b *VADAnalyzer) isActiveSpeech(audio []byte) bool {
	return b.IVoiceConfidenceProvider.IsActiveSpeech(audio)
}

// AnalyzeAudio 分析音频
func (b *VADAnalyzer) AnalyzeAudio(buffer []byte) *localframes.VADStateAudioRawFrame {
	if len(buffer) != b.vadFramesNumBytes {
		logger.Warnf("VADAnalyzer: buffer size MisMatch: %d != %d", len(buffer), b.vadFramesNumBytes)
		return nil
	}

	b.curAtS = math.Round(float64(b.accumulateSpeechBytesLen)/float64(b.sampleNumBytes)*1000) / 1000
	b.accumulateSpeechBytesLen += len(buffer)

	speaking := b.isActiveSpeech(buffer)
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

	vadStateAudioRawFrame := &localframes.VADStateAudioRawFrame{
		AudioRawFrame: &pipelineframes.AudioRawFrame{
			DataFrame:   pipelineframes.NewDataFrameWithName("VADStateAudioRawFrame"),
			Audio:       buffer,
			SampleRate:  b.args.SampleRate,
			NumChannels: b.args.NumChannels,
			SampleWidth: b.args.SampleWidth,
			NumFrames:   len(buffer) / (b.args.NumChannels * b.args.SampleWidth),
		},
		State:    b.vadState,
		SpeechID: b.speechID,
		IsFinal:  b.isFinal,
		StartAtS: b.startAtS,
		CurAtS:   b.curAtS,
		EndAtS:   b.endAtS,
	}
	return vadStateAudioRawFrame
}
