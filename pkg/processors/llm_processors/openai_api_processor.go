package llm_processors

import (
	"context"

	"github.com/go-viper/mapstructure/v2"
	"github.com/google/uuid"
	"github.com/openai/openai-go/v3"
	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/logger"
	"github.com/weedge/pipeline-go/pkg/processors"

	"achatbot/pkg/common"
	"achatbot/pkg/types"
	achatbot_frames "achatbot/pkg/types/frames"
)

type LLMOpenAIApiProcessor struct {
	*processors.AsyncFrameProcessor
	provider common.IOpenAILLMProvider
	session  *common.Session
	mode     string
	stream   bool
	args     types.LMGenerateArgs
}

func NewLLMOpenAIApiProcessor(
	provider common.IOpenAILLMProvider, session *common.Session,
	mode string, stream bool, args types.LMGenerateArgs,
) *LLMOpenAIApiProcessor {
	if session == nil {
		session = common.NewSession(uuid.NewString(), nil)
	}
	p := &LLMOpenAIApiProcessor{
		AsyncFrameProcessor: processors.NewAsyncFrameProcessorWithPushQueueSize("LLMOpenAIApiProcessor", 128, 128),
		provider:            provider,
		session:             session,
		mode:                mode,
		stream:              stream,
		args:                args,
	}

	return p
}

// ProcessFrame processes a frame
func (p *LLMOpenAIApiProcessor) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	// call frame processor to init star frame init
	p.AsyncFrameProcessor.WithPorcessFrameAllowPush(false).ProcessFrame(frame, direction)
	switch f := frame.(type) {
	case *frames.StartFrame:
		logger.Info("LLMOpenAIApiProcessor Start")
		p.PushFrame(f, direction)
	case *frames.EndFrame:
		logger.Info("LLMOpenAIApiProcessor End")
		p.PushFrame(f, direction)
	case *frames.CancelFrame:
		logger.Info("LLMOpenAIApiProcessor Cancel")
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

func (p *LLMOpenAIApiProcessor) chat(frame *frames.TextFrame, direction processors.FrameDirection) {
	chatHistory := p.session.GetChatHistory()
	chatHistory.Append(map[string]any{"role": "user", "content": frame.Text})
	historyList := chatHistory.ToListWithoutTools() // init tools in provider
	messages := make([]types.Message, 0)
	err := mapstructure.Decode(historyList, &messages)
	if err != nil {
		logger.Error("chat", "err", err)
	}

	genContent := ""
	if !p.stream {
		p.provider.Chat(context.Background(), p.args, messages, func(resp *openai.ChatCompletion) error {
			if resp.Choices[0].Message.Content != "" {
				p.QueueFrame(frames.NewTextFrame(resp.Choices[0].Message.Content), direction)
				genContent += resp.Choices[0].Message.Content
			}
			return nil
		})
		if genContent != "" {
			chatHistory.Append(map[string]any{"role": "assistant", "content": genContent})
		}
	} else {
		p.provider.ChatStream(context.Background(), p.args, messages, func(resp *openai.ChatCompletionChunk) error {
			if resp.Choices[0].Delta.Content != "" {
				p.QueueFrame(frames.NewTextFrame(resp.Choices[0].Delta.Content), direction)
				genContent += resp.Choices[0].Delta.Content
			}
			return nil
		})
		if genContent != "" {
			chatHistory.Append(map[string]any{"role": "assistant", "content": genContent})
		}
	}

	p.QueueFrame(achatbot_frames.NewTurnEndFrame(), direction)
	logger.Debugf("ChatHistory: %+v", chatHistory.ToList())
}

func (p *LLMOpenAIApiProcessor) generate(frame *frames.TextFrame, direction processors.FrameDirection) {
	if !p.stream {
		p.provider.Generate(context.Background(), p.args, frame.Text, func(resp *openai.Completion) error {
			if resp.Choices[0].Text != "" {
				p.QueueFrame(frames.NewTextFrame(resp.Choices[0].Text), direction)
			}
			return nil
		})
	} else {
		p.provider.GenerateStream(context.Background(), p.args, frame.Text, func(resp *openai.Completion) error {
			if resp.Choices[0].Text != "" {
				p.QueueFrame(frames.NewTextFrame(resp.Choices[0].Text), direction)
			}
			return nil
		})
	}
}
