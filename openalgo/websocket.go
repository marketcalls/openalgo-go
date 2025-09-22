package openalgo

import (
	"encoding/json"
	"fmt"
	"log"
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