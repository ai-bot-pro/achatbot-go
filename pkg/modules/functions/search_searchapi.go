package functions

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"achatbot/pkg/params"

	"github.com/go-viper/mapstructure/v2"
)

const (
	SearchApiBaseUrl = "https://www.searchapi.io/api/v1/search"
	SearchApiTag     = "search_api"
)

type SearchApi struct {
	args params.SearchApiArgs
}

// NewSearchApi creates a new Search API instance
func NewSearchApi(args params.SearchApiArgs) *SearchApi {
	return &SearchApi{
		args: args,
	}
}

func (s *SearchApi) GetToolCall() map[string]any {
	return SearchToolSchema
}

func (s *SearchApi) GetOllamaAPIToolCall() map[string]any {
	return OllamaAPISearchToolSchema
}

func (s *SearchApi) Execute(args map[string]any) (string, error) {
	err := mapstructure.Decode(args, &s.args)
	if err != nil {
		return "", err
	}
	return s.WebSearch(s.args.Query)
}

// WebSearch implements the IWebSearch interface
func (s *SearchApi) WebSearch(query string) (string, error) {
	apiKey := os.Getenv("SEARCH_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("SEARCH_API_KEY environment variable not set")
	}

	// Build query parameters
	params := map[string]string{
		"engine":  s.args.Engine,
		"api_key": apiKey,
		"q":       query,
		"gl":      s.args.GL,
		"hl":      s.args.HL,
		"page":    fmt.Sprintf("%d", s.args.Page),
		"num":     fmt.Sprintf("%d", s.args.Num),
	}

	// Construct URL with query parameters
	url := SearchApiBaseUrl + "?" + buildQueryString(params)

	// Create HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	// Extract snippets from organic results
	organicResults, ok := result["organic_results"].([]any)
	if !ok {
		return "", fmt.Errorf("invalid response format: 'organic_results' field not found or not an array")
	}

	snippets := make([]string, 0)
	for _, item := range organicResults {
		if itemMap, ok := item.(map[string]any); ok {
			if snippet, ok := itemMap["snippet"].(string); ok {
				snippets = append(snippets, snippet)
			}
		}
	}

	// Convert snippets to JSON
	snippetsJson, err := json.Marshal(snippets)
	if err != nil {
		return "", fmt.Errorf("failed to marshal snippets: %v", err)
	}

	return string(snippetsJson), nil
}

// buildQueryString constructs a query string from a map of parameters
func buildQueryString(params map[string]string) string {
	values := url.Values{}
	for key, value := range params {
		if value != "" {
			values.Add(key, value)
		}
	}
	return values.Encode()
}
