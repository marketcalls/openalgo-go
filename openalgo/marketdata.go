package openalgo

import (
	"fmt"
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
	return c.makeRequest("POST", "quotes", payload)
}

func (c *Client) Depth(symbol, exchange string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey":   c.apiKey,
		"symbol":   symbol,
		"exchange": exchange,
	}
	return c.makeRequest("POST", "depth", payload)
}

// History fetches historical OHLCV data for a symbol.
//
// source is optional: "api" (default) fetches from the broker API, "db"
// fetches from the OpenAlgo DuckDB/Historify database (required for custom
// intraday intervals and multiplier-based daily intervals like "2W"/"3M").
func (c *Client) History(symbol, exchange, interval, startDate, endDate string, source ...string) (map[string]interface{}, error) {
	src := "api"
	if len(source) > 0 && source[0] != "" {
		src = source[0]
	}

	payload := map[string]interface{}{
		"apikey":     c.apiKey,
		"symbol":     symbol,
		"exchange":   exchange,
		"interval":   interval,
		"start_date": startDate,
		"end_date":   endDate,
		"source":     src,
	}
	return c.makeRequest("POST", "history", payload)
}

func (c *Client) Intervals() (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey": c.apiKey,
	}
	return c.makeRequest("POST", "intervals", payload)
}

func (c *Client) Symbol(symbol, exchange string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey":   c.apiKey,
		"symbol":   symbol,
		"exchange": exchange,
	}
	return c.makeRequest("POST", "symbol", payload)
}

func (c *Client) Search(query, exchange string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey": c.apiKey,
		"query":  query,
	}
	if exchange != "" {
		payload["exchange"] = exchange
	}
	return c.makeRequest("POST", "search", payload)
}

func (c *Client) Expiry(symbol, exchange, instrumentType string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey":         c.apiKey,
		"symbol":         symbol,
		"exchange":       exchange,
		"instrumenttype": instrumentType,
	}
	return c.makeRequest("POST", "expiry", payload)
}

// MultiQuotes gets quotes for multiple symbols
func (c *Client) MultiQuotes(symbols []map[string]string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey":  c.apiKey,
		"symbols": symbols,
	}
	return c.makeRequest("POST", "multiquotes", payload)
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
	return c.makeRequest("POST", "optionchain", payload)
}

// OptionSymbol resolves the option symbol, lot size, and tick size for an
// underlying and strike offset (ATM, ITM1-ITM50, OTM1-OTM50) without
// placing an order.
//
// optionalParams may include:
//   - "strike_int" (int): DEPRECATED - strike interval override (e.g. 50 for NIFTY).
func (c *Client) OptionSymbol(underlying, exchange, expiryDate, offset, optionType string, optionalParams ...map[string]interface{}) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey":      c.apiKey,
		"underlying":  underlying,
		"exchange":    exchange,
		"expiry_date": expiryDate,
		"offset":      offset,
		"option_type": optionType,
	}

	if len(optionalParams) > 0 {
		for key, value := range optionalParams[0] {
			if value == nil {
				continue
			}
			switch v := value.(type) {
			case string:
				payload[key] = v
			case int:
				payload[key] = fmt.Sprintf("%d", v)
			case float64:
				payload[key] = fmt.Sprintf("%g", v)
			default:
				payload[key] = fmt.Sprintf("%v", v)
			}
		}
	}

	return c.makeRequest("POST", "optionsymbol", payload)
}

// SyntheticFuture calculates synthetic future price from ATM options
func (c *Client) SyntheticFuture(underlying, exchange, expiryDate string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey":      c.apiKey,
		"underlying":  underlying,
		"exchange":    exchange,
		"expiry_date": expiryDate,
	}
	return c.makeRequest("POST", "syntheticfuture", payload)
}

// OptionGreeks calculates Option Greeks (Delta, Gamma, Theta, Vega, Rho) and
// Implied Volatility for an option using the Black-76 model. Requires
// real-time LTP for the underlying and option, unless "forward_price" is
// supplied.
//
// optionalParams may include:
//   - "interest_rate" (float64): risk-free annualized rate. Defaults to 0 server-side.
//   - "forward_price" (float64): custom forward/synthetic futures price; skips the
//     underlying price fetch (useful for FINNIFTY/MIDCPNIFTY or scenario analysis).
//   - "underlying_symbol" (string): override the auto-detected underlying symbol.
//   - "underlying_exchange" (string): override the auto-detected underlying exchange.
//   - "expiry_time" (string): custom expiry time "HH:MM", required for MCX contracts
//     with non-standard expiry times.
func (c *Client) OptionGreeks(symbol, exchange string, optionalParams ...map[string]interface{}) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey":   c.apiKey,
		"symbol":   symbol,
		"exchange": exchange,
	}

	if len(optionalParams) > 0 {
		for key, value := range optionalParams[0] {
			if value != nil {
				payload[key] = value
			}
		}
	}

	return c.makeRequest("POST", "optiongreeks", payload)
}

// allInstrumentExchanges is queried and merged when Instruments is called
// without an exchange filter, matching the Python SDK's "download all
// exchanges" behaviour.
var allInstrumentExchanges = []string{"NSE", "BSE", "NFO", "BFO", "MCX", "CDS", "BCD", "NSE_INDEX", "BSE_INDEX"}

// Instruments downloads trading symbols and instruments for an exchange.
//
// Pass an empty string to download and merge instruments across ALL
// supported exchanges (NSE, BSE, NFO, BFO, MCX, CDS, BCD, NSE_INDEX, BSE_INDEX).
func (c *Client) Instruments(exchange string) (map[string]interface{}, error) {
	if exchange != "" {
		payload := map[string]interface{}{
			"apikey":   c.apiKey,
			"exchange": exchange,
		}
		return c.makeRequest("POST", "instruments", payload)
	}

	// No exchange specified - fetch every supported exchange and combine.
	var combined []interface{}
	for _, ex := range allInstrumentExchanges {
		payload := map[string]interface{}{
			"apikey":   c.apiKey,
			"exchange": ex,
		}
		result, err := c.makeRequest("POST", "instruments", payload)
		if err != nil {
			continue // skip exchanges that fail, matching Python SDK behaviour
		}
		if data, ok := result["data"].([]interface{}); ok {
			combined = append(combined, data...)
		}
	}

	if len(combined) == 0 {
		return nil, fmt.Errorf("failed to fetch instruments from any exchange")
	}

	return map[string]interface{}{
		"status": "success",
		"data":   combined,
	}, nil
}
