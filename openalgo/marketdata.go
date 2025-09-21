package openalgo

import (
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
	req := QuotesRequest{
		Symbol:   symbol,
		Exchange: exchange,
	}
	return c.makeRequest("POST", "/api/v1/quotes", req)
}

func (c *Client) Depth(symbol, exchange string) (map[string]interface{}, error) {
	req := DepthRequest{
		Symbol:   symbol,
		Exchange: exchange,
	}
	return c.makeRequest("POST", "/api/v1/depth", req)
}

func (c *Client) History(symbol, exchange, interval, startDate, endDate string) (map[string]interface{}, error) {
	req := HistoryRequest{
		Symbol:    symbol,
		Exchange:  exchange,
		Interval:  interval,
		StartDate: startDate,
		EndDate:   endDate,
	}
	return c.makeRequest("POST", "/api/v1/history", req)
}

func (c *Client) Intervals() (map[string]interface{}, error) {
	return c.makeRequest("GET", "/api/v1/intervals", nil)
}

func (c *Client) Symbol(symbol, exchange string) (map[string]interface{}, error) {
	req := SymbolRequest{
		Symbol:   symbol,
		Exchange: exchange,
	}
	return c.makeRequest("POST", "/api/v1/symbol", req)
}

func (c *Client) Search(query, exchange string) (map[string]interface{}, error) {
	req := SearchRequest{
		Query:    query,
		Exchange: exchange,
	}
	return c.makeRequest("POST", "/api/v1/search", req)
}

func (c *Client) Expiry(symbol, exchange, instrumentType string) (map[string]interface{}, error) {
	req := ExpiryRequest{
		Symbol:         symbol,
		Exchange:       exchange,
		InstrumentType: instrumentType,
	}
	return c.makeRequest("POST", "/api/v1/expiry", req)
}