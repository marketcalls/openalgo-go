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

type Client struct {
	apiKey    string
	host      string
	wsURL     string
	client    *http.Client
	wsConn    *websocket.Conn
	callbacks map[string]func(interface{})
}

func NewClient(apiKey string, host string, wsURL ...string) *Client {
	c := &Client{
		apiKey:    apiKey,
		host:      host,
		client:    &http.Client{Timeout: 30 * time.Second},
		callbacks: make(map[string]func(interface{})),
	}

	if len(wsURL) > 0 {
		c.wsURL = wsURL[0]
	}

	return c
}

func (c *Client) Ping() (map[string]interface{}, error) {
	return c.makeRequest("GET", "/api/v1/ping", nil)
}

func (c *Client) makeRequest(method, endpoint string, payload interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s%s", c.host, endpoint)

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
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
	}

	req.Header.Set("X-API-KEY", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}