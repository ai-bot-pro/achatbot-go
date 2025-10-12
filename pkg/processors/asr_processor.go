package processors

import (
	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/logger"
	"github.com/weedge/pipeline-go/pkg/processors"

	"achatbot/pkg/common"
	achatbot_frames "achatbot/pkg/types/frames"
)

type ASRProcessor struct {
	*processors.AsyncFrameProcessor
	provider     common.IASRProvider
	passRawAudio bool
}

func NewASRProcessor(provider common.IASRProvider) *ASRProcessor {
	return &ASRProcessor{
		AsyncFrameProcessor: processors.NewAsyncFrameProcessor("ASRProcessor"),
		provider:            provider,
		passRawAudio:        false,
	}
}

func (p *ASRProcessor) WithPassRawAudio(passRawAudio bool) *ASRProcessor {
	p.passRawAudio = passRawAudio
	return p
}

func (p *ASRProcessor) Start(frame *frames.StartFrame) {
	logger.Info("ASRProcessor Start")
}

func (p *ASRProcessor) Stop(frame *frames.EndFrame) {
	logger.Info("ASRProcessor Stop")
}

func (p *ASRProcessor) Cancel(frame *frames.CancelFrame) {
	p.provider.Release()
	logger.Info("ASRProcessor Cancel")
}

// ProcessFrame processes a frame
func (p *ASRProcessor) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
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
	case *frames.AudioRawFrame:
		if p.passRawAudio {
			p.QueueFrame(f, direction)
		}
		text := p.provider.Transcribe(f.Audio)
		p.PushDownstreamFrame(frames.NewTextFrame(text))
	case *achatbot_frames.VADStateAudioRawFrame:
		if p.passRawAudio {
			p.QueueFrame(f, direction)
		}
		text := p.provider.Transcribe(f.Audio)
		p.PushDownstreamFrame(frames.NewTextFrame(text))
	case *achatbot_frames.AnimationAudioRawFrame:
		if p.passRawAudio {
			p.QueueFrame(f, direction)
		}
		text := p.provider.Transcribe(f.Audio)
		p.PushDownstreamFrame(frames.NewTextFrame(text))
	default:
		p.QueueFrame(f, direction)
	}

}
