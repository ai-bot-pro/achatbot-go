package main

import (
	"fmt"
	"log"
	"os"

	"achatbot/pkg/modules/functions"
	"achatbot/pkg/params"
)

func ExampleSerperAPI() {
	// Set your SERPER_API_KEY in the environment variables
	// os.Setenv("SERPER_API_KEY", "your-api-key")

	// Check if the API key is set
	apiKey := os.Getenv("SERPER_API_KEY")
	if apiKey == "" {
		log.Println("Note: SERPER_API_KEY environment variable not set. Set it to run actual searches.")
	}

	// Also check for SEARCH_API_KEY
	searchApiKey := os.Getenv("SEARCH_API_KEY")
	if searchApiKey == "" {
		log.Println("Note: SEARCH_API_KEY environment variable not set. Set it to run actual searches with SearchAPI.")
	}

	// Create search parameters for Serper
	serperArgs := params.SerperApiArgs{
		GL:   "us", // Country code
		HL:   "en", // Language code
		Page: 1,    // Page number
		Num:  10,   // Number of results
	}

	// Create a new Serper API instance
	_ = functions.NewSerperApi(serperArgs)

	// Create search parameters for SearchAPI
	searchApiArgs := params.SearchApiArgs{
		Engine: "google", // Search engine
		GL:     "us",     // Country code
		HL:     "en",     // Language code
		Page:   1,        // Page number
		Num:    10,       // Number of results
	}

	// Create a new Search API instance
	_ = functions.NewSearchApi(searchApiArgs)

	// In a real application, you would pass a proper session
	// For this example, we'll just demonstrate the interface usage
	fmt.Println("Serper API instance created successfully")
	fmt.Printf("Serper parameters: GL=%s, HL=%s, Page=%d, Num=%d\n",
		serperArgs.GL, serperArgs.HL, serperArgs.Page, serperArgs.Num)

	fmt.Println("Search API instance created successfully")
	fmt.Printf("Search API parameters: Engine=%s, GL=%s, HL=%s, Page=%d, Num=%d\n",
		searchApiArgs.Engine, searchApiArgs.GL, searchApiArgs.HL, searchApiArgs.Page, searchApiArgs.Num)
}
