package main

import (
	"fmt"
	"log"
	"os"

	"achatbot/pkg/modules/functions"
	"achatbot/pkg/params"
)

func main() {
	// Set your SERPER_API_KEY in the environment variables
	// os.Setenv("SERPER_API_KEY", "your-api-key")

	// Check if the API key is set
	apiKey := os.Getenv("SERPER_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set the SERPER_API_KEY environment variable")
	}

	// Create search parameters
	searchArgs := params.SerperApiArgs{
		GL:   "us", // Country code
		HL:   "en", // Language code
		Page: 1,    // Page number
		Num:  10,   // Number of results
	}

	// Create a new Serper API instance
	serper := functions.NewSerperApi(searchArgs)

	// Perform a web search
	query := "latest AI news"
	results, err := serper.WebSearch(query)
	if err != nil {
		log.Fatalf("Error performing search: %v", err)
	}

	fmt.Printf("Search results for '%s':\n%s\n", query, results)
}
