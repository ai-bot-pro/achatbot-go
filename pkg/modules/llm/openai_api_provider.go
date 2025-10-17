package llm

//https://github.com/openai/openai-go

import (
	"achatbot/pkg/common"
	"achatbot/pkg/types"
	"context"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/shared"
	"github.com/openai/openai-go/v3/shared/constant"
	"github.com/weedge/pipeline-go/pkg/logger"
)

type OpenAIAPIProvider struct {
	name   string
	model  string
	client openai.Client
}

const (
	OpenAIAPIProviderName    = "openai_api"
	OpenAIAPIProviderBaseUrl = "https://api.openai.com/v1"

	OllamaAPIProviderBaseUrl = "http://127.0.0.1:11434/v1"

	OpenRouterAIAPIProviderBaseUrl               = "https://openrouter.ai/api/v1"
	OpenRouterAIAPIProviderModelQwen3_235b_free  = "qwen/qwen3-235b-a22b:free" //have some issue,don't match openaiapi(think&tools)
	OpenRouterAIAPIProviderModelQwen2_5_72b_free = "qwen/qwen-2.5-72b-instruct:free"
)

func NewOpenAIAPIProvider(name, baseUrl, model string) *OpenAIAPIProvider {
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseUrl),
	)

	p := &OpenAIAPIProvider{
		name:   name,
		model:  model,
		client: client,
	}

	return p
}

// Generate 生成文本token
// call /v1/completions
func (p *OpenAIAPIProvider) Generate(ctx context.Context, args types.LMGenerateArgs, prompt string, respFunc common.OpenAICompletionRespFunc) {
	completion, err := p.client.Completions.New(
		context.Background(), openai.CompletionNewParams{
			Prompt:           openai.CompletionNewParamsPromptUnion{OfString: param.Opt[string]{Value: prompt}},
			Model:            openai.CompletionNewParamsModel(p.model),
			N:                param.Opt[int64]{Value: args.LmN},
			Seed:             param.Opt[int64]{Value: args.LmGenSeed},
			MaxTokens:        param.Opt[int64]{Value: args.LmMaxTokens},
			FrequencyPenalty: param.Opt[float64]{Value: args.LmGenFrequencyPenalty},
			Temperature:      param.Opt[float64]{Value: args.LmGenTemperature},
			TopP:             param.Opt[float64]{Value: args.LmGenTopP},
			Stop:             openai.CompletionNewParamsStopUnion{OfStringArray: args.LmGenStops},
		},
		// Override the header
		option.WithHeader("HTTP-Referer", "github.com/weedge"),
		option.WithHeader("X-Title", "achatbot-go"),
		option.WithMaxRetries(2), // Override the default max retries
	)
	if err != nil {
		logger.Errorf("Generate failed: %v", err)
		return
	}
	logger.Infof("%+s", completion.RawJSON())

	err = respFunc(completion)
	if err != nil {
		logger.Errorf("Generate failed: %v", err)
		return
	}
}

func (p *OpenAIAPIProvider) convertMessages(messages []types.Message) []openai.ChatCompletionMessageParamUnion {
	msgUnion := []openai.ChatCompletionMessageParamUnion{}

	for _, msg := range messages {
		switch msg.Role {
		case "system":
			msgUnion = append(msgUnion, openai.SystemMessage(msg.Content))
		case "user":
			msgUnion = append(msgUnion, openai.UserMessage(msg.Content))
		case "assistant":
			if msg.Content != "" {
				msgUnion = append(msgUnion, openai.AssistantMessage(msg.Content))
			}
			if msg.ToolCalls != nil { // tool_calls
				toolCalls := []openai.ChatCompletionMessageToolCallUnionParam{}
				for _, toolCall := range msg.ToolCalls {
					toolCalls = append(toolCalls, toolCall.ToParam())
				}
				msgUnion = append(msgUnion, openai.ChatCompletionMessageParamUnion{OfAssistant: &openai.ChatCompletionAssistantMessageParam{
					Role:      msg.Role,
					ToolCalls: toolCalls,
				}})
			}
		case "tool":
			msgUnion = append(msgUnion, openai.ChatCompletionMessageParamUnion{
				OfTool: &openai.ChatCompletionToolMessageParam{
					Role:       constant.Tool(msg.Role),
					ToolCallID: msg.ToolCallID,
					Content:    openai.ChatCompletionToolMessageParamContentUnion{OfString: param.Opt[string]{Value: msg.Content}},
				},
			})
		}
	}

	return msgUnion
}

func (p *OpenAIAPIProvider) getChatCompletionNewParams(messages []types.Message, args types.LMGenerateArgs) openai.ChatCompletionNewParams {
	params := openai.ChatCompletionNewParams{
		Messages:            p.convertMessages(messages),
		Model:               shared.ChatModel(p.model),
		PromptCacheKey:      param.Opt[string]{Value: args.PromptCacheKey},
		N:                   param.Opt[int64]{Value: args.LmN},
		Seed:                param.Opt[int64]{Value: args.LmGenSeed},
		MaxTokens:           param.Opt[int64]{Value: args.LmMaxTokens},
		MaxCompletionTokens: param.Opt[int64]{Value: args.LmGenMaxTokens},
		FrequencyPenalty:    param.Opt[float64]{Value: args.LmGenFrequencyPenalty},
		Temperature:         param.Opt[float64]{Value: args.LmGenTemperature},
		TopP:                param.Opt[float64]{Value: args.LmGenTopP},
		Stop:                openai.ChatCompletionNewParamsStopUnion{OfStringArray: args.LmGenStops},
	}
	if p.name == OpenAIAPIProviderName { //think for openai(the same as)
		if args.LmGenThinking != nil {
			switch *args.LmGenThinking {
			case "minimal":
				params.ReasoningEffort = shared.ReasoningEffortMinimal
			case "low":
				params.ReasoningEffort = shared.ReasoningEffortLow
			case "medium":
				params.ReasoningEffort = shared.ReasoningEffortMedium
			case "high":
				params.ReasoningEffort = shared.ReasoningEffortHigh
			default:
				params.ReasoningEffort = shared.ReasoningEffortMinimal
			}
		}
	}
	return params
}

// Chat 上下文chat_template 指令生成文本token
// call /v1/chat/completions
func (p *OpenAIAPIProvider) Chat(ctx context.Context,
	args types.LMGenerateArgs, messages []types.Message,
	respFunc common.OpenAIChatCompletionRespFunc,
) {
	params := p.getChatCompletionNewParams(messages, args)
	chatCompletion, err := p.client.Chat.Completions.New(ctx, params,
		// Override the header
		option.WithHeader("HTTP-Referer", "github.com/weedge"),
		option.WithHeader("X-Title", "achatbot-go"),
		option.WithMaxRetries(2), // Override the default max retries
	)
	if err != nil {
		logger.Infof("Chat failed: %v", err)
		return
	}
	logger.Infof("%s", chatCompletion.RawJSON())

	err = respFunc(chatCompletion)
	if err != nil {
		logger.Infof("Chat failed: %v", err)
		return
	}
}

// GenerateStream stream generate 生成文本token
func (p *OpenAIAPIProvider) GenerateStream(ctx context.Context, args types.LMGenerateArgs, prompt string, respFunc common.OpenAIStreamCompletionRespFunc) {
	stream := p.client.Completions.NewStreaming(
		ctx, openai.CompletionNewParams{
			Prompt:           openai.CompletionNewParamsPromptUnion{OfString: param.Opt[string]{Value: prompt}},
			N:                param.Opt[int64]{Value: args.LmN},
			Seed:             param.Opt[int64]{Value: args.LmGenSeed},
			MaxTokens:        param.Opt[int64]{Value: args.LmMaxTokens},
			FrequencyPenalty: param.Opt[float64]{Value: args.LmGenFrequencyPenalty},
			Temperature:      param.Opt[float64]{Value: args.LmGenTemperature},
			TopP:             param.Opt[float64]{Value: args.LmGenTopP},
			Stop:             openai.CompletionNewParamsStopUnion{OfStringArray: args.LmGenStops},
		},

		// Override the header
		option.WithHeader("HTTP-Referer", "github.com/weedge"),
		option.WithHeader("X-Title", "achatbot-go"),
		option.WithMaxRetries(2), // Override the default max retries
	)
	for stream.Next() {
		chunk := stream.Current()
		err := respFunc(&chunk)
		if err != nil {
			logger.Errorf("generate stream error: %v", err)
		}
	}
}

// ChatStream stream chat 上下文chat_template 指令生成文本token
func (p *OpenAIAPIProvider) ChatStream(ctx context.Context, args types.LMGenerateArgs, messages []types.Message, respFunc common.OpenAIStreamChatCompletionRespFunc) {

	params := p.getChatCompletionNewParams(messages, args)
	stream := p.client.Chat.Completions.NewStreaming(ctx, params,
		// Override the header
		option.WithHeader("HTTP-Referer", "github.com/weedge"),
		option.WithHeader("X-Title", "achatbot-go"),
		option.WithMaxRetries(2), // Override the default max retries
	)

	for stream.Next() {
		chunk := stream.Current()
		err := respFunc(&chunk)
		if err != nil {
			logger.Errorf("chat stream error: %v", err)
		}
	}

	if stream.Err() != nil {
		logger.Errorf("chat stream error: %v", stream.Err())
		return
	}
}

func (p *OpenAIAPIProvider) Name() string {
	return p.name
}
