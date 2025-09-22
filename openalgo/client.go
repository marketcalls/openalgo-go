package openalgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Client is the main OpenAlgo API client
type Client struct {
	apiKey    string
	host      string
	baseURL   string
	wsURL     string
	wsPort    int
	client    *http.Client
	wsConn    *websocket.Conn
	callbacks map[string]func(interface{})
}

// NewClient creates a new OpenAlgo API client
func NewClient(apiKey string, host string, optionalArgs ...interface{}) *Client {
	version := "v1"
	wsPort := 8765
	var wsURL string

	// Parse optional arguments
	for i, arg := range optionalArgs {
		switch v := arg.(type) {
		case string:
			if i == 0 && v != "" {
				version = v
			} else if i > 0 && v != "" {
				wsURL = v
			}
		case int:
			wsPort = v
		}
	}

	c := &Client{
		apiKey:    apiKey,
		host:      host,
		baseURL:   fmt.Sprintf("%s/api/%s/", host, version),
		wsPort:    wsPort,
		client:    &http.Client{Timeout: 30 * time.Second},
		callbacks: make(map[string]func(interface{})),
	}

	// Set WebSocket URL
	if wsURL != "" {
		c.wsURL = wsURL
	} else {
		// Extract host without protocol for WebSocket
		wsHost := host
		if len(host) > 7 && host[:7] == "http://" {
			wsHost = host[7:]
		} else if len(host) > 8 && host[:8] == "https://" {
			wsHost = host[8:]
		}
		// Remove port if present
		for i, ch := range wsHost {
			if ch == ':' || ch == '/' {
				wsHost = wsHost[:i]
				break
			}
		}
		c.wsURL = fmt.Sprintf("ws://%s:%d", wsHost, wsPort)
	}

	return c
}

// makeRequest performs an HTTP request to the OpenAlgo API
func (c *Client) makeRequest(method, endpoint string, payload interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)

	var req *http.Request
	var err error

	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check if response is JSON
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		// If response is not JSON, include the actual response in error for debugging
		if len(body) > 200 {
			return nil, fmt.Errorf("failed to unmarshal response: %w (response: %s...)", err, string(body[:200]))
		}
		return nil, fmt.Errorf("failed to unmarshal response: %w (response: %s)", err, string(body))
	}

	// Check if API returned an error
	if status, ok := result["status"].(string); ok && status == "error" {
		if msg, ok := result["message"].(string); ok {
			return nil, fmt.Errorf("API error: %s", msg)
		}
		return nil, fmt.Errorf("API error: %v", result)
	}

	return result, nil
}