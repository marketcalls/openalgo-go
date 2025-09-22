package main

import (
	"fmt"
	"log"
	"strings"

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

	fmt.Println("=== OpenAlgo Basket & Split Order Examples ===\n")

	// =====================================
	// EXAMPLE 1: BASKET ORDER
	// =====================================
	fmt.Println("1. BASKET ORDER - Place Multiple Orders at Once")
	fmt.Println("=" + strings.Repeat("=", 50))

	// Define multiple orders to place simultaneously
	basketOrders := []map[string]interface{}{
		{
			"symbol":    "RELIANCE",
			"exchange":  "NSE",
			"action":    "BUY",
			"quantity":  5,
			"pricetype": "MARKET",
			"product":   "MIS",
		},
		{
			"symbol":    "INFY",
			"exchange":  "NSE",
			"action":    "BUY",
			"quantity":  10,
			"pricetype": "LIMIT",
			"product":   "CNC",
			"price":     "1500.00", // Limit price for INFY
		},
		{
			"symbol":    "TCS",
			"exchange":  "NSE",
			"action":    "SELL",
			"quantity":  3,
			"pricetype": "LIMIT",
			"product":   "MIS",
			"price":     "3100.00", // Limit price for TCS
		},
		{
			"symbol":        "HDFCBANK",
			"exchange":      "NSE",
			"action":        "BUY",
			"quantity":      2,
			"pricetype":     "SL", // Stop Loss order
			"product":       "MIS",
			"price":         "1450.00", // Limit price
			"trigger_price": "1448.00", // Trigger price
		},
	}

	fmt.Println("Placing basket order with 4 different orders:")
	fmt.Println("  1. RELIANCE - BUY 5 qty at MARKET")
	fmt.Println("  2. INFY - BUY 10 qty at LIMIT ₹1500")
	fmt.Println("  3. TCS - SELL 3 qty at LIMIT ₹3100")
	fmt.Println("  4. HDFC - BUY 2 qty STOP LOSS (trigger: ₹1448, limit: ₹1450)")

	// Place basket order
	basketResp, err := client.BasketOrder("GO Strategy", basketOrders)
	if err != nil {
		log.Printf("❌ Error placing basket order: %v", err)
	} else {
		if basketResp["status"] == "success" {
			fmt.Println("\n✅ Basket Order placed successfully!")

			// Display results for each order
			if results, ok := basketResp["results"].([]interface{}); ok {
				fmt.Println("\nOrder Results:")
				for i, result := range results {
					if r, ok := result.(map[string]interface{}); ok {
						fmt.Printf("  Order %d: Symbol=%v, Status=%v, OrderID=%v\n",
							i+1, r["symbol"], r["status"], r["orderid"])
					}
				}
			}
		} else {
			fmt.Printf("❌ Basket order failed: %v\n", basketResp)
		}
	}

	// =====================================
	// EXAMPLE 2: SPLIT ORDER
	// =====================================
	fmt.Println("\n2. SPLIT ORDER - Split Large Order into Smaller Chunks")
	fmt.Println("=" + strings.Repeat("=", 50))

	// Split order parameters
	totalQuantity := 100
	splitSize := 20

	fmt.Printf("Splitting order: ICICIBANK - BUY %d qty into chunks of %d\n",
		totalQuantity, splitSize)
	fmt.Printf("This will create %d orders of %d qty each\n",
		totalQuantity/splitSize, splitSize)

	// Place split order - Market order
	splitResp, err := client.SplitOrder(
		"GO Strategy", // strategy
		"ICICIBANK",   // symbol
		"NSE",         // exchange
		"BUY",         // action
		totalQuantity, // total quantity to split
		splitSize,     // size of each split
		"MARKET",      // price_type
		"MIS",         // product
	)
	if err != nil {
		log.Printf("❌ Error placing split order: %v", err)
	} else {
		if splitResp["status"] == "success" {
			fmt.Println("\n✅ Split Order placed successfully!")
			fmt.Printf("Total Quantity: %v\n", splitResp["total_quantity"])
			fmt.Printf("Split Size: %v\n", splitResp["split_size"])

			// Display each split order result
			if results, ok := splitResp["results"].([]interface{}); ok {
				fmt.Println("\nSplit Order Results:")
				for _, result := range results {
					if r, ok := result.(map[string]interface{}); ok {
						fmt.Printf("  Order #%v: Qty=%v, OrderID=%v, Status=%v\n",
							r["order_num"], r["quantity"], r["orderid"], r["status"])
					}
				}
			}
		} else {
			fmt.Printf("❌ Split order failed: %v\n", splitResp)
		}
	}

	// =====================================
	// EXAMPLE 3: SPLIT ORDER WITH LIMIT PRICE
	// =====================================
	fmt.Println("\n3. SPLIT ORDER - With Limit Price")
	fmt.Println("=" + strings.Repeat("=", 50))

	fmt.Println("Splitting LIMIT order: SBIN - SELL 50 qty at ₹450 into chunks of 10")

	// Place split order with limit price
	splitLimitResp, err := client.SplitOrder(
		"GO Strategy", // strategy
		"SBIN",   // symbol
		"NSE",    // exchange
		"SELL",   // action
		50,       // total quantity
		10,       // split size (5 orders of 10 each)
		"LIMIT",  // price_type
		"CNC",    // product
		map[string]interface{}{
			"price": "450.00", // Limit price
		},
	)
	if err != nil {
		log.Printf("❌ Error placing split limit order: %v", err)
	} else {
		if splitLimitResp["status"] == "success" {
			fmt.Println("\n✅ Split Limit Order placed successfully!")
			fmt.Printf("Total Quantity: %v\n", splitLimitResp["total_quantity"])
			fmt.Printf("Split Size: %v\n", splitLimitResp["split_size"])

			if results, ok := splitLimitResp["results"].([]interface{}); ok {
				fmt.Println("\nSplit Order Results:")
				for _, result := range results {
					if r, ok := result.(map[string]interface{}); ok {
						fmt.Printf("  Order #%v: Qty=%v at ₹450, OrderID=%v, Status=%v\n",
							r["order_num"], r["quantity"], r["orderid"], r["status"])
					}
				}
			}
		} else {
			fmt.Printf("❌ Split limit order failed: %v\n", splitLimitResp)
		}
	}

	// =====================================
	// EXAMPLE 4: ADVANCED BASKET ORDER
	// =====================================
	fmt.Println("\n4. ADVANCED BASKET ORDER - Mixed Order Types")
	fmt.Println("=" + strings.Repeat("=", 50))

	// Create a basket with different product types and order types
	advancedBasket := []map[string]interface{}{
		// Intraday Market Order
		{
			"symbol":    "NIFTY",
			"exchange":  "NSE",
			"action":    "BUY",
			"quantity":  25,
			"pricetype": "MARKET",
			"product":   "MIS", // Intraday
		},
		// Delivery Limit Order
		{
			"symbol":    "WIPRO",
			"exchange":  "NSE",
			"action":    "BUY",
			"quantity":  20,
			"pricetype": "LIMIT",
			"product":   "CNC", // Delivery
			"price":     "400.50",
		},
		// Stop Loss Market Order
		{
			"symbol":        "TATAMOTORS",
			"exchange":      "NSE",
			"action":        "SELL",
			"quantity":      15,
			"pricetype":     "SL-M", // Stop Loss Market
			"product":       "MIS",
			"trigger_price": "950.00",
		},
	}

	fmt.Println("Placing advanced basket order:")
	fmt.Println("  1. NIFTY - BUY 25 qty MARKET (MIS)")
	fmt.Println("  2. WIPRO - BUY 20 qty LIMIT ₹400.50 (CNC)")
	fmt.Println("  3. TATAMOTORS - SELL 15 qty SL-M trigger ₹950 (MIS)")

	advBasketResp, err := client.BasketOrder("GO Strategy", advancedBasket)
	if err != nil {
		log.Printf("❌ Error placing advanced basket order: %v", err)
	} else {
		if advBasketResp["status"] == "success" {
			fmt.Println("\n✅ Advanced Basket Order placed successfully!")

			if results, ok := advBasketResp["results"].([]interface{}); ok {
				fmt.Println("\nOrder Results:")
				for i, result := range results {
					if r, ok := result.(map[string]interface{}); ok {
						status := "✅"
						if r["status"] != "success" {
							status = "❌"
						}
						fmt.Printf("  %s Order %d: %v - OrderID=%v\n",
							status, i+1, r["symbol"], r["orderid"])
					}
				}
			}
		} else {
			fmt.Printf("❌ Advanced basket order failed: %v\n", advBasketResp)
		}
	}

	fmt.Println("\n=== Examples Complete ===")
	fmt.Println("\nKey Points:")
	fmt.Println("• Basket Orders: Execute multiple orders simultaneously")
	fmt.Println("• Split Orders: Break large orders into smaller chunks")
	fmt.Println("• Both support all order types (MARKET, LIMIT, SL, SL-M)")
	fmt.Println("• Both support all products (MIS, CNC, NRML)")
}
