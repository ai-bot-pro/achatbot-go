package params

import "achatbot/pkg/consts"

// VADAnalyzerArgs VAD分析器参数
type VADAnalyzerArgs struct {
	SampleRate  int     `json:"sample_rate"`
	NumChannels int     `json:"num_channels"`
	SampleWidth int     `json:"sample_width"`
	StartSecs   float64 `json:"start_secs"`
	StopSecs    float64 `json:"stop_secs"`
}

// NewVADAnalyzerArgs 创建一个新的VADAnalyzerArgs实例，带有默认值
func NewVADAnalyzerArgs() *VADAnalyzerArgs {
	return &VADAnalyzerArgs{
		SampleRate:  consts.DefaultRate,
		NumChannels: consts.DefaultChannels,
		SampleWidth: consts.DefaultSampleWidth,
		StartSecs:   0.032, // default use SileroVAD 32ms start once for 16000 samples, 512 frames per second, accumulate 1 times
		StopSecs:    0.32,  // default use SileroVAD  32ms stop once for 16000 samples, 512 frames per second, accumulate 10 times
	}
}

// WithSampleRate 设置采样率
func (args *VADAnalyzerArgs) WithSampleRate(sampleRate int) *VADAnalyzerArgs {
	args.SampleRate = sampleRate
	return args
}

// WithNumChannels 设置声道数
func (args *VADAnalyzerArgs) WithNumChannels(numChannels int) *VADAnalyzerArgs {
	args.NumChannels = numChannels
	return args
}

// WithSampleWidth 设置采样宽度
func (args *VADAnalyzerArgs) WithSampleWidth(sampleWidth int) *VADAnalyzerArgs {
	args.SampleWidth = sampleWidth
	return args
}

// WithStartSecs 设置开始秒数
func (args *VADAnalyzerArgs) WithStartSecs(startSecs float64) *VADAnalyzerArgs {
	args.StartSecs = startSecs
	return args
}

// WithStopSecs 设置停止秒数
func (args *VADAnalyzerArgs) WithStopSecs(stopSecs float64) *VADAnalyzerArgs {
	args.StopSecs = stopSecs
	return args
}
