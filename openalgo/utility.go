package openalgo

import "time"

// Ping checks API connectivity
func (c *Client) Ping() (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey": c.apiKey,
	}
	return c.makeRequest("POST", "ping", payload)
}

// Holidays gets trading holidays for a year (2020-2050).
//
// year is optional - pass no argument (or 0) to default to the current year.
func (c *Client) Holidays(year ...int) (map[string]interface{}, error) {
	y := time.Now().Year()
	if len(year) > 0 && year[0] != 0 {
		y = year[0]
	}

	payload := map[string]interface{}{
		"apikey": c.apiKey,
		"year":   y,
	}
	return c.makeRequest("POST", "market/holidays", payload)
}

// Timings gets exchange trading timings for a specific date.
//
// date is optional (format "YYYY-MM-DD") - pass no argument (or an empty
// string) to default to today's date.
func (c *Client) Timings(date ...string) (map[string]interface{}, error) {
	d := time.Now().Format("2006-01-02")
	if len(date) > 0 && date[0] != "" {
		d = date[0]
	}

	payload := map[string]interface{}{
		"apikey": c.apiKey,
		"date":   d,
	}
	return c.makeRequest("POST", "market/timings", payload)
}

// Telegram sends a Telegram notification with default priority (5)
func (c *Client) Telegram(username, message string) (map[string]interface{}, error) {
	return c.TelegramWithPriority(username, message, 5)
}

// TelegramWithPriority sends a Telegram notification with custom priority
func (c *Client) TelegramWithPriority(username, message string, priority int) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey":   c.apiKey,
		"username": username,
		"message":  message,
		"priority": priority,
	}
	return c.makeRequest("POST", "telegram/notify", payload)
}