package functions

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/ollama/ollama/api"
	"github.com/openai/openai-go/v3"
)

var SearchToolSchema = map[string]any{
	"type": "function",
	"function": map[string]any{
		"name":        "web_search",
		"description": "web search by query",
		"parameters": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"query": map[string]any{
					"type":        "string",
					"description": "web search query",
				},
			},
			"required": []string{"query"},
		},
	},
}

func AdapteOllamaSearchToolSchema(schemas []map[string]any) (api.Tools, error) {
	tools := api.Tools{}
	for _, schema := range schemas {
		function := schema["function"].(map[string]any)
		parameters := function["parameters"].(map[string]any)
		properties := parameters["properties"].(map[string]any)
		query := properties["query"].(map[string]any)
		query["type"] = []string{"string"}

		tool := api.Tool{}
		err := mapstructure.Decode(SearchToolSchema, &tool)
		if err != nil {
			return nil, err
		}
		tools = append(tools, tool)
	}

	return tools, nil
}

func AdapteOpenAISearchToolSchema(schemas []map[string]any) ([]openai.ChatCompletionToolUnionParam, error) {
	tools := []openai.ChatCompletionToolUnionParam{}
	for _, schema := range schemas {
		function := schema["function"].(map[string]any)
		function["parameters"] = openai.String(function["parameters"].(string))

		tool := openai.ChatCompletionFunctionToolParam{}
		err := mapstructure.Decode(SearchToolSchema, &tool)
		if err != nil {
			return nil, err
		}
		tools = append(tools, openai.ChatCompletionToolUnionParam{OfFunction: &tool})
	}

	return tools, nil
}
