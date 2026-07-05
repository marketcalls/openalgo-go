package openalgo

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Instrument represents a trading instrument for WebSocket subscriptions
type Instrument struct {
	Exchange      string `json:"exchange"`
	Symbol        string `json:"symbol"`
	ExchangeToken string `json:"exchange_token,omitempty"`
}

// SubscriptionMessage represents the WebSocket subscription message format
type SubscriptionMessage struct {
	Action   string `json:"action"`
	Symbol   string `json:"symbol"`
	Exchange string `json:"exchange"`
	Mode     int    `json:"mode"`
	Depth    int    `json:"depth,omitempty"`
}

// AuthMessage represents the WebSocket authentication message
type AuthMessage struct {
	Action string `json:"action"`
	APIKey string `json:"api_key"`
}

// MarketData represents the market data received from WebSocket
type MarketData struct {
	Type      string      `json:"type"`
	Symbol    string      `json:"symbol"`
	Exchange  string      `json:"exchange"`
	Mode      int         `json:"mode"`
	Data      interface{} `json:"data"`
}

// ltpCacheEntry mirrors the Python SDK's ltp_data[...] snapshot: {'price': ..., 'timestamp': ...}
type ltpCacheEntry struct {
	Price     float64
	Timestamp int64
}

// quoteCacheEntry mirrors the Python SDK's quotes_data[...] snapshot.
type quoteCacheEntry struct {
	Open              float64
	High              float64
	Low               float64
	Close             float64
	LTP               float64
	Volume            float64
	LastTradeQuantity float64
	AvgTradePrice     float64
	Change            float64
	ChangePercent     float64
	Timestamp         int64
}

// depthCacheEntry mirrors the Python SDK's depth_data[...] snapshot. Buy/Sell
// hold the raw per-level maps (each with "price"/"quantity"/"orders") as
// received from the server.
type depthCacheEntry struct {
	LTP       float64
	Timestamp int64
	Buy       []map[string]interface{}
	Sell      []map[string]interface{}
}

// toFloat64 best-effort converts a decoded JSON value to float64.
func toFloat64(v interface{}) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case int:
		return float64(t)
	case int64:
		return float64(t)
	case string:
		f, _ := strconv.ParseFloat(t, 64)
		return f
	default:
		return 0
	}
}

// toInt64OrDefault best-effort converts a decoded JSON value to int64,
// falling back to def when v is nil or of an unexpected type.
func toInt64OrDefault(v interface{}, def int64) int64 {
	if v == nil {
		return def
	}
	switch t := v.(type) {
	case float64:
		return int64(t)
	case int64:
		return t
	case int:
		return int64(t)
	default:
		return def
	}
}

// toDepthLevels extracts a "buy"/"sell" depth level array from the raw
// decoded "depth" object.
func toDepthLevels(v interface{}) []map[string]interface{} {
	arr, ok := v.([]interface{})
	if !ok {
		return nil
	}
	levels := make([]map[string]interface{}, 0, len(arr))
	for _, item := range arr {
		if m, ok := item.(map[string]interface{}); ok {
			levels = append(levels, m)
		}
	}
	return levels
}

// splitSymbolKey splits an "EXCHANGE:SYMBOL" cache key back into its parts.
func splitSymbolKey(key string) (exchange, symbol string, ok bool) {
	idx := strings.Index(key, ":")
	if idx < 0 {
		return "", "", false
	}
	return key[:idx], key[idx+1:], true
}

// storeMarketData caches an incoming market_data message in ltpData/
// quotesData/depthData depending on mode, mirroring the Python SDK's
// _process_message bookkeeping. Safe for concurrent use.
func (c *Client) storeMarketData(exchange, symbol string, mode int, marketData map[string]interface{}) {
	if exchange == "" || symbol == "" || marketData == nil {
		return
	}
	key := exchange + ":" + symbol
	now := time.Now().UnixMilli()

	switch mode {
	case 1: // LTP
		if _, ok := marketData["ltp"]; !ok {
			return
		}
		entry := ltpCacheEntry{
			Price:     toFloat64(marketData["ltp"]),
			Timestamp: toInt64OrDefault(marketData["timestamp"], now),
		}
		c.dataMu.Lock()
		c.ltpData[key] = entry
		c.dataMu.Unlock()

	case 2: // Quote
		entry := quoteCacheEntry{
			Open:              toFloat64(marketData["open"]),
			High:              toFloat64(marketData["high"]),
			Low:               toFloat64(marketData["low"]),
			Close:             toFloat64(marketData["close"]),
			LTP:               toFloat64(marketData["ltp"]),
			Volume:            toFloat64(marketData["volume"]),
			LastTradeQuantity: toFloat64(marketData["last_trade_quantity"]),
			AvgTradePrice:     toFloat64(marketData["avg_trade_price"]),
			Change:            toFloat64(marketData["change"]),
			ChangePercent:     toFloat64(marketData["change_percent"]),
			Timestamp:         toInt64OrDefault(marketData["timestamp"], now),
		}
		c.dataMu.Lock()
		c.quotesData[key] = entry
		c.dataMu.Unlock()

	case 3: // Depth
		depthRaw, ok := marketData["depth"].(map[string]interface{})
		if !ok {
			return
		}
		entry := depthCacheEntry{
			LTP:       toFloat64(marketData["ltp"]),
			Timestamp: toInt64OrDefault(marketData["timestamp"], now),
			Buy:       toDepthLevels(depthRaw["buy"]),
			Sell:      toDepthLevels(depthRaw["sell"]),
		}
		c.dataMu.Lock()
		c.depthData[key] = entry
		c.dataMu.Unlock()
	}
}

// GetLTP returns the latest cached LTP snapshots in nested format:
//
//	{"ltp": {"EXCHANGE": {"SYMBOL": {"timestamp": ..., "ltp": ...}}}}
//
// Pass "" for exchange/symbol to skip that filter (both empty returns
// everything currently cached).
func (c *Client) GetLTP(exchange, symbol string) map[string]interface{} {
	c.dataMu.Lock()
	defer c.dataMu.Unlock()

	ltp := map[string]interface{}{}
	for key, entry := range c.ltpData {
		ex, sym, ok := splitSymbolKey(key)
		if !ok {
			continue
		}
		if exchange != "" && ex != exchange {
			continue
		}
		if symbol != "" && sym != symbol {
			continue
		}

		exMap, ok := ltp[ex].(map[string]interface{})
		if !ok {
			exMap = map[string]interface{}{}
			ltp[ex] = exMap
		}
		exMap[sym] = map[string]interface{}{
			"timestamp": entry.Timestamp,
			"ltp":       entry.Price,
		}
	}

	return map[string]interface{}{"ltp": ltp}
}

// GetQuotes returns the latest cached Quote snapshots in nested format:
//
//	{"quote": {"EXCHANGE": {"SYMBOL": {...ohlc + ltp + volume fields...}}}}
//
// Pass "" for exchange/symbol to skip that filter.
func (c *Client) GetQuotes(exchange, symbol string) map[string]interface{} {
	c.dataMu.Lock()
	defer c.dataMu.Unlock()

	quote := map[string]interface{}{}
	for key, entry := range c.quotesData {
		ex, sym, ok := splitSymbolKey(key)
		if !ok {
			continue
		}
		if exchange != "" && ex != exchange {
			continue
		}
		if symbol != "" && sym != symbol {
			continue
		}

		exMap, ok := quote[ex].(map[string]interface{})
		if !ok {
			exMap = map[string]interface{}{}
			quote[ex] = exMap
		}
		exMap[sym] = map[string]interface{}{
			"timestamp":           entry.Timestamp,
			"open":                entry.Open,
			"high":                entry.High,
			"low":                 entry.Low,
			"close":               entry.Close,
			"ltp":                 entry.LTP,
			"volume":              entry.Volume,
			"last_trade_quantity": entry.LastTradeQuantity,
			"avg_trade_price":     entry.AvgTradePrice,
			"change":              entry.Change,
			"change_percent":      entry.ChangePercent,
		}
	}

	return map[string]interface{}{"quote": quote}
}

// buildDepthBook converts raw depth levels into the "1".."5" level map shape
// used by GetDepth, padding missing levels with zeroes (matching the Python SDK).
func buildDepthBook(levels []map[string]interface{}) map[string]interface{} {
	book := map[string]interface{}{}
	for i := 0; i < 5; i++ {
		level := map[string]interface{}{"price": 0.0, "qty": 0, "orders": 0}
		if i < len(levels) {
			l := levels[i]
			level = map[string]interface{}{
				"price":  toFloat64(l["price"]),
				"qty":    int64(toFloat64(l["quantity"])),
				"orders": int64(toFloat64(l["orders"])),
			}
		}
		book[strconv.Itoa(i+1)] = level
	}
	return book
}

// GetDepth returns the latest cached Market Depth snapshots in nested format:
//
//	{"depth": {"EXCHANGE": {"SYMBOL": {"timestamp":..., "ltp":..., "buyBook": {...}, "sellBook": {...}}}}}
//
// buyBook/sellBook always contain levels "1".."5", zero-padded when fewer
// levels were received. Pass "" for exchange/symbol to skip that filter.
func (c *Client) GetDepth(exchange, symbol string) map[string]interface{} {
	c.dataMu.Lock()
	defer c.dataMu.Unlock()

	depth := map[string]interface{}{}
	for key, entry := range c.depthData {
		ex, sym, ok := splitSymbolKey(key)
		if !ok {
			continue
		}
		if exchange != "" && ex != exchange {
			continue
		}
		if symbol != "" && sym != symbol {
			continue
		}

		exMap, ok := depth[ex].(map[string]interface{})
		if !ok {
			exMap = map[string]interface{}{}
			depth[ex] = exMap
		}
		exMap[sym] = map[string]interface{}{
			"timestamp": entry.Timestamp,
			"ltp":       entry.LTP,
			"buyBook":   buildDepthBook(entry.Buy),
			"sellBook":  buildDepthBook(entry.Sell),
		}
	}

	return map[string]interface{}{"depth": depth}
}

// Connect establishes a WebSocket connection and authenticates
func (c *Client) Connect() error {
	if c.wsURL == "" {
		return fmt.Errorf("WebSocket URL not provided")
	}

	conn, _, err := websocket.DefaultDialer.Dial(c.wsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}
	c.wsConn = conn

	// Authenticate using the same format as Python SDK
	authMsg := AuthMessage{
		Action: "authenticate",
		APIKey: c.apiKey,
	}

	if err := c.wsConn.WriteJSON(authMsg); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// Start message reader
	go c.readMessages()

	log.Printf("Connected to %s", c.wsURL)
	return nil
}

// Disconnect closes the WebSocket connection
func (c *Client) Disconnect() error {
	if c.wsConn != nil {
		log.Printf("Disconnected from %s", c.wsURL)
		return c.wsConn.Close()
	}
	return nil
}

// readMessages reads and processes incoming WebSocket messages
func (c *Client) readMessages() {
	for {
		_, message, err := c.wsConn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		// Check for message type
		if msgType, ok := data["type"].(string); ok && msgType == "market_data" {
			// Extract mode to determine data type
			mode := 0
			if m, ok := data["mode"].(float64); ok {
				mode = int(m)
			}

			// Cache the snapshot so GetLTP/GetQuotes/GetDepth can serve it
			// without requiring the caller to wire up a callback.
			exchange, _ := data["exchange"].(string)
			symbol, _ := data["symbol"].(string)
			if marketData, ok := data["data"].(map[string]interface{}); ok {
				c.storeMarketData(exchange, symbol, mode, marketData)
			}

			// Route to appropriate callback based on mode
			switch mode {
			case 1: // LTP
				if callback, exists := c.callbacks["ltp"]; exists {
					callback(data)
				}
			case 2: // Quote
				if callback, exists := c.callbacks["quote"]; exists {
					callback(data)
				}
			case 3: // Depth
				if callback, exists := c.callbacks["depth"]; exists {
					callback(data)
				}
			}
		} else if status, ok := data["status"].(string); ok {
			// Handle status messages
			if message, ok := data["message"].(string); ok {
				log.Printf("WebSocket status: %s - %s", status, message)
			}
		}
	}
}

// SubscribeLTP subscribes to Last Traded Price updates
func (c *Client) SubscribeLTP(instruments []Instrument, onDataReceived func(interface{})) error {
	if c.wsConn == nil {
		return fmt.Errorf("not connected to WebSocket server")
	}

	// Set callback
	if onDataReceived != nil {
		c.callbacks["ltp"] = onDataReceived
	}

	// Subscribe to each instrument individually (matching Python SDK)
	for _, instrument := range instruments {
		symbol := instrument.Symbol
		exchange := instrument.Exchange

		// Use exchange_token as symbol if symbol is not provided
		if symbol == "" && instrument.ExchangeToken != "" {
			symbol = instrument.ExchangeToken
		}

		if exchange == "" || symbol == "" {
			log.Printf("Invalid instrument: %+v", instrument)
			continue
		}

		msg := SubscriptionMessage{
			Action:   "subscribe",
			Symbol:   symbol,
			Exchange: exchange,
			Mode:     1, // 1 for LTP
			Depth:    5, // Default depth level
		}

		log.Printf("Subscribing to %s:%s LTP", exchange, symbol)
		if err := c.wsConn.WriteJSON(msg); err != nil {
			return fmt.Errorf("error subscribing to %s:%s: %w", exchange, symbol, err)
		}

		// Small delay to ensure message is processed separately
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

// UnsubscribeLTP unsubscribes from LTP updates
func (c *Client) UnsubscribeLTP(instruments []Instrument) error {
	if c.wsConn == nil {
		return fmt.Errorf("not connected to WebSocket server")
	}

	// Unsubscribe from each instrument individually
	for _, instrument := range instruments {
		symbol := instrument.Symbol
		exchange := instrument.Exchange

		if symbol == "" && instrument.ExchangeToken != "" {
			symbol = instrument.ExchangeToken
		}

		if exchange == "" || symbol == "" {
			log.Printf("Invalid instrument: %+v", instrument)
			continue
		}

		msg := SubscriptionMessage{
			Action:   "unsubscribe",
			Symbol:   symbol,
			Exchange: exchange,
			Mode:     1, // 1 for LTP
		}

		log.Printf("Unsubscribing from %s:%s LTP", exchange, symbol)
		if err := c.wsConn.WriteJSON(msg); err != nil {
			return fmt.Errorf("error unsubscribing from %s:%s: %w", exchange, symbol, err)
		}

		// Clean up cached LTP data for this symbol (mirrors Python SDK).
		c.dataMu.Lock()
		delete(c.ltpData, exchange+":"+symbol)
		c.dataMu.Unlock()

		time.Sleep(100 * time.Millisecond)
	}

	// Remove callback
	delete(c.callbacks, "ltp")
	return nil
}

// SubscribeQuote subscribes to Quote updates
func (c *Client) SubscribeQuote(instruments []Instrument, onDataReceived func(interface{})) error {
	if c.wsConn == nil {
		return fmt.Errorf("not connected to WebSocket server")
	}

	// Set callback
	if onDataReceived != nil {
		c.callbacks["quote"] = onDataReceived
	}

	// Subscribe to each instrument individually
	for _, instrument := range instruments {
		symbol := instrument.Symbol
		exchange := instrument.Exchange

		if symbol == "" && instrument.ExchangeToken != "" {
			symbol = instrument.ExchangeToken
		}

		if exchange == "" || symbol == "" {
			log.Printf("Invalid instrument: %+v", instrument)
			continue
		}

		msg := SubscriptionMessage{
			Action:   "subscribe",
			Symbol:   symbol,
			Exchange: exchange,
			Mode:     2, // 2 for Quote
			Depth:    5,
		}

		log.Printf("Subscribing to %s:%s Quote", exchange, symbol)
		if err := c.wsConn.WriteJSON(msg); err != nil {
			return fmt.Errorf("error subscribing to %s:%s: %w", exchange, symbol, err)
		}

		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

// UnsubscribeQuote unsubscribes from Quote updates
func (c *Client) UnsubscribeQuote(instruments []Instrument) error {
	if c.wsConn == nil {
		return fmt.Errorf("not connected to WebSocket server")
	}

	// Unsubscribe from each instrument individually
	for _, instrument := range instruments {
		symbol := instrument.Symbol
		exchange := instrument.Exchange

		if symbol == "" && instrument.ExchangeToken != "" {
			symbol = instrument.ExchangeToken
		}

		if exchange == "" || symbol == "" {
			log.Printf("Invalid instrument: %+v", instrument)
			continue
		}

		msg := SubscriptionMessage{
			Action:   "unsubscribe",
			Symbol:   symbol,
			Exchange: exchange,
			Mode:     2, // 2 for Quote
		}

		log.Printf("Unsubscribing from %s:%s Quote", exchange, symbol)
		if err := c.wsConn.WriteJSON(msg); err != nil {
			return fmt.Errorf("error unsubscribing from %s:%s: %w", exchange, symbol, err)
		}

		// Clean up cached quote data for this symbol (mirrors Python SDK).
		c.dataMu.Lock()
		delete(c.quotesData, exchange+":"+symbol)
		c.dataMu.Unlock()

		time.Sleep(100 * time.Millisecond)
	}

	// Remove callback
	delete(c.callbacks, "quote")
	return nil
}

// SubscribeDepth subscribes to Market Depth updates
func (c *Client) SubscribeDepth(instruments []Instrument, onDataReceived func(interface{})) error {
	if c.wsConn == nil {
		return fmt.Errorf("not connected to WebSocket server")
	}

	// Set callback
	if onDataReceived != nil {
		c.callbacks["depth"] = onDataReceived
	}

	// Subscribe to each instrument individually
	for _, instrument := range instruments {
		symbol := instrument.Symbol
		exchange := instrument.Exchange

		if symbol == "" && instrument.ExchangeToken != "" {
			symbol = instrument.ExchangeToken
		}

		if exchange == "" || symbol == "" {
			log.Printf("Invalid instrument: %+v", instrument)
			continue
		}

		msg := SubscriptionMessage{
			Action:   "subscribe",
			Symbol:   symbol,
			Exchange: exchange,
			Mode:     3, // 3 for Depth
			Depth:    5,
		}

		log.Printf("Subscribing to %s:%s Depth", exchange, symbol)
		if err := c.wsConn.WriteJSON(msg); err != nil {
			return fmt.Errorf("error subscribing to %s:%s: %w", exchange, symbol, err)
		}

		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

// UnsubscribeDepth unsubscribes from Market Depth updates
func (c *Client) UnsubscribeDepth(instruments []Instrument) error {
	if c.wsConn == nil {
		return fmt.Errorf("not connected to WebSocket server")
	}

	// Unsubscribe from each instrument individually
	for _, instrument := range instruments {
		symbol := instrument.Symbol
		exchange := instrument.Exchange

		if symbol == "" && instrument.ExchangeToken != "" {
			symbol = instrument.ExchangeToken
		}

		if exchange == "" || symbol == "" {
			log.Printf("Invalid instrument: %+v", instrument)
			continue
		}

		msg := SubscriptionMessage{
			Action:   "unsubscribe",
			Symbol:   symbol,
			Exchange: exchange,
			Mode:     3, // 3 for Depth
		}

		log.Printf("Unsubscribing from %s:%s Depth", exchange, symbol)
		if err := c.wsConn.WriteJSON(msg); err != nil {
			return fmt.Errorf("error unsubscribing from %s:%s: %w", exchange, symbol, err)
		}

		// Clean up cached depth data for this symbol (mirrors Python SDK).
		c.dataMu.Lock()
		delete(c.depthData, exchange+":"+symbol)
		c.dataMu.Unlock()

		time.Sleep(100 * time.Millisecond)
	}

	// Remove callback
	delete(c.callbacks, "depth")
	return nil
}

// SafeWSClient provides thread-safe WebSocket operations
type SafeWSClient struct {
	*Client
	mu sync.Mutex
}

// NewSafeWSClient creates a new thread-safe WebSocket client
func NewSafeWSClient(apiKey string, host string, optionalArgs ...interface{}) *SafeWSClient {
	return &SafeWSClient{
		Client: NewClient(apiKey, host, optionalArgs...),
	}
}

// SubscribeLTP thread-safe LTP subscription
func (s *SafeWSClient) SubscribeLTP(instruments []Instrument, onDataReceived func(interface{})) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Client.SubscribeLTP(instruments, onDataReceived)
}

// UnsubscribeLTP thread-safe LTP unsubscription
func (s *SafeWSClient) UnsubscribeLTP(instruments []Instrument) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Client.UnsubscribeLTP(instruments)
}

// SubscribeQuote thread-safe Quote subscription
func (s *SafeWSClient) SubscribeQuote(instruments []Instrument, onDataReceived func(interface{})) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Client.SubscribeQuote(instruments, onDataReceived)
}

// UnsubscribeQuote thread-safe Quote unsubscription
func (s *SafeWSClient) UnsubscribeQuote(instruments []Instrument) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Client.UnsubscribeQuote(instruments)
}

// SubscribeDepth thread-safe Depth subscription
func (s *SafeWSClient) SubscribeDepth(instruments []Instrument, onDataReceived func(interface{})) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Client.SubscribeDepth(instruments, onDataReceived)
}

// UnsubscribeDepth thread-safe Depth unsubscription
func (s *SafeWSClient) UnsubscribeDepth(instruments []Instrument) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Client.UnsubscribeDepth(instruments)
}