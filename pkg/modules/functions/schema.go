package functions

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/ollama/ollama/api"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/shared"
)

// 标准 openai 定义的 tool schema
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

var OllamaAPISearchToolSchema = map[string]any{
	"type": "function",
	"function": map[string]any{
		"name":        "web_search",
		"description": "web search by query",
		"parameters": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"query": map[string]any{
					"type":        []string{"string"},
					"description": "web search query",
				},
			},
			"required": []string{"query"},
		},
	},
}

func AdapteOllamaToolSchema(schemas []map[string]any) (api.Tools, error) {
	tools := api.Tools{}
	err := mapstructure.Decode(schemas, &tools)
	if err != nil {
		return nil, err
	}
	return tools, nil
}

func AdapteOpenAIToolSchema(schemas []map[string]any) ([]openai.ChatCompletionToolUnionParam, error) {
	tools := []openai.ChatCompletionToolUnionParam{}
	for _, schema := range schemas {
		functionData := schema["function"].(map[string]any)

		// Extract function properties
		name := functionData["name"].(string)
		description := functionData["description"].(string)

		// Handle parameters
		var parameters shared.FunctionParameters
		if params, ok := functionData["parameters"].(map[string]any); ok {
			paramsMap := shared.FunctionParameters(params)
			parameters = paramsMap
		}

		// Create the function definition
		functionDef := shared.FunctionDefinitionParam{
			Name:        name,
			Description: param.Opt[string]{Value: description},
			Parameters:  parameters,
		}

		// Create the tool using the helper function
		tool := openai.ChatCompletionFunctionTool(functionDef)
		tools = append(tools, tool)
	}

	return tools, nil
}
