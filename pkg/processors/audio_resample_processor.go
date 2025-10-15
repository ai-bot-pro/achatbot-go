package processors

import (
	"achatbot/pkg/utils"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/processors"
)

type AudioResampleProcessor struct {
	*processors.AsyncFrameProcessor
	outRate int
}

func NewAudioResampleProcessor(outRate int) *AudioResampleProcessor {
	return &AudioResampleProcessor{
		AsyncFrameProcessor: processors.NewAsyncFrameProcessor("AudioResampleProcessor"),
		outRate:             outRate,
	}
}

// ProcessFrame processes a frame
func (p *AudioResampleProcessor) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	// call frame processor to init star frame init
	p.AsyncFrameProcessor.WithPorcessFrameAllowPush(false).ProcessFrame(frame, direction)

	switch f := frame.(type) {
	case *frames.StartFrame:
		p.PushFrame(f, direction)
	case *frames.EndFrame:
		p.PushFrame(f, direction)
	case *frames.CancelFrame:
		p.PushFrame(f, direction)
	case *frames.AudioRawFrame:
		if f.SampleRate != p.outRate {
			f.Audio = utils.ResampleBytes(f.Audio, f.SampleRate, p.outRate)
			f.SampleRate = p.outRate
			if f.NumChannels > 0 && f.SampleWidth > 0 {
				f.NumFrames = len(f.Audio) / (f.NumChannels * f.SampleWidth)
			}
		}
		p.QueueFrame(f, direction)
	default:
		p.QueueFrame(f, direction)
	}

}
