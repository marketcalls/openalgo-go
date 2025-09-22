package main

import (
	"fmt"
	"log"
	"time"

	"github.com/marketcalls/openalgo-go/openalgo"
)

func main() {
	// Initialize OpenAlgo client
	// Replace 'your_api_key_here' with your actual API key
	// Specify the host URL with your hosted domain or ngrok domain.
	// If running locally in windows then use the default host value.
	client := openalgo.NewClient(
		"your_api_key_here",
		"http://127.0.0.1:5000",
		"ws://127.0.0.1:8765", // Optional WebSocket URL
	)

	// Check OpenAlgo Version
	fmt.Printf("OpenAlgo Go SDK Version: %s\n", openalgo.Version)

	// PlaceOrder example - Market Order
	marketOrderResp, err := client.PlaceOrder(
		"Python",    // strategy
		"NHPC",      // symbol
		"BUY",       // action
		"NSE",       // exchange
		"MARKET",    // price_type
		"MIS",       // product
		1,           // quantity
	)
	if err != nil {
		log.Printf("Error placing market order: %v", err)
	} else {
		fmt.Printf("Market Order Response: %v\n", marketOrderResp)
	}

	// PlaceOrder example - Limit Order with optional parameters
	limitOrderResp, err := client.PlaceOrder(
		"Python",    // strategy
		"YESBANK",   // symbol
		"BUY",       // action
		"NSE",       // exchange
		"LIMIT",     // price_type
		"MIS",       // product
		"1",         // quantity
		map[string]interface{}{
			"price":               16.0,
			"trigger_price":       0.0,
			"disclosed_quantity":  "0",
		},
	)
	if err != nil {
		log.Printf("Error placing limit order: %v", err)
	} else {
		fmt.Printf("Limit Order Response: %v\n", limitOrderResp)
	}

	// PlaceSmartOrder example
	smartOrderResp, err := client.PlaceSmartOrder(
		"Python",      // strategy
		"TATAMOTORS",  // symbol
		"SELL",        // action
		"NSE",         // exchange
		"MARKET",      // price_type
		"MIS",         // product
		1,             // quantity
		5,             // position_size
	)
	if err != nil {
		log.Printf("Error placing smart order: %v", err)
	} else {
		fmt.Printf("Smart Order Response: %v\n", smartOrderResp)
	}

	// BasketOrder example
	basketOrders := []map[string]interface{}{
		{
			"symbol":    "BHEL",
			"exchange":  "NSE",
			"action":    "BUY",
			"quantity":  1,
			"pricetype": "MARKET",
			"product":   "MIS",
		},
		{
			"symbol":    "ZOMATO",
			"exchange":  "NSE",
			"action":    "SELL",
			"quantity":  1,
			"pricetype": "MARKET",
			"product":   "MIS",
		},
	}
	basketResp, err := client.BasketOrder("Python", basketOrders)
	if err != nil {
		log.Printf("Error placing basket order: %v", err)
	} else {
		fmt.Printf("Basket Order Response: %v\n", basketResp)
	}

	// SplitOrder example
	splitResp, err := client.SplitOrder(
		"Python",    // strategy
		"YESBANK",   // symbol
		"NSE",       // exchange
		"SELL",      // action
		105,         // quantity
		20,          // split_size
		"MARKET",    // price_type
		"MIS",       // product
	)
	if err != nil {
		log.Printf("Error placing split order: %v", err)
	} else {
		fmt.Printf("Split Order Response: %v\n", splitResp)
	}

	// ModifyOrder example
	modifyResp, err := client.ModifyOrder(
		"250408001002736", // order_id
		"Python",          // strategy
		"YESBANK",         // symbol
		"BUY",             // action
		"NSE",             // exchange
		"LIMIT",           // price_type
		"CNC",             // product
		1,                 // quantity
		"16.5",            // price (required)
		"0",               // disclosed_quantity (required)
		"0",               // trigger_price (required)
	)
	if err != nil {
		log.Printf("Error modifying order: %v", err)
	} else {
		fmt.Printf("Modify Order Response: %v\n", modifyResp)
	}

	// CancelOrder example
	cancelResp, err := client.CancelOrder(
		"250408001002736", // order_id
		"Python",          // strategy
	)
	if err != nil {
		log.Printf("Error canceling order: %v", err)
	} else {
		fmt.Printf("Cancel Order Response: %v\n", cancelResp)
	}

	// CancelAllOrder example
	cancelAllResp, err := client.CancelAllOrder("Python")
	if err != nil {
		log.Printf("Error canceling all orders: %v", err)
	} else {
		fmt.Printf("Cancel All Orders Response: %v\n", cancelAllResp)
	}

	// ClosePosition example
	closeResp, err := client.ClosePosition("Python")
	if err != nil {
		log.Printf("Error closing positions: %v", err)
	} else {
		fmt.Printf("Close Position Response: %v\n", closeResp)
	}

	// OrderStatus example
	statusResp, err := client.OrderStatus(
		"250408000989443", // order_id
		"Test Strategy",   // strategy
	)
	if err != nil {
		log.Printf("Error getting order status: %v", err)
	} else {
		fmt.Printf("Order Status Response: %v\n", statusResp)
	}

	// OpenPosition example
	posResp, err := client.OpenPosition(
		"Test Strategy", // strategy
		"YESBANK",       // symbol
		"NSE",           // exchange
		"MIS",           // product
	)
	if err != nil {
		log.Printf("Error getting open position: %v", err)
	} else {
		fmt.Printf("Open Position Response: %v\n", posResp)
	}

	// Quotes example
	quotesResp, err := client.Quotes("RELIANCE", "NSE")
	if err != nil {
		log.Printf("Error getting quotes: %v", err)
	} else {
		fmt.Printf("Quotes Response: %v\n", quotesResp)
	}

	// Depth example
	depthResp, err := client.Depth("SBIN", "NSE")
	if err != nil {
		log.Printf("Error getting depth: %v", err)
	} else {
		fmt.Printf("Depth Response: %v\n", depthResp)
	}

	// History example
	historyResp, err := client.History(
		"SBIN",       // symbol
		"NSE",        // exchange
		"5m",         // interval
		"2025-04-01", // start_date
		"2025-04-08", // end_date
	)
	if err != nil {
		log.Printf("Error getting history: %v", err)
	} else {
		fmt.Printf("History Response: %v\n", historyResp)
	}

	// Intervals example
	intervalsResp, err := client.Intervals()
	if err != nil {
		log.Printf("Error getting intervals: %v", err)
	} else {
		fmt.Printf("Intervals Response: %v\n", intervalsResp)
	}

	// Symbol example
	symbolResp, err := client.Symbol("RELIANCE", "NSE")
	if err != nil {
		log.Printf("Error getting symbol info: %v", err)
	} else {
		fmt.Printf("Symbol Response: %v\n", symbolResp)
	}

	// Search example
	searchResp, err := client.Search("NIFTY 25000 JUL CE", "NFO")
	if err != nil {
		log.Printf("Error searching: %v", err)
	} else {
		fmt.Printf("Search Response: %v\n", searchResp)
	}

	// Expiry example
	expiryResp, err := client.Expiry(
		"NIFTY",   // symbol
		"NFO",     // exchange
		"options", // instrument_type
	)
	if err != nil {
		log.Printf("Error getting expiry: %v", err)
	} else {
		fmt.Printf("Expiry Response: %v\n", expiryResp)
	}

	// Funds example
	fundsResp, err := client.Funds()
	if err != nil {
		log.Printf("Error getting funds: %v", err)
	} else {
		fmt.Printf("Funds Response: %v\n", fundsResp)
	}

	// OrderBook example
	orderBookResp, err := client.OrderBook()
	if err != nil {
		log.Printf("Error getting order book: %v", err)
	} else {
		fmt.Printf("OrderBook Response: %v\n", orderBookResp)
	}

	// TradeBook example
	tradeBookResp, err := client.TradeBook()
	if err != nil {
		log.Printf("Error getting trade book: %v", err)
	} else {
		fmt.Printf("TradeBook Response: %v\n", tradeBookResp)
	}

	// PositionBook example
	positionBookResp, err := client.PositionBook()
	if err != nil {
		log.Printf("Error getting position book: %v", err)
	} else {
		fmt.Printf("PositionBook Response: %v\n", positionBookResp)
	}

	// Holdings example
	holdingsResp, err := client.Holdings()
	if err != nil {
		log.Printf("Error getting holdings: %v", err)
	} else {
		fmt.Printf("Holdings Response: %v\n", holdingsResp)
	}

	// Analyzer Status example
	analyzerStatusResp, err := client.AnalyzerStatus()
	if err != nil {
		log.Printf("Error getting analyzer status: %v", err)
	} else {
		fmt.Printf("Analyzer Status Response: %v\n", analyzerStatusResp)
	}

	// Analyzer Toggle example - switch to analyze mode
	analyzerToggleResp, err := client.AnalyzerToggle(true)
	if err != nil {
		log.Printf("Error toggling analyzer: %v", err)
	} else {
		fmt.Printf("Analyzer Toggle Response: %v\n", analyzerToggleResp)
	}

	// WebSocket Streaming Examples
	// Initialize client with WebSocket URL
	wsClient := openalgo.NewClient(
		"your_api_key",
		"http://127.0.0.1:5000",
		"ws://127.0.0.1:8765",
	)

	// Connect to WebSocket
	if err := wsClient.Connect(); err != nil {
		log.Printf("Error connecting to WebSocket: %v", err)
		return
	}
	defer wsClient.Disconnect()

	// Define instruments to subscribe
	instruments := []openalgo.Instrument{
		{Exchange: "NSE", Symbol: "RELIANCE"},
		{Exchange: "NSE", Symbol: "INFY"},
	}

	// Subscribe to LTP updates
	ltpCallback := func(data interface{}) {
		fmt.Println("LTP Update Received:")
		fmt.Println(data)
	}
	if err := wsClient.SubscribeLTP(instruments, ltpCallback); err != nil {
		log.Printf("Error subscribing to LTP: %v", err)
	}

	// Subscribe to Quote updates
	quoteCallback := func(data interface{}) {
		fmt.Println("Quote Update Received:")
		fmt.Println(data)
	}
	if err := wsClient.SubscribeQuote(instruments, quoteCallback); err != nil {
		log.Printf("Error subscribing to quotes: %v", err)
	}

	// Subscribe to Depth updates
	depthCallback := func(data interface{}) {
		fmt.Println("Market Depth Update Received:")
		fmt.Println(data)
	}
	if err := wsClient.SubscribeDepth(instruments, depthCallback); err != nil {
		log.Printf("Error subscribing to depth: %v", err)
	}

	// Run for a few seconds to receive data
	time.Sleep(10 * time.Second)

	// Unsubscribe from streams
	wsClient.UnsubscribeLTP(instruments)
	wsClient.UnsubscribeQuote(instruments)
	wsClient.UnsubscribeDepth(instruments)
}