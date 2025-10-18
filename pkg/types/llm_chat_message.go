package types

import (
	"github.com/openai/openai-go/v3"
)

type Message struct {
	openai.ChatCompletionMessage `mapstructure:",squash"`
	ToolCallID                   string `json:"tool_call_id"` // hook a tool_call_id
}
