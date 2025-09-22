package openalgo

// Ping checks API connectivity
func (c *Client) Ping() (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey": c.apiKey,
	}
	return c.makeRequest("POST", "ping", payload)
}