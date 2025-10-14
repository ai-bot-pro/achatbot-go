package llm_processors

import (
	"context"

	"github.com/go-viper/mapstructure/v2"
	"github.com/google/uuid"
	"github.com/ollama/ollama/api"
	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/logger"
	"github.com/weedge/pipeline-go/pkg/processors"

	"achatbot/pkg/common"
	"achatbot/pkg/modules/llm"
	achatbot_frames "achatbot/pkg/types/frames"
)

type LLMOllamaApiProcessor struct {
	*processors.AsyncFrameProcessor
	provider *llm.OllamaAPIProvider
	session  *common.Session
	mode     string
}

const (
	Mode_Generate = "generate"
	Mode_Chat     = "chat"
)

func NewLLMOllamaApiProcessor(provider *llm.OllamaAPIProvider, session *common.Session, mode string) *LLMOllamaApiProcessor {
	if session == nil {
		session = common.NewSession(uuid.NewString(), nil)
	}
	p := &LLMOllamaApiProcessor{
		AsyncFrameProcessor: processors.NewAsyncFrameProcessorWithPushQueueSize("LLMOllamaApiProcessor", 128, 128),
		provider:            provider,
		session:             session,
		mode:                mode,
	}

	return p
}

// ProcessFrame processes a frame
func (p *LLMOllamaApiProcessor) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	// call frame processor to init star frame init
	p.AsyncFrameProcessor.WithPorcessFrameAllowPush(false).ProcessFrame(frame, direction)
	switch f := frame.(type) {
	case *frames.StartFrame:
		logger.Info("LLMOllamaApiProcessor Start")
		p.PushFrame(f, direction)
	case *frames.EndFrame:
		logger.Info("LLMOllamaApiProcessor End")
		p.PushFrame(f, direction)
	case *frames.CancelFrame:
		logger.Info("LLMOllamaApiProcessor Cancel")
		p.PushFrame(f, direction)
	case *frames.TextFrame:
		switch p.mode {
		case "chat":
			p.chat(f, direction)
		case "generate":
			p.generate(f, direction)
		}
	default:
		p.QueueFrame(f, direction)
	}
}

func (p *LLMOllamaApiProcessor) chat(frame *frames.TextFrame, direction processors.FrameDirection) {
	chatHistory := p.session.GetChatHistory()
	chatHistory.Append(map[string]any{"role": "user", "content": frame.Text})
	historyList := chatHistory.ToListWithoutTools() // init tools in provider
	messages := make([]api.Message, 0)
	err := mapstructure.Decode(historyList, &messages)
	if err != nil {
		logger.Error("chat", "err", err)
	}

	genThinking, genContent := "", ""
	p.provider.Chat(context.Background(), messages, func(resp api.ChatResponse) error {
		if resp.Message.Thinking != "" {
			p.QueueFrame(achatbot_frames.NewThinkTextFrame(resp.Message.Thinking), direction)
			genThinking += resp.Message.Thinking
		}
		if resp.Message.Content != "" {
			p.QueueFrame(frames.NewTextFrame(resp.Message.Content), direction)
			genContent += resp.Message.Content
		}
		return nil
	})
	if genContent != "" {
		chatHistory.Append(map[string]any{"role": "assistant", "content": genContent})
	}
	if genThinking != "" {
		chatHistory.Append(map[string]any{"role": "assistant", "thinking": genThinking})
	}
	p.QueueFrame(achatbot_frames.NewTurnEndFrame(), direction)
}

func (p *LLMOllamaApiProcessor) generate(frame *frames.TextFrame, direction processors.FrameDirection) {
	p.provider.Generate(context.Background(), frame.Text, func(resp api.GenerateResponse) error {
		if resp.Thinking != "" {
			p.QueueFrame(achatbot_frames.NewThinkTextFrame(resp.Thinking), direction)
		}
		if resp.Response != "" {
			p.QueueFrame(frames.NewTextFrame(resp.Response), direction)
		}
		return nil
	})
}
