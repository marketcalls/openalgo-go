package main

import (
	"fmt"
	"log"

	"github.com/marketcalls/openalgo-go/openalgo"
)

func main() {
	// Initialize OpenAlgo client
	client := openalgo.NewClient(
		"3fa0620b435555c2c73720872288b626be96103a02edcb95e8ad864ee160b343", // API Key
		"http://127.0.0.1:5000", // Host URL
		"v1",                    // API Version
		"ws://127.0.0.1:8765",   // WebSocket URL (optional)
	)

	// 1. Fetch Account Funds
	fmt.Println("=== 1. ACCOUNT FUNDS ===")
	funds, err := client.Funds()
	if err != nil {
		log.Printf("Error fetching funds: %v", err)
	} else {
		if fundsData, ok := funds["data"].(map[string]interface{}); ok {
			fmt.Printf("Available Cash: ₹%v\n", fundsData["availablecash"])
			fmt.Printf("Collateral: ₹%v\n", fundsData["collateral"])
			fmt.Printf("M2M Realized: ₹%v\n", fundsData["m2mrealized"])
			fmt.Printf("M2M Unrealized: ₹%v\n", fundsData["m2munrealized"])
			fmt.Printf("Utilised Debits: ₹%v\n", fundsData["utiliseddebits"])
		}
	}

	// 2. Place Market Order for NHPC
	fmt.Println("\n=== 2. PLACE MARKET ORDER (NHPC) ===")
	marketOrder, err := client.PlaceOrder(
		"Python", // strategy
		"NHPC",   // symbol
		"BUY",    // action
		"NSE",    // exchange
		"MARKET", // price_type
		"MIS",    // product
		1,        // quantity
	)
	if err != nil {
		log.Printf("Error placing market order: %v", err)
	} else {
		if marketOrder["status"] == "success" {
			fmt.Printf("✅ Market Order placed successfully\n")
			fmt.Printf("Order ID: %v\n", marketOrder["orderid"])
		}
	}

	// 3. Fetch Holdings
	fmt.Println("\n=== 3. HOLDINGS ===")
	holdings, err := client.Holdings()
	if err != nil {
		log.Printf("Error fetching holdings: %v", err)
	} else {
		if holdData, ok := holdings["data"].(map[string]interface{}); ok {
			if holdingsList, ok := holdData["holdings"].([]interface{}); ok {
				fmt.Printf("Total Holdings: %d\n", len(holdingsList))
				fmt.Println("------------------------")
				for i, holding := range holdingsList {
					if h, ok := holding.(map[string]interface{}); ok {
						fmt.Printf("%d. %v - Qty: %v\n", i+1, h["symbol"], h["quantity"])
					}
				}
			}
		}
	}

	// 4. Fetch Order Book
	fmt.Println("\n=== 4. ORDER BOOK ===")
	orderBook, err := client.OrderBook()
	if err != nil {
		log.Printf("Error fetching order book: %v", err)
	} else {
		if orderData, ok := orderBook["data"].(map[string]interface{}); ok {
			if orders, ok := orderData["orders"].([]interface{}); ok {
				fmt.Printf("Total Orders: %d\n", len(orders))
				if len(orders) > 0 {
					fmt.Println("Recent Orders:")
					fmt.Println("------------------------")
					// Show last 5 orders
					start := 0
					if len(orders) > 5 {
						start = len(orders) - 5
					}
					for i := start; i < len(orders); i++ {
						if order, ok := orders[i].(map[string]interface{}); ok {
							fmt.Printf("%d. %v | %v | %v | Qty: %v | Status: %v | OrderID: %v\n",
								i-start+1,
								order["symbol"],
								order["action"],
								order["pricetype"],
								order["quantity"],
								order["order_status"],
								order["orderid"])
						}
					}
				}
			}
		}
	}

	// 5. Fetch Trade Book
	fmt.Println("\n=== 5. TRADE BOOK ===")
	tradeBook, err := client.TradeBook()
	if err != nil {
		log.Printf("Error fetching trade book: %v", err)
	} else {
		if trades, ok := tradeBook["data"].([]interface{}); ok {
			fmt.Printf("Total Trades: %d\n", len(trades))
			if len(trades) > 0 {
				fmt.Println("Recent Trades:")
				fmt.Println("------------------------")
				// Show last 5 trades
				start := 0
				if len(trades) > 5 {
					start = len(trades) - 5
				}
				for i := start; i < len(trades); i++ {
					if trade, ok := trades[i].(map[string]interface{}); ok {
						fmt.Printf("%d. %v | %v | Qty: %v | Price: %v | OrderID: %v\n",
							i-start+1,
							trade["symbol"],
							trade["action"],
							trade["quantity"],
							trade["average_price"],
							trade["orderid"])
					}
				}
			}
		}
	}

	// 6. Fetch Position Book
	fmt.Println("\n=== 6. POSITION BOOK ===")
	positionBook, err := client.PositionBook()
	if err != nil {
		log.Printf("Error fetching position book: %v", err)
	} else {
		if posData, ok := positionBook["data"].(map[string]interface{}); ok {
			// Check for net positions
			if netPos, ok := posData["net"].([]interface{}); ok && len(netPos) > 0 {
				fmt.Printf("Net Positions: %d\n", len(netPos))
				fmt.Println("------------------------")
				for i, pos := range netPos {
					if position, ok := pos.(map[string]interface{}); ok {
						fmt.Printf("%d. %v | Qty: %v | Product: %v | Avg Price: %v\n",
							i+1,
							position["symbol"],
							position["quantity"],
							position["product"],
							position["average_price"])
					}
				}
			} else {
				fmt.Println("No open positions")
			}
		}
	}

	fmt.Println("\n=== COMPLETE ===")
}