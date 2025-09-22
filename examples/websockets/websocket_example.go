package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/marketcalls/openalgo-go/openalgo"
)

func main() {
	// Initialize OpenAlgo client
	client := openalgo.NewClient(
		"YOUR_API_KEY",          // Replace with your API key
		"http://127.0.0.1:5000", // OpenAlgo server URL
		"v1",                    // API version
		"ws://127.0.0.1:8765",   // WebSocket URL for live data
	)

	// Connect to WebSocket
	fmt.Println("=== Connecting to WebSocket ===")
	err := client.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	fmt.Println("✅ WebSocket connected successfully\n")

	// Define instruments to subscribe
	instruments := []openalgo.Instrument{
		{Symbol: "RELIANCE", Exchange: "NSE"},
		{Symbol: "INFY", Exchange: "NSE"},
		{Symbol: "TCS", Exchange: "NSE"},
	}

	// Subscribe to LTP with raw data printing
	fmt.Println("=== Subscribing to LTP Updates (RAW DATA) ===")
	err = client.SubscribeLTP(instruments, func(data interface{}) {
		// Print raw data as JSON
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Printf("[LTP RAW] Error: %v\n", err)
		} else {
			fmt.Printf("[LTP RAW DATA]:\n%s\n---\n", string(jsonData))
		}
	})
	if err != nil {
		log.Printf("Failed to subscribe to LTP: %v", err)
	} else {
		fmt.Println("✅ Subscribed to LTP updates\n")
	}

	// Subscribe to Quotes with raw data printing
	quoteInstruments := []openalgo.Instrument{
		{Symbol: "RELIANCE", Exchange: "NSE"},
	}
	fmt.Println("=== Subscribing to Quote Updates (RAW DATA) ===")
	err = client.SubscribeQuote(quoteInstruments, func(data interface{}) {
		// Print raw data as JSON
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Printf("[QUOTE RAW] Error: %v\n", err)
		} else {
			fmt.Printf("[QUOTE RAW DATA]:\n%s\n---\n", string(jsonData))
		}
	})
	if err != nil {
		log.Printf("Failed to subscribe to quotes: %v", err)
	} else {
		fmt.Println("✅ Subscribed to quote updates\n")
	}

	// Keep the connection alive
	fmt.Println("=== Streaming Live RAW Data ===")
	fmt.Println("Press Ctrl+C to stop...\n")

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signal
	<-sigChan

	// Cleanup
	fmt.Println("\n=== Cleaning up ===")
	client.UnsubscribeLTP(instruments)
	client.UnsubscribeQuote(quoteInstruments)
	client.Disconnect()
	fmt.Println("✅ Disconnected")
}