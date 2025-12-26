package openalgo

import (
	"net/http"
	"time"
)

type QuotesRequest struct {
	Symbol   string `json:"symbol"`
	Exchange string `json:"exchange"`
}

type DepthRequest struct {
	Symbol   string `json:"symbol"`
	Exchange string `json:"exchange"`
}

type HistoryRequest struct {
	Symbol    string `json:"symbol"`
	Exchange  string `json:"exchange"`
	Interval  string `json:"interval"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type SymbolRequest struct {
	Symbol   string `json:"symbol"`
	Exchange string `json:"exchange"`
}

type SearchRequest struct {
	Query    string `json:"query"`
	Exchange string `json:"exchange"`
}

type ExpiryRequest struct {
	Symbol         string `json:"symbol"`
	Exchange       string `json:"exchange"`
	InstrumentType string `json:"instrumenttype"`
}

type QuotesResponse struct {
	Status string `json:"status"`
	Data   struct {
		Open      float64 `json:"open"`
		High      float64 `json:"high"`
		Low       float64 `json:"low"`
		LTP       float64 `json:"ltp"`
		Ask       float64 `json:"ask"`
		Bid       float64 `json:"bid"`
		PrevClose float64 `json:"prev_close"`
		Volume    int64   `json:"volume"`
	} `json:"data"`
}

type DepthResponse struct {
	Status string `json:"status"`
	Data   struct {
		Open         float64 `json:"open"`
		High         float64 `json:"high"`
		Low          float64 `json:"low"`
		LTP          float64 `json:"ltp"`
		LTQ          int     `json:"ltq"`
		PrevClose    float64 `json:"prev_close"`
		Volume       int64   `json:"volume"`
		OI           int64   `json:"oi"`
		TotalBuyQty  int64   `json:"totalbuyqty"`
		TotalSellQty int64   `json:"totalsellqty"`
		Asks         []struct {
			Price    float64 `json:"price"`
			Quantity int     `json:"quantity"`
		} `json:"asks"`
		Bids []struct {
			Price    float64 `json:"price"`
			Quantity int     `json:"quantity"`
		} `json:"bids"`
	} `json:"data"`
}

type HistoryBar struct {
	Timestamp time.Time `json:"timestamp"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    int64     `json:"volume"`
}

func (c *Client) Quotes(symbol, exchange string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey":   c.apiKey,
		"symbol":   symbol,
		"exchange": exchange,
	}
	return c.makeRequest(http.MethodPost, "quotes", payload)
}

func (c *Client) Depth(symbol, exchange string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey":   c.apiKey,
		"symbol":   symbol,
		"exchange": exchange,
	}
	return c.makeRequest(http.MethodPost, "depth", payload)
}

func (c *Client) History(symbol, exchange, interval, startDate, endDate string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey":     c.apiKey,
		"symbol":     symbol,
		"exchange":   exchange,
		"interval":   interval,
		"start_date": startDate,
		"end_date":   endDate,
	}
	return c.makeRequest(http.MethodPost, "history", payload)
}

func (c *Client) Intervals() (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey": c.apiKey,
	}
	return c.makeRequest(http.MethodPost, "intervals", payload)
}

func (c *Client) Symbol(symbol, exchange string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey":   c.apiKey,
		"symbol":   symbol,
		"exchange": exchange,
	}
	return c.makeRequest(http.MethodPost, "symbol", payload)
}

func (c *Client) Search(query, exchange string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey": c.apiKey,
		"query":  query,
	}
	if exchange != "" {
		payload["exchange"] = exchange
	}
	return c.makeRequest(http.MethodPost, "search", payload)
}

func (c *Client) Expiry(symbol, exchange, instrumentType string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey":         c.apiKey,
		"symbol":         symbol,
		"exchange":       exchange,
		"instrumenttype": instrumentType,
	}
	return c.makeRequest(http.MethodPost, "expiry", payload)
}

// MultiQuotes gets quotes for multiple symbols
func (c *Client) MultiQuotes(symbols []map[string]string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey":  c.apiKey,
		"symbols": symbols,
	}
	return c.makeRequest(http.MethodPost, "multiquotes", payload)
}

// OptionChain gets option chain data for an underlying
func (c *Client) OptionChain(underlying, exchange, expiryDate string, strikeCount ...int) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey":      c.apiKey,
		"underlying":  underlying,
		"exchange":    exchange,
		"expiry_date": expiryDate,
	}
	if len(strikeCount) > 0 {
		payload["strike_count"] = strikeCount[0]
	}
	return c.makeRequest(http.MethodPost, "optionchain", payload)
}

// OptionSymbol gets option symbol based on offset from ATM
func (c *Client) OptionSymbol(underlying, exchange, expiryDate, offset, optionType string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey":      c.apiKey,
		"underlying":  underlying,
		"exchange":    exchange,
		"expiry_date": expiryDate,
		"offset":      offset,
		"option_type": optionType,
	}
	return c.makeRequest(http.MethodPost, "optionsymbol", payload)
}

// SyntheticFuture calculates synthetic future price from ATM options
func (c *Client) SyntheticFuture(underlying, exchange, expiryDate string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey":      c.apiKey,
		"underlying":  underlying,
		"exchange":    exchange,
		"expiry_date": expiryDate,
	}
	return c.makeRequest(http.MethodPost, "syntheticfuture", payload)
}

// OptionGreeks calculates option Greeks for a symbol
func (c *Client) OptionGreeks(symbol, exchange string, interestRate float64, underlyingSymbol, underlyingExchange string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey":              c.apiKey,
		"symbol":              symbol,
		"exchange":            exchange,
		"interest_rate":       interestRate,
		"underlying_symbol":   underlyingSymbol,
		"underlying_exchange": underlyingExchange,
	}
	return c.makeRequest(http.MethodPost, "optiongreeks", payload)
}

// Instruments gets all instruments for an exchange
func (c *Client) Instruments(exchange string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey":   c.apiKey,
		"exchange": exchange,
	}
	return c.makeRequest(http.MethodPost, "instruments", payload)
}
