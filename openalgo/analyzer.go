package openalgo

type AnalyzerStatusResponse struct {
	Status string `json:"status"`
	Data   struct {
		AnalyzeMode bool   `json:"analyze_mode"`
		Mode        string `json:"mode"`
		TotalLogs   int    `json:"total_logs"`
	} `json:"data"`
}

type AnalyzerToggleRequest struct {
	Mode bool `json:"mode"`
}

type AnalyzerToggleResponse struct {
	Status string `json:"status"`
	Data   struct {
		AnalyzeMode bool   `json:"analyze_mode"`
		Message     string `json:"message"`
		Mode        string `json:"mode"`
		TotalLogs   int    `json:"total_logs"`
	} `json:"data"`
}

func (c *Client) AnalyzerStatus() (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey": c.apiKey,
	}
	return c.makeRequest("POST", "analyzer", payload)
}

func (c *Client) AnalyzerToggle(mode bool) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey": c.apiKey,
		"mode":   mode,
	}
	return c.makeRequest("POST", "analyzer/toggle", payload)
}