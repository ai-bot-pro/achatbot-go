package common

import (
	"context"

	"github.com/openai/openai-go/v3"
	"github.com/weedge/pipeline-go/pkg/frames"

	"achatbot/pkg/consts"
	"achatbot/pkg/types"
	achatbot_frames "achatbot/pkg/types/frames"
)

type IPoolInstance interface {
	// Reset 重置池中的实例状态
	Reset() error

	// Destroy 销毁池中的实例
	Release() error
}

// IVADAnalyzer 定义了语音活动检测（VAD）分析器的接口。
type IVADAnalyzer interface {
	// AnalyzeAudio 对输入的音频缓冲区进行分析，返回 VAD Frame
	AnalyzeAudio(buffer []byte) *achatbot_frames.VADStateAudioRawFrame

	// GetSampleRate 返回 VAD 的采样率。
	GetSampleRate() int

	// GetWindowSize 返回 VAD 的采样窗口大小。
	GetWindowSize() int

	// Reset 重置 VAD 的统计信息和模型状态。
	// Release 释放资源。
	IPoolInstance
}

// IVoiceConfidenceProvider local语音置信度提供者接口
type IVoiceConfidenceProvider interface {
	// IsActiveSpeech 判断当前音频是否是活跃的语音。
	IsActiveSpeech(audio []byte) bool

	// Warmup 预热
	Warmup()

	// GetSampleInfo 返回语音置信度提供者的采样率信息和采样窗口大小。
	GetSampleInfo() (int, int)

	// Name 返回语音置信度提供者的名称。
	Name() string

	// Reset 重置 VAD 模型状态(单轮识别)。
	// Release 释放资源。
	IPoolInstance
}

// ------------------------------------------------------------

// IASRProvider local语音识别提供者接口
type IASRProvider interface {
	// Transcribe 语音转录文本
	Transcribe(audio []byte) string

	// Warmup 预热
	Warmup()

	// Name 返回语音识别提供者的名称。
	Name() string

	// Release 释放资源。
	// Reset 重置状态。
	IPoolInstance
}

// ------------------------------------------------------------

type OpenAIStreamChatCompletionRespFunc func(*openai.ChatCompletionChunk) error
type OpenAIChatCompletionRespFunc func(*openai.ChatCompletion) error
type OpenAIStreamCompletionRespFunc func(*openai.Completion) error
type OpenAICompletionRespFunc func(*openai.Completion) error

// IOpenAILLMProvider remote生成模型提供者接口
type IOpenAILLMProvider interface {
	// generate 生成文本token
	Generate(ctx context.Context, args types.LMGenerateArgs, prompt string, respFunc OpenAICompletionRespFunc)

	// chat 上下文chat_template 指令生成文本token
	Chat(ctx context.Context, args types.LMGenerateArgs, messages []types.Message, respFunc OpenAIChatCompletionRespFunc)

	// stream generate 生成文本token
	GenerateStream(ctx context.Context, args types.LMGenerateArgs, prompt string, respFunc OpenAIStreamCompletionRespFunc)

	// stream chat 上下文chat_template 指令生成文本token
	ChatStream(ctx context.Context, args types.LMGenerateArgs, messages []types.Message, respFunc OpenAIStreamChatCompletionRespFunc)

	// Name 返回生成文本token提供者的名称。
	Name() string
}

// ------------------------------------------------------------

// ITTSProvider 文本合成语音提供者接口
type ITTSProvider interface {
	// Synthesize 文本合成语音
	Synthesize(text string) []byte

	// Warmup 预热
	Warmup()

	// GetSampleInfo 返回合成语音采样率信息(sample_rate), 通道数(channels), sample_width/bit_depth(样本宽度/位深度)
	GetSampleInfo() (int, int, int)

	// SetPromptAudio 设置Prompt Audio shot示例样本去生成对应特征的声音
	SetPromptAudio(string, []byte) error

	// Name 返回文本合成语音提供者的名称。
	Name() string

	// Release 释放资源。
	// Reset 重置状态。
	IPoolInstance
}

// --------------------------------------------------------------------

// We'll use the standard net/http package for WebSocket support
// You may need to add the gorilla/websocket dependency or use standard library
// For now, we'll define a generic interface
// IWebSocketConn defines the interface for WebSocket connections
type IWebSocketConn interface {
	ReadMessage() (messageType consts.MessageType, p []byte, err error)
	WriteMessage(messageType consts.MessageType, data []byte) error
	Close() error
}

type ITransportWriter interface {
	WriteRawAudio(data []byte) error

	WriteFrame(frame frames.Frame) error

	//WriteAnimationAudioFrame(frame *achatbot_frames.AnimationAudioRawFrame) error
	//SendMessage(frame *achatbot_frames.TransportMessageFrame) error
	//SendText(frame *frames.TextFrame) error
	//WriteImageFrame(frame *frames.ImageRawFrame) error
}

type IFunction interface {
	// Execute 执行函数
	Execute(args map[string]any) (string, error)

	// GetToolCall 获取 openai 标准 tool call schema
	GetToolCall() map[string]any

	// GetOllamaAPIToolCall 获取 ollama 自定义的 toolcall schema
	GetOllamaAPIToolCall() map[string]any
}
