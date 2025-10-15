package llm

// LMGenerateArgs lm generation sampling arguments
type LMGenerateArgs struct {
	// LmGenSeed The seed to use for the language model. Default is 42.
	LmGenSeed int `json:"lm_gen_seed" default:"42"`

	// LmMaxLength Corresponds to the length of the input prompt + max_new_tokens.
	// Its effect is overridden by max_new_tokens, if also set. Default is 2048.
	LmMaxLength int `json:"lm_max_length" default:"2048"`

	// LmGenMaxTokens Maximum number of new tokens to generate in a single completion. Default is 2048.
	LmGenMaxTokens int `json:"lm_gen_max_tokens" default:"2048"`

	// LmGenReasoningMaxTokens Maximum number of new tokens to generate in a single completion when thinking.
	// reasoning_max_tokens must be less than max_tokens, Default is 1024.
	LmGenReasoningMaxTokens int `json:"lm_gen_reasoning_max_tokens" default:"1024"`

	// LmGenMaxNewTokens Maximum number of new tokens to generate in a single completion. Default is 1024.
	LmGenMaxNewTokens int `json:"lm_gen_max_new_tokens" default:"1024"`

	// LmGenMinNewTokens Minimum number of new tokens to generate in a single completion. Default is 1.
	LmGenMinNewTokens int `json:"lm_gen_min_new_tokens" default:"1"`

	// LmGenDoSample Whether to use sampling; set this to False for deterministic outputs. Default is False.
	LmGenDoSample bool `json:"lm_gen_do_sample" default:"true"`

	// LmGenTemperature Controls the randomness of the output. Set to 0.0 for deterministic (repeatable) outputs. Default is 0.1.
	LmGenTemperature float64 `json:"lm_gen_temperature" default:"0.6"`

	// LmGenTopK Changing the top - k parameter sets the size of the shortlist the model samples from as it outputs each token.
	// Setting top - k to 1 gives us greedy decoding. Default is 1
	LmGenTopK int `json:"lm_gen_top_k" default:"10"`

	// LmGenTopP Top-p is usually set to a high value (like 0.75) with the purpose of limiting the long tail of low-probability tokens that may be sampled.
	// We can use both top-k and top-p together. If both k and p are enabled, p acts after k. Default is 0.8.
	LmGenTopP float64 `json:"lm_gen_top_p" default:"0.9"`

	// LmGenMinP samples from tokens with probability larger than min_p * highest_token_probability. Default is 0.0.
	LmGenMinP float64 `json:"lm_gen_min_p" default:"0.0"`

	// LmGenRepetitionPenalty Controls the token repetition pealty.  no repetition Default is 1.0, >1.0: low repetition, <1.0: high repetition
	LmGenRepetitionPenalty float64 `json:"lm_gen_repetition_penalty" default:"1.0"`

	// LmGenStops A list of strings that will stop the generation. Default is []. If the stop word is a substring of the generated text, the generation will stop.
	LmGenStops []string `json:"lm_gen_stops"`

	// LmGenStopIds A list of token ids that will stop the generation. Default is []. If the stop id is a substring token id of the generated text, the generation will stop.
	LmGenStopIds []int `json:"lm_gen_stop_ids"`

	// LmGenEndId The end token id. Default is 0. If the end id is a substring token id of the generated text, the generation will stop.
	LmGenEndId int `json:"lm_gen_end_id" default:"0"`

	// LmGenPadId The pad token id. Default is 0. If the pad id is a substring token id of the generated text, the generation will stop.
	LmGenPadId int `json:"lm_gen_pad_id" default:"0"`

	// LmGenThinking if use think, Whether to output
	// content; set this to False for deterministic outputs. Default is None. auto thinking
	LmGenThinking *bool `json:"lm_gen_thinking"`

	// LmGenThinkOutput if use RL model, Whether to output
	// content; set this to False for deterministic outputs. Default is True.
	LmGenThinkOutput bool `json:"lm_gen_think_output" default:"true"`

	// LmGenThinkIntervalTime The think interval time to tip user. Default is 0<=. no tip,
	LmGenThinkIntervalTime int `json:"lm_gen_think_interval_time" default:"0"`
}

// NewLMGenerateArgs creates a new LMGenerateArgs with default values
func NewLMGenerateArgs() *LMGenerateArgs {
	return &LMGenerateArgs{
		LmGenSeed:               42,
		LmMaxLength:             2048,
		LmGenMaxTokens:          2048,
		LmGenReasoningMaxTokens: 1024,
		LmGenMaxNewTokens:       1024,
		LmGenMinNewTokens:       1,
		LmGenDoSample:           true,
		LmGenTemperature:        0.6,
		LmGenTopK:               10,
		LmGenTopP:               0.9,
		LmGenMinP:               0.0,
		LmGenRepetitionPenalty:  1.0,
		LmGenStops:              []string{},
		LmGenStopIds:            []int{},
		LmGenEndId:              0,
		LmGenPadId:              0,
		LmGenThinkOutput:        true,
		LmGenThinkIntervalTime:  0,
	}
}
