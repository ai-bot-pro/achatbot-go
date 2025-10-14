package llm

import (
	"context"
	"maps"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/weedge/pipeline-go/pkg/logger"
)

type OllamaAPIProvider struct {
	name     string
	model    string
	stream   bool
	thinking *string   // nil, "high", "medium", "low"
	tools    api.Tools // for chat with tools
	genArgs  map[string]any
	client   *api.Client
}

const (
	OllamaAPIProviderName            = "ollama_api"
	OllamaAPIProviderModel_QWEN3_0_6 = "qwen3:0.6b"
)

func NewOllamaAPIProviderWithoutTools(
	name, model string, stream bool, thinking *string, genArgs map[string]any) *OllamaAPIProvider {
	return NewOllamaAPIProvider(name, model, stream, thinking, nil, genArgs)
}

func NewOllamaAPIProvider(name, model string, stream bool, thinking *string, tools api.Tools, genArgs map[string]any) *OllamaAPIProvider {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		logger.Error("NewOllamaAPIProvider failed", "error", err)
		return nil
	}

	p := &OllamaAPIProvider{
		name:     name,
		model:    model,
		stream:   stream,
		thinking: thinking,
		client:   client,
		tools:    tools,
		genArgs:  genArgs,
	}

	return p
}

// Generate call /api/generate
func (p *OllamaAPIProvider) Generate(ctx context.Context, prompt string, respFunc api.GenerateResponseFunc) {
	think := &api.ThinkValue{Value: false} // no thinking
	if p.thinking != nil {
		think = &api.ThinkValue{Value: strings.ToLower(*p.thinking)}
	}
	req := &api.GenerateRequest{
		Model:   p.model,
		Prompt:  prompt,
		Think:   think,
		Options: p.genArgs,
	}
	if !p.stream {
		// set streaming to false
		req.Stream = new(bool)
	}

	err := p.client.Generate(ctx, req, respFunc)
	if err != nil {
		logger.Error("Generate failed", "req", req, "error", err)
	}
}

// Generate call /api/chat
func (p *OllamaAPIProvider) Chat(ctx context.Context, messages []api.Message, respFunc api.ChatResponseFunc) {
	think := &api.ThinkValue{Value: false} // no thinking
	if p.thinking != nil {
		think = &api.ThinkValue{Value: strings.ToLower(*p.thinking)}
	}
	req := &api.ChatRequest{
		Model:    p.model,
		Messages: messages,
		Think:    think,
		Options:  p.genArgs,
		Tools:    p.tools,
	}
	if !p.stream {
		// set streaming to false
		req.Stream = new(bool)
	}

	err := p.client.Chat(ctx, req, respFunc)
	if err != nil {
		logger.Error("Chat failed", "req", req, "error", err)
	}
}

func (p *OllamaAPIProvider) Name() string {
	return p.name
}

func (p *OllamaAPIProvider) Release() {
}

func (p *OllamaAPIProvider) Warmup() {
}

// UpdateGenArgs updates the GenerateArgs with the provided values
func (p *OllamaAPIProvider) UpdateGenArgs(values map[string]any) {
	maps.Copy(p.genArgs, values)
}
