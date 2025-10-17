package main

import (
	"fmt"
	"log"
	"os"

	"achatbot/pkg/modules/functions"
	"achatbot/pkg/params"
)

// ExampleSearchAPI demonstrates how to use the SearchAPI
func ExampleSearchAPI() {
	// Set your SEARCH_API_KEY in the environment variables
	// os.Setenv("SEARCH_API_KEY", "your_searchapi_key_here")

	// Check if the API key is set
	apiKey := os.Getenv("SEARCH_API_KEY")
	if apiKey == "" {
		log.Println("Warning: SEARCH_API_KEY environment variable not set")
	}

	// Create SearchApiArgs with search parameters
	args := params.SearchApiArgs{
		Engine: "google", // or "bing", "duckduckgo", etc.
		GL:     "us",     // Country code
		HL:     "en",     // Language code
		Page:   1,        // Page number
		Num:    10,       // Number of results
	}

	// Create a new SearchApi instance
	searchApi := functions.NewSearchApi(args)

	// Example of getting tool call definition
	toolCall := searchApi.GetToolCall()
	fmt.Printf("Tool Call Definition: %+v\n\n", toolCall)

	// Example search query
	query := "latest AI news"

	// Execute the search (this will fail if no API key is set)
	result, err := searchApi.WebSearch(query)
	if err != nil {
		log.Printf("Error performing search: %v", err)
		return
	}

	fmt.Printf("Search Results for '%s':\n%s\n", query, result)
}

