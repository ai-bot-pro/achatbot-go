package llm

// https://ai.google.dev/gemini-api/docs/quickstart?hl=zh-cn
type GeminiAPIProvider struct {
	name   string
	model  string
	stream bool
	// nil, "high", "medium", "low"
	thinking *string
}

func NewGeminiAPIProvider(name, model string, stream bool, thinking *string) *GeminiAPIProvider {

	p := &GeminiAPIProvider{
		name:     name,
		model:    model,
		stream:   stream,
		thinking: thinking,
	}

	return p
}

func (p *GeminiAPIProvider) Generate(prompt string) string {
	return ""
}

func (p *GeminiAPIProvider) Chat(prompt string) string {
	// TODO: implement me
	return ""
}

func (p *GeminiAPIProvider) Name() string {
	return p.name
}

func (p *GeminiAPIProvider) Release() {
	// TODO: implement me
}

func (p *GeminiAPIProvider) Warmup() {
	// TODO: implement me
}
