package processors

import (
	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/logger"
	"github.com/weedge/pipeline-go/pkg/processors"

	"achatbot/pkg/common"
)

type TTSProcessor struct {
	*processors.AsyncFrameProcessor
	provider common.ITTSProvider
}

func NewTTSProcessor(provider common.ITTSProvider) *TTSProcessor {
	return &TTSProcessor{
		AsyncFrameProcessor: processors.NewAsyncFrameProcessor("TTSProcessor"),
		provider:            provider,
	}
}

func (p *TTSProcessor) WithPassText(passText bool) *TTSProcessor {
	p.AsyncFrameProcessor = p.AsyncFrameProcessor.WithPassText(passText)
	return p
}

func (p *TTSProcessor) Start(frame *frames.StartFrame) {
	logger.Info("TTSProcessor Start")
}

func (p *TTSProcessor) Stop(frame *frames.EndFrame) {
	logger.Info("TTSProcessor Stop")
}

func (p *TTSProcessor) Cancel(frame *frames.CancelFrame) {
	p.provider.Release()
	logger.Info("TTSProcessor Cancel")
}

// ProcessFrame processes a frame
func (p *TTSProcessor) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
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
	case *frames.TextFrame:
		if p.PassText() {
			p.QueueFrame(f, direction)
		}
		audio := p.provider.Synthesize(f.Text)
		rate, channels, sampleWidth := p.provider.GetSampleInfo()
		p.PushDownstreamFrame(frames.NewAudioRawFrame(audio, rate, channels, sampleWidth))
	default:
		p.QueueFrame(f, direction)
	}

}
