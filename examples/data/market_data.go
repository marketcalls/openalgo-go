package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/marketcalls/openalgo-go/openalgo"
)

func main() {
	// Initialize OpenAlgo client with provided API key
	client := openalgo.NewClient(
		"openalgo-api-key",      // API Key
		"http://127.0.0.1:5000", // Host URL
		"v1",                    // API Version
		"ws://127.0.0.1:8765",   // WebSocket URL
	)

	fmt.Println("=== Testing Quote, Depth, and History APIs (RAW OUTPUT) ===\n")

	// 1. Test QUOTES API
	fmt.Println("1. QUOTES API - RELIANCE")
	fmt.Println("====================================================")

	quotesResp, err := client.Quotes("RELIANCE", "NSE")
	if err != nil {
		log.Printf("❌ Error fetching quotes: %v", err)
	} else {
		// Pretty print raw JSON response
		jsonData, _ := json.MarshalIndent(quotesResp, "", "  ")
		fmt.Printf("Raw Response:\n%s\n", string(jsonData))
	}

	// 2. Test DEPTH API (Market Depth / Order Book)
	fmt.Println("\n2. DEPTH API - SBIN")
	fmt.Println("====================================================")

	depthResp, err := client.Depth("SBIN", "NSE")
	if err != nil {
		log.Printf("❌ Error fetching depth: %v", err)
	} else {
		// Pretty print raw JSON response
		jsonData, _ := json.MarshalIndent(depthResp, "", "  ")
		fmt.Printf("Raw Response:\n%s\n", string(jsonData))
	}

	// 3. Test HISTORY API - Intraday data
	fmt.Println("\n3. HISTORY API - INFY (5-minute candles)")
	fmt.Println("====================================================")

	historyResp, err := client.History(
		"INFY",       // symbol
		"NSE",        // exchange
		"5m",         // interval (5 minutes)
		"2025-09-20", // start_date (2 days ago)
		"2025-09-22", // end_date (today)
	)
	if err != nil {
		log.Printf("❌ Error fetching history: %v", err)
	} else {
		// Pretty print raw JSON response
		jsonData, _ := json.MarshalIndent(historyResp, "", "  ")
		fmt.Printf("Raw Response:\n%s\n", string(jsonData))
	}

	// 4. Test HISTORY API - Daily data
	fmt.Println("\n4. HISTORY API - TCS (Daily candles)")
	fmt.Println("====================================================")

	dailyHistoryResp, err := client.History(
		"TCS",        // symbol
		"NSE",        // exchange
		"1d",         // interval (daily)
		"2025-09-15", // start_date (1 week ago)
		"2025-09-22", // end_date (today)
	)
	if err != nil {
		log.Printf("❌ Error fetching daily history: %v", err)
	} else {
		// Pretty print raw JSON response
		jsonData, _ := json.MarshalIndent(dailyHistoryResp, "", "  ")
		fmt.Printf("Raw Response:\n%s\n", string(jsonData))
	}

	// 5. Test multiple symbols for quotes
	fmt.Println("\n5. QUOTES API - Multiple Symbols")
	fmt.Println("====================================================")

	symbols := []struct {
		symbol   string
		exchange string
	}{
		{"NHPC", "NSE"},
		{"YESBANK", "NSE"},
		{"IDEA", "NSE"},
	}

	for _, s := range symbols {
		fmt.Printf("\n--- %s ---\n", s.symbol)
		quotes, err := client.Quotes(s.symbol, s.exchange)
		if err != nil {
			log.Printf("❌ Error: %v", err)
		} else {
			jsonData, _ := json.MarshalIndent(quotes, "", "  ")
			fmt.Printf("%s\n", string(jsonData))
		}
	}

	// 6. Test INTERVALS API to see available intervals
	fmt.Println("\n6. INTERVALS API - Available Time Intervals")
	fmt.Println("====================================================")

	intervalsResp, err := client.Intervals()
	if err != nil {
		log.Printf("❌ Error fetching intervals: %v", err)
	} else {
		jsonData, _ := json.MarshalIndent(intervalsResp, "", "  ")
		fmt.Printf("Raw Response:\n%s\n", string(jsonData))
	}

	fmt.Println("\n=== Market Data Test Complete ===")
}
