package llm

//https://github.com/openai/openai-go

import (
	"os"

	"github.com/openai/openai-go/v3" // imported as openai
	"github.com/openai/openai-go/v3/option"
)

type OpenAIAPIProvider struct {
	name   string
	model  string
	stream bool
	// nil, "high", "medium", "low"
	thinking *string
	client   openai.Client
}

func NewOpenAIAPIProvider(name, model string, stream bool, thinking *string) *OpenAIAPIProvider {
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	p := &OpenAIAPIProvider{
		name:     name,
		model:    model,
		stream:   stream,
		thinking: thinking,
		client:   client,
	}

	return p
}

// call /v1/completions
func (p *OpenAIAPIProvider) Completions(prompt string) string {
	return ""
}

// call /v1/chat/completions
func (p *OpenAIAPIProvider) ChatCompletions(prompt string) string {
	// TODO: implement me
	return ""
}

func (p *OpenAIAPIProvider) Name() string {
	return p.name
}

func (p *OpenAIAPIProvider) Release() {
}

func (p *OpenAIAPIProvider) Warmup() {
}
