package functions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"achatbot/pkg/params"
	"io"

	"github.com/go-viper/mapstructure/v2"
)

const (
	SerperApiBaseUrl = "https://google.serper.dev/search"
	SerperApiTag     = "serper_api"
)

type SerperApi struct {
	args params.SerperApiArgs
}

// NewSerperApi creates a new Serper API instance
func NewSerperApi(args params.SerperApiArgs) *SerperApi {
	return &SerperApi{
		args: args,
	}
}

func (s *SerperApi) GetToolCall() map[string]any {
	return SearchToolSchema
}

func (s *SerperApi) GetOllamaAPIToolCall() map[string]any {
	return OllamaAPISearchToolSchema
}

func (s *SerperApi) Execute(args map[string]any) (string, error) {
	err := mapstructure.Decode(args, &s.args)
	if err != nil {
		return "", err
	}
	return s.WebSearch(s.args.Query)
}

// WebSearch implements the SearchBaseApi interface
func (s *SerperApi) WebSearch(query string) (string, error) {
	apiKey := os.Getenv("SERPER_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("SERPER_API_KEY environment variable not set")
	}

	url := SerperApiBaseUrl
	payload := map[string]any{
		"q":    query,
		"gl":   s.args.GL,
		"hl":   s.args.HL,
		"page": s.args.Page,
		"num":  s.args.Num,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("X-API-KEY", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	// Extract snippets from organic results
	organics, ok := result["organic"].([]any)
	if !ok {
		return "", fmt.Errorf("invalid response format: 'organics' field not found or not an array")
	}

	snippets := make([]string, 0)
	for _, item := range organics {
		if itemMap, ok := item.(map[string]any); ok {
			if snippet, ok := itemMap["snippet"].(string); ok {
				snippets = append(snippets, snippet)
			}
		}
	}

	snippetsJson, err := json.Marshal(snippets)
	if err != nil {
		return "", fmt.Errorf("failed to marshal snippets: %v", err)
	}

	return string(snippetsJson), nil
}
