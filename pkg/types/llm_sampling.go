package types

// LMGenerateArgs lm generation sampling arguments
type LMGenerateArgs struct {
	LmN int64 `json:"lm_n" default:"1"`

	// LmGenSeed The seed to use for the language model. Default is 42.
	LmGenSeed int64 `json:"lm_gen_seed" default:"42"`

	// LmMaxTokens Maximum number of tokens(promts+new tokens) to generate in a single completion. Default is 2048.
	LmMaxTokens int64 `json:"lm_max_tokens" default:"2048"`

	// LmGenMaxTokens Maximum number of new tokens to generate in a single completion. Default is 2048.
	LmGenMaxTokens int64 `json:"lm_gen_max_tokens" default:"2048"`

	// LmGenReasoningMaxTokens Maximum number of new tokens to generate in a single completion when thinking.
	// reasoning_max_tokens must be less than max_tokens, Default is 1024.
	LmGenReasoningMaxTokens int64 `json:"lm_gen_reasoning_max_tokens" default:"1024"`

	// LmGenTemperature Controls the randomness of the output. Set to 0.0 for deterministic (repeatable) outputs. Default is 0.1.
	LmGenTemperature float64 `json:"lm_gen_temperature" default:"0.6"`

	// LmGenTopP Top-p is usually set to a high value (like 0.75) with the purpose of limiting the long tail of low-probability tokens that may be sampled.
	// We can use both top-k and top-p together. If both k and p are enabled, p acts after k. Default is 0.8.
	LmGenTopP float64 `json:"lm_gen_top_p" default:"0.9"`

	// LmGenStops A list of strings that will stop the generation. Default is []. If the stop word is a substring of the generated text, the generation will stop.
	LmGenStops []string `json:"lm_gen_stops"`

	// Number between -2.0 and 2.0. Positive values penalize new tokens based on their
	// existing frequency in the text so far, decreasing the model's likelihood to
	// repeat the same line verbatim.
	LmGenFrequencyPenalty float64 `json:"lm_gen_frequency_penalty" default:"0.0"`
	LmGenPresencePenalty  float64 `json:"lm_gen_presence_penalty" default:"0.0"`

	// nil(if llm api support no thinking) or Any of "minimal", "low", "medium", "high".
	LmGenThinking *string `json:"lm_gen_thinking"`

	// [Learn more](https://platform.openai.com/docs/guides/prompt-caching).
	PromptCacheKey string `json:"prompt_cache_key,omitzero" default:""`
}

// NewLMGenerateArgs creates a new LMGenerateArgs with default values
func NewLMGenerateArgs() *LMGenerateArgs {
	return &LMGenerateArgs{
		LmN:                     1,
		LmGenSeed:               42,
		LmGenMaxTokens:          2048,
		LmGenReasoningMaxTokens: 1024,
		LmGenTemperature:        0.6,
		LmGenTopP:               0.9,
		LmGenStops:              []string{},
		LmGenFrequencyPenalty:   0.0,
		LmGenPresencePenalty:    0.0,
		LmGenThinking:           nil,
		PromptCacheKey:          "",
	}
}
