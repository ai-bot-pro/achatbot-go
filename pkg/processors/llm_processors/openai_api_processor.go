package llm_processors

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/google/uuid"
	"github.com/openai/openai-go/v3"
	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/logger"
	"github.com/weedge/pipeline-go/pkg/processors"

	"achatbot/pkg/common"
	"achatbot/pkg/modules/functions"
	"achatbot/pkg/types"
	achatbot_frames "achatbot/pkg/types/frames"
)

type LLMOpenAIApiProcessor struct {
	*processors.AsyncFrameProcessor
	provider       common.IOpenAILLMProvider
	session        *common.Session
	mode           string
	stream         bool
	args           types.LMGenerateArgs
	isHistoryThink bool
}

func NewLLMOpenAIApiProcessor(
	provider common.IOpenAILLMProvider, session *common.Session,
	mode string, stream bool, args types.LMGenerateArgs,
) *LLMOpenAIApiProcessor {
	if session == nil {
		session = common.NewSession(uuid.NewString(), nil)
	}
	p := &LLMOpenAIApiProcessor{
		AsyncFrameProcessor: processors.NewAsyncFrameProcessorWithPushQueueSize("LLMOpenAIApiProcessor", 1024, 1024),
		provider:            provider,
		session:             session,
		mode:                mode,
		stream:              stream,
		args:                args,
		isHistoryThink:      false,
	}

	return p
}

func (p *LLMOpenAIApiProcessor) WithIsHistoryThink(isHistoryThink bool) *LLMOpenAIApiProcessor {
	p.isHistoryThink = isHistoryThink
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

// appendHistoryChatMessages message(api.Message) append to history list([]map[string]any)
func (p *LLMOpenAIApiProcessor) appendHistoryChatMessages(msgs []types.Message) {
	for _, msg := range msgs {
		mapMsg := map[string]any{}
		err := mapstructure.Decode(msg, &mapMsg)
		if err != nil {
			logger.Errorf("mapstructure.Decode error: %v", err)
			continue
		}
		p.session.GetChatHistory().Append(mapMsg)
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

	isToolCalls := true
	cnToolCalls := 0
	for isToolCalls {
		if cnToolCalls > 3 {
			logger.Error("chat", "err", "too many tool calls")
			break
		}
		cnToolCalls++
		if !p.stream {
			p.provider.Chat(context.Background(), p.args, messages, func(resp *openai.ChatCompletion) error {
				toolMsgs := []types.Message{}
				for i, toolCall := range resp.Choices[0].Message.ToolCalls {
					// Extract the location from the function call arguments
					funcArgs := strings.ReplaceAll(toolCall.Function.Arguments, "{}", "")
					resp.Choices[0].Message.ToolCalls[i].Function.Arguments = funcArgs

					var args map[string]any
					err := json.Unmarshal([]byte(funcArgs), &args)
					if err != nil {
						logger.Errorf("Failed to unmarshal function arguments: %v err: %v", funcArgs, err)
						continue
					}
					result, err := functions.RegisterFuncs.Execute(toolCall.Function.Name, args)
					if err != nil {
						logger.Errorf("Failed to execute function: %v err: %v", toolCall.Function.Name, err)
						continue
					}
					toolMsgs = append(toolMsgs, types.Message{
						ChatCompletionMessage: openai.ChatCompletionMessage{Role: "tool", Content: result},
						ToolCallID:            toolCall.ID,
					})
					p.QueueFrame(achatbot_frames.NewFunctionCallFrame(toolCall.ID, toolCall.Function.Name, args, i), direction)
				}
				// If there is a was a function call, continue the conversation
				if len(toolMsgs) > 0 { //call_tools
					if !p.isHistoryThink {
						resp.Choices[0].Message.Reasoning = ""
					}
					msg := types.Message{ChatCompletionMessage: resp.Choices[0].Message}
					messages = append(messages, msg)
					p.appendHistoryChatMessages([]types.Message{msg})
					messages = append(messages, toolMsgs...)
					p.appendHistoryChatMessages(toolMsgs)
					isToolCalls = true
				}
				if resp.Choices[0].Message.Reasoning != "" {
					p.QueueFrame(achatbot_frames.NewThinkTextFrame(resp.Choices[0].Message.Reasoning), direction)
				}
				if resp.Choices[0].Message.Content != "" {
					isToolCalls = false
					if !p.isHistoryThink {
						resp.Choices[0].Message.Reasoning = ""
					}
					msg := types.Message{ChatCompletionMessage: resp.Choices[0].Message}
					messages = append(messages, msg)
					p.appendHistoryChatMessages([]types.Message{msg})
					p.QueueFrame(frames.NewTextFrame(resp.Choices[0].Message.Content), direction)
				}
				return nil
			})
		} else { //stream
			acc := openai.ChatCompletionAccumulator{}
			toolMsgs := []types.Message{}
			p.provider.ChatStream(context.Background(), p.args, messages, func(chunk *openai.ChatCompletionChunk) error {
				acc.AddChunk(*chunk)
				if len(chunk.Choices) == 0 {
					return nil
				}

				if chunk.Choices[0].Delta.Reasoning != "" {
					p.QueueFrame(achatbot_frames.NewThinkTextFrame(chunk.Choices[0].Delta.Reasoning), direction)
				}
				if chunk.Choices[0].Delta.Content != "" {
					p.QueueFrame(frames.NewTextFrame(chunk.Choices[0].Delta.Content), direction)
				}

				if chunk.Choices[0].Delta.ToolCalls != nil {
					for _, tool := range chunk.Choices[0].Delta.ToolCalls {
						tool.Function.Arguments = strings.ReplaceAll(tool.Function.Arguments, "{}", "")
						var args map[string]any
						err := json.Unmarshal([]byte(tool.Function.Arguments), &args)
						if err != nil {
							logger.Errorf("Failed to Unmarshal err: %v", err)
							continue
						}
						result, err := functions.RegisterFuncs.Execute(tool.Function.Name, args)
						if err != nil {
							logger.Error("Execute", "err", err, "funcName", tool.Function.Name, "funcArgs", tool.Function.Arguments)
							continue
						}
						toolMsgs = append(toolMsgs, types.Message{
							ChatCompletionMessage: openai.ChatCompletionMessage{Role: "tool", Content: result},
							ToolCallID:            tool.ID,
						})
						p.QueueFrame(achatbot_frames.NewFunctionCallFrame(tool.ID, tool.Function.Name, args, int(tool.Index)), direction)
					}
				}
				return nil
			})
			// If there is a was a function call, continue the conversation
			if len(toolMsgs) > 0 { //call_tools
				if !p.isHistoryThink {
					acc.Choices[0].Message.Reasoning = ""
				}
				msg := types.Message{ChatCompletionMessage: acc.Choices[0].Message}
				messages = append(messages, msg)
				p.appendHistoryChatMessages([]types.Message{msg})
				messages = append(messages, toolMsgs...)
				p.appendHistoryChatMessages(toolMsgs)
				isToolCalls = true
				cnToolCalls++
			}

			if len(acc.Choices) > 0 && acc.Choices[0].Message.Content != "" {
				if !p.isHistoryThink {
					acc.Choices[0].Message.Reasoning = ""
				}
				msg := types.Message{ChatCompletionMessage: acc.Choices[0].Message}
				messages = append(messages, msg)
				p.appendHistoryChatMessages([]types.Message{msg})
				isToolCalls = false
			}
		} //end stream
	} //end call

	p.QueueFrame(achatbot_frames.NewTurnEndFrame(), direction)
	logger.Infof("ChatHistory: %+v", p.session.GetChatHistory().ToList())
	p.session.IncrementChatRound()
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
	p.QueueFrame(achatbot_frames.NewTurnEndFrame(), direction)
}
