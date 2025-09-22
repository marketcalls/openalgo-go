package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/marketcalls/openalgo-go/openalgo"
)

func main() {
	// Initialize OpenAlgo client
	client := openalgo.NewClient(
		"YOUR_API_KEY",          // Replace with your API key
		"http://127.0.0.1:5000", // OpenAlgo server URL
		"v1",                    // API version
		"ws://127.0.0.1:8765",   // WebSocket URL (optional)
	)

	fmt.Println("=== OpenAlgo API Analyzer Examples ===\n")

	// =====================================
	// EXAMPLE 1: API ANALYZER STATUS
	// =====================================
	fmt.Println("1. API ANALYZER - Check Analyzer Status")
	fmt.Println("=" + strings.Repeat("=", 50))

	fmt.Println("Fetching API Analyzer status...")
	analyzerStatus, err := client.AnalyzerStatus()
	if err != nil {
		log.Printf("❌ Error fetching analyzer status: %v", err)
	} else {
		if analyzerStatus["status"] == "success" {
			fmt.Println("✅ Analyzer Status retrieved successfully!")
			fmt.Printf("Full Response: %+v\n", analyzerStatus)

			// Display analyzer data if available
			if analyzerData, ok := analyzerStatus["data"].(map[string]interface{}); ok {
				fmt.Println("\nAnalyzer Data:")
				for key, value := range analyzerData {
					fmt.Printf("  %s: %v\n", key, value)
				}
			}
		} else {
			fmt.Printf("⚠️ Analyzer status check failed: %v\n", analyzerStatus)
		}
	}

	// =====================================
	// EXAMPLE 2: ENABLE API ANALYZER
	// =====================================
	fmt.Println("\n2. API ANALYZER - Enable Analyzer")
	fmt.Println("=" + strings.Repeat("=", 50))

	fmt.Println("Enabling API Analyzer...")
	enableResp, err := client.AnalyzerToggle(true) // true to enable
	if err != nil {
		log.Printf("❌ Error enabling analyzer: %v", err)
	} else {
		if enableResp["status"] == "success" {
			fmt.Println("✅ Analyzer enabled successfully!")
			fmt.Printf("Response: %+v\n", enableResp)
		} else {
			fmt.Printf("⚠️ Failed to enable analyzer: %v\n", enableResp)
		}
	}

	// =====================================
	// EXAMPLE 3: GENERATE API METRICS
	// =====================================
	fmt.Println("\n3. GENERATE API METRICS")
	fmt.Println("=" + strings.Repeat("=", 50))

	fmt.Println("Making API calls to generate metrics...")

	// Make several API calls to generate analyzer data
	apiCalls := []struct {
		name string
		fn   func() (map[string]interface{}, error)
	}{
		{"Funds", client.Funds},
		{"OrderBook", client.OrderBook},
		{"TradeBook", client.TradeBook},
		{"Holdings", client.Holdings},
		{"PositionBook", client.PositionBook},
	}

	for _, call := range apiCalls {
		startTime := time.Now()
		resp, err := call.fn()
		latency := time.Since(startTime)

		if err != nil {
			fmt.Printf("  ❌ %s: Error - %v\n", call.name, err)
		} else if resp["status"] == "success" {
			fmt.Printf("  ✅ %s: Success (%.2fms)\n", call.name, latency.Seconds()*1000)
		} else {
			fmt.Printf("  ⚠️ %s: Failed - %v\n", call.name, resp)
		}

		// Small delay between calls
		time.Sleep(100 * time.Millisecond)
	}

	// =====================================
	// EXAMPLE 4: CHECK METRICS AFTER CALLS
	// =====================================
	fmt.Println("\n4. CHECK ANALYZER METRICS")
	fmt.Println("=" + strings.Repeat("=", 50))

	fmt.Println("Fetching analyzer metrics after API calls...")
	metricsStatus, err := client.AnalyzerStatus()
	if err != nil {
		log.Printf("❌ Error fetching metrics: %v", err)
	} else {
		if metricsStatus["status"] == "success" {
			fmt.Println("✅ Metrics retrieved successfully!")

			if data, ok := metricsStatus["data"].(map[string]interface{}); ok {
				fmt.Println("\nAnalyzer Metrics:")
				for key, value := range data {
					// Format the output based on the type of metric
					switch key {
					case "total_calls", "success_count", "error_count":
						fmt.Printf("  %s: %.0f\n", key, toFloat64(value))
					case "avg_response_time", "min_response_time", "max_response_time":
						fmt.Printf("  %s: %.2fms\n", key, toFloat64(value))
					case "success_rate":
						fmt.Printf("  %s: %.1f%%\n", key, toFloat64(value))
					default:
						fmt.Printf("  %s: %v\n", key, value)
					}
				}
			}
		}
	}

	// =====================================
	// EXAMPLE 5: DISABLE API ANALYZER
	// =====================================
	fmt.Println("\n5. API ANALYZER - Disable Analyzer")
	fmt.Println("=" + strings.Repeat("=", 50))

	fmt.Println("Disabling API Analyzer...")
	disableResp, err := client.AnalyzerToggle(false) // false to disable
	if err != nil {
		log.Printf("❌ Error disabling analyzer: %v", err)
	} else {
		if disableResp["status"] == "success" {
			fmt.Println("✅ Analyzer disabled successfully!")
			fmt.Printf("Response: %+v\n", disableResp)
		} else {
			fmt.Printf("⚠️ Failed to disable analyzer: %v\n", disableResp)
		}
	}

	// =====================================
	// EXAMPLE 6: HEALTH CHECK SIMULATION
	// =====================================
	fmt.Println("\n6. HEALTH CHECK - Test API Connectivity")
	fmt.Println("=" + strings.Repeat("=", 50))

	fmt.Println("Performing health check with 5 API calls...")

	successCount := 0
	var totalLatency time.Duration
	checkCount := 5

	for i := 1; i <= checkCount; i++ {
		startTime := time.Now()

		// Use Funds endpoint as a lightweight health check
		resp, err := client.Funds()
		latency := time.Since(startTime)

		if err != nil {
			fmt.Printf("  Check %d: ❌ Failed (%v)\n", i, err)
		} else if resp["status"] == "success" {
			fmt.Printf("  Check %d: ✅ Success (%.2fms)\n", i, latency.Seconds()*1000)
			successCount++
			totalLatency += latency
		} else {
			fmt.Printf("  Check %d: ⚠️ API Error: %v\n", i, resp["message"])
		}

		// Wait between checks
		if i < checkCount {
			time.Sleep(500 * time.Millisecond)
		}
	}

	// Display health check statistics
	fmt.Printf("\n📊 Health Check Statistics:\n")
	fmt.Printf("  Total checks: %d\n", checkCount)
	fmt.Printf("  Successful: %d\n", successCount)
	fmt.Printf("  Failed: %d\n", checkCount-successCount)
	fmt.Printf("  Success rate: %.1f%%\n", float64(successCount)*100/float64(checkCount))

	if successCount > 0 {
		avgLatency := totalLatency / time.Duration(successCount)
		fmt.Printf("  Average latency: %.2fms\n", avgLatency.Seconds()*1000)
	}

	// Determine API health
	if successCount == checkCount {
		fmt.Println("\n✅ API Health: EXCELLENT")
	} else if float64(successCount)/float64(checkCount) >= 0.8 {
		fmt.Println("\n⚠️ API Health: GOOD (some issues)")
	} else if successCount > 0 {
		fmt.Println("\n⚠️ API Health: POOR (significant issues)")
	} else {
		fmt.Println("\n❌ API Health: CRITICAL (API unreachable)")
	}

	// =====================================
	// EXAMPLE 7: MARKET DATA LATENCY TEST
	// =====================================
	fmt.Println("\n7. MARKET DATA LATENCY TEST")
	fmt.Println("=" + strings.Repeat("=", 50))

	symbols := []struct {
		symbol   string
		exchange string
	}{
		{"RELIANCE", "NSE"},
		{"INFY", "NSE"},
		{"TCS", "NSE"},
	}

	fmt.Println("Testing market data endpoint latency...")

	for _, s := range symbols {
		startTime := time.Now()
		quotes, err := client.Quotes(s.symbol, s.exchange)
		latency := time.Since(startTime)

		if err != nil {
			fmt.Printf("  %s: ❌ Error - %v\n", s.symbol, err)
		} else if quotes["status"] == "success" {
			fmt.Printf("  %s: ✅ Success (%.2fms)\n", s.symbol, latency.Seconds()*1000)

			// Display LTP if available
			if data, ok := quotes["data"].(map[string]interface{}); ok {
				if ltp, ok := data["ltp"]; ok {
					fmt.Printf("    LTP: ₹%v\n", ltp)
				}
			}
		}
	}

	fmt.Println("\n=== Examples Complete ===")
	fmt.Println("\nKey Features:")
	fmt.Println("• Analyzer Status: View current analyzer state")
	fmt.Println("• Analyzer Toggle: Enable/disable API performance tracking")
	fmt.Println("• Metrics Collection: Track API usage and response times")
	fmt.Println("• Health Check: Test API connectivity and latency")
	fmt.Println("• Market Data Latency: Test real-time data endpoints")
}

// Helper function to convert interface{} to float64
func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case int:
		return float64(val)
	case int64:
		return float64(val)
	default:
		return 0
	}
}