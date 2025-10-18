package functions

import (
	"achatbot/pkg/common"
	"achatbot/pkg/params"
	"os"
	"slices"
)

type RegisteredFunctions map[string]common.IFunction

var RegisterFuncs = NewRegisteredFunctions()

func init() {
	searchFunc := os.Getenv("SEARCH_FUNC")
	switch searchFunc {
	case "serperapi":
		RegisterFuncs.Register("web_search", NewSerperApi(params.SerperApiArgs{
			GL:   "us", // Country code
			HL:   "en", // Language code
			Page: 1,    // Page number
			Num:  5,    // Number of results
		}))
	case "searchapi":
		RegisterFuncs.Register("web_search", NewSearchApi(params.SearchApiArgs{
			Engine: "google", // Search engine
			GL:     "us",     // Country code
			HL:     "en",     // Language code
			Page:   1,        // Page number
			Num:    5,        // Number of results
		}))
	default:
		RegisterFuncs.Register("web_search", NewSerperApi(params.SerperApiArgs{
			GL:   "us", // Country code
			HL:   "en", // Language code
			Page: 1,    // Page number
			Num:  5,    // Number of results
		}))
	}
}

func NewRegisteredFunctions() *RegisteredFunctions {
	return &RegisteredFunctions{}
}

func (r *RegisteredFunctions) Register(name string, value common.IFunction) {
	(*r)[name] = value
}

func (r *RegisteredFunctions) Get(name string) common.IFunction {
	return (*r)[name]
}

func (r *RegisteredFunctions) GetToolCall(name string) map[string]any {
	return (*r)[name].GetToolCall()
}

func (r *RegisteredFunctions) GetToolCalls() []map[string]any {
	toolCalls := make([]map[string]any, 0)
	for name := range *r {
		toolCall := (*r)[name].GetToolCall()
		toolCalls = append(toolCalls, toolCall)
	}
	return toolCalls
}

func (r *RegisteredFunctions) GetToolCallsByName(names []string) []map[string]any {
	toolCalls := make([]map[string]any, 0)
	for name := range *r {
		if slices.Contains(names, name) {
			toolCall := (*r)[name].GetToolCall()
			toolCalls = append(toolCalls, toolCall)
		}
	}
	return toolCalls
}

func (r *RegisteredFunctions) GetOllamaAPIToolCall(name string) map[string]any {
	return (*r)[name].GetOllamaAPIToolCall()
}

func (r *RegisteredFunctions) GetOllamaAPIToolCalls() []map[string]any {
	toolCalls := make([]map[string]any, 0)
	for name := range *r {
		toolCall := (*r)[name].GetOllamaAPIToolCall()
		toolCalls = append(toolCalls, toolCall)
	}
	return toolCalls
}

func (r *RegisteredFunctions) GetOllamaAPIToolCallsByName(names []string) []map[string]any {
	toolCalls := make([]map[string]any, 0)
	for name := range *r {
		if slices.Contains(names, name) {
			toolCall := (*r)[name].GetOllamaAPIToolCall()
			toolCalls = append(toolCalls, toolCall)
		}
	}
	return toolCalls
}

func (r *RegisteredFunctions) Execute(name string, args map[string]any) (string, error) {
	return (*r)[name].Execute(args)
}
