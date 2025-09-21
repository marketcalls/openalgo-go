package openalgo

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Instrument struct {
	Exchange string `json:"exchange"`
	Symbol   string `json:"symbol"`
}

type WSMessage struct {
	Action      string       `json:"action"`
	Instruments []Instrument `json:"instruments,omitempty"`
	APIKey      string       `json:"apikey,omitempty"`
}

type StreamData struct {
	Exchange string      `json:"exchange"`
	Symbol   string      `json:"symbol"`
	LTP      float64     `json:"ltp,omitempty"`
	Quote    interface{} `json:"quote,omitempty"`
	Depth    interface{} `json:"depth,omitempty"`
}

func (c *Client) Connect() error {
	if c.wsURL == "" {
		return fmt.Errorf("WebSocket URL not provided")
	}

	conn, _, err := websocket.DefaultDialer.Dial(c.wsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}
	c.wsConn = conn

	authMsg := WSMessage{
		Action: "auth",
		APIKey: c.apiKey,
	}

	if err := c.wsConn.WriteJSON(authMsg); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	go c.readMessages()

	return nil
}

func (c *Client) Disconnect() error {
	if c.wsConn != nil {
		return c.wsConn.Close()
	}
	return nil
}

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

		if msgType, ok := data["type"].(string); ok {
			switch msgType {
			case "ltp":
				if callback, exists := c.callbacks["ltp"]; exists {
					callback(data)
				}
			case "quote":
				if callback, exists := c.callbacks["quote"]; exists {
					callback(data)
				}
			case "depth":
				if callback, exists := c.callbacks["depth"]; exists {
					callback(data)
				}
			}
		}
	}
}

func (c *Client) SubscribeLTP(instruments []Instrument, onDataReceived func(interface{})) error {
	c.callbacks["ltp"] = onDataReceived

	msg := WSMessage{
		Action:      "subscribe_ltp",
		Instruments: instruments,
	}

	return c.wsConn.WriteJSON(msg)
}

func (c *Client) UnsubscribeLTP(instruments []Instrument) error {
	delete(c.callbacks, "ltp")

	msg := WSMessage{
		Action:      "unsubscribe_ltp",
		Instruments: instruments,
	}

	return c.wsConn.WriteJSON(msg)
}

func (c *Client) SubscribeQuote(instruments []Instrument, onDataReceived func(interface{})) error {
	c.callbacks["quote"] = onDataReceived

	msg := WSMessage{
		Action:      "subscribe_quote",
		Instruments: instruments,
	}

	return c.wsConn.WriteJSON(msg)
}

func (c *Client) UnsubscribeQuote(instruments []Instrument) error {
	delete(c.callbacks, "quote")

	msg := WSMessage{
		Action:      "unsubscribe_quote",
		Instruments: instruments,
	}

	return c.wsConn.WriteJSON(msg)
}

func (c *Client) SubscribeDepth(instruments []Instrument, onDataReceived func(interface{})) error {
	c.callbacks["depth"] = onDataReceived

	msg := WSMessage{
		Action:      "subscribe_depth",
		Instruments: instruments,
	}

	return c.wsConn.WriteJSON(msg)
}

func (c *Client) UnsubscribeDepth(instruments []Instrument) error {
	delete(c.callbacks, "depth")

	msg := WSMessage{
		Action:      "unsubscribe_depth",
		Instruments: instruments,
	}

	return c.wsConn.WriteJSON(msg)
}

type SafeWSClient struct {
	*Client
	mu sync.Mutex
}

func NewSafeWSClient(apiKey string, host string, wsURL string) *SafeWSClient {
	return &SafeWSClient{
		Client: NewClient(apiKey, host, wsURL),
	}
}

func (s *SafeWSClient) SubscribeLTP(instruments []Instrument, onDataReceived func(interface{})) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Client.SubscribeLTP(instruments, onDataReceived)
}

func (s *SafeWSClient) UnsubscribeLTP(instruments []Instrument) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Client.UnsubscribeLTP(instruments)
}

func (s *SafeWSClient) SubscribeQuote(instruments []Instrument, onDataReceived func(interface{})) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Client.SubscribeQuote(instruments, onDataReceived)
}

func (s *SafeWSClient) UnsubscribeQuote(instruments []Instrument) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Client.UnsubscribeQuote(instruments)
}

func (s *SafeWSClient) SubscribeDepth(instruments []Instrument, onDataReceived func(interface{})) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Client.SubscribeDepth(instruments, onDataReceived)
}

func (s *SafeWSClient) UnsubscribeDepth(instruments []Instrument) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Client.UnsubscribeDepth(instruments)
}