package openalgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Strategy is a standalone TradingView-style webhook client for sending
// strategy signals to an OpenAlgo server. Unlike Client, it does not need an
// API key - the per-strategy webhook ID (shown on the OpenAlgo "Strategies"
// page) authorizes the request instead.
type Strategy struct {
	hostURL    string
	webhookID  string
	webhookURL string
	httpClient *http.Client
}

// NewStrategy creates a new Strategy webhook client.
//
// hostURL is the OpenAlgo server URL (e.g. "http://127.0.0.1:5000").
// webhookID is the strategy's webhook ID as shown in the OpenAlgo UI.
func NewStrategy(hostURL, webhookID string) *Strategy {
	return &Strategy{
		hostURL:   strings.TrimRight(hostURL, "/"),
		webhookID: webhookID,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// WebhookURL returns the fully-qualified webhook URL, building and caching
// it the first time it is requested.
func (s *Strategy) WebhookURL() string {
	if s.webhookURL == "" {
		s.webhookURL = fmt.Sprintf("%s/strategy/webhook/%s", s.hostURL, s.webhookID)
	}
	return s.webhookURL
}

// StrategyOrder sends a strategy order signal via webhook to OpenAlgo. The
// strategy mode (LONG_ONLY, SHORT_ONLY, BOTH) is configured on the OpenAlgo
// server for this webhook/strategy.
//
// symbol is the trading symbol (e.g. "RELIANCE", "NIFTY").
// action is "BUY" or "SELL".
// positionSize is optional - pass nil to omit it. Required for BOTH mode.
// Accepts a string, int, or float64.
func (s *Strategy) StrategyOrder(symbol, action string, positionSize interface{}) (map[string]interface{}, error) {
	postMessage := map[string]interface{}{
		"symbol": symbol,
		"action": strings.ToUpper(action),
	}

	if positionSize != nil {
		switch v := positionSize.(type) {
		case string:
			postMessage["position_size"] = v
		case int:
			postMessage["position_size"] = fmt.Sprintf("%d", v)
		case float64:
			postMessage["position_size"] = fmt.Sprintf("%.0f", v)
		default:
			postMessage["position_size"] = fmt.Sprintf("%v", v)
		}
	}

	jsonData, err := json.Marshal(postMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	req, err := http.NewRequest("POST", s.WebhookURL(), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create webhook request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("strategy order failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read webhook response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("strategy order failed: HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal webhook response: %w", err)
	}

	return result, nil
}
