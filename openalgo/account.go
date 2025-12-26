package openalgo

import "net/http"

type FundsResponse struct {
	Status string `json:"status"`
	Data   struct {
		AvailableCash  string `json:"availablecash"`
		Collateral     string `json:"collateral"`
		M2MRealized    string `json:"m2mrealized"`
		M2MUnrealized  string `json:"m2munrealized"`
		UtilisedDebits string `json:"utiliseddebits"`
	} `json:"data"`
}

type OrderBookResponse struct {
	Status string `json:"status"`
	Data   struct {
		Orders []struct {
			Action       string  `json:"action"`
			Symbol       string  `json:"symbol"`
			Exchange     string  `json:"exchange"`
			OrderID      string  `json:"orderid"`
			Product      string  `json:"product"`
			Quantity     string  `json:"quantity"`
			Price        float64 `json:"price"`
			PriceType    string  `json:"pricetype"`
			OrderStatus  string  `json:"order_status"`
			TriggerPrice float64 `json:"trigger_price"`
			Timestamp    string  `json:"timestamp"`
		} `json:"orders"`
		Statistics struct {
			TotalBuyOrders       float64 `json:"total_buy_orders"`
			TotalSellOrders      float64 `json:"total_sell_orders"`
			TotalCompletedOrders float64 `json:"total_completed_orders"`
			TotalOpenOrders      float64 `json:"total_open_orders"`
			TotalRejectedOrders  float64 `json:"total_rejected_orders"`
		} `json:"statistics"`
	} `json:"data"`
}

type TradeBookResponse struct {
	Status string `json:"status"`
	Data   []struct {
		Action       string  `json:"action"`
		Symbol       string  `json:"symbol"`
		Exchange     string  `json:"exchange"`
		OrderID      string  `json:"orderid"`
		Product      string  `json:"product"`
		Quantity     float64 `json:"quantity"`
		AveragePrice float64 `json:"average_price"`
		Timestamp    string  `json:"timestamp"`
		TradeValue   float64 `json:"trade_value"`
	} `json:"data"`
}

type PositionBookResponse struct {
	Status string `json:"status"`
	Data   []struct {
		Symbol       string `json:"symbol"`
		Exchange     string `json:"exchange"`
		Product      string `json:"product"`
		Quantity     string `json:"quantity"`
		AveragePrice string `json:"average_price"`
		LTP          string `json:"ltp"`
		PnL          string `json:"pnl"`
	} `json:"data"`
}

type HoldingsResponse struct {
	Status string `json:"status"`
	Data   struct {
		Holdings []struct {
			Symbol     string  `json:"symbol"`
			Exchange   string  `json:"exchange"`
			Product    string  `json:"product"`
			Quantity   int     `json:"quantity"`
			PnL        float64 `json:"pnl"`
			PnLPercent float64 `json:"pnlpercent"`
		} `json:"holdings"`
		Statistics struct {
			TotalHoldingValue float64 `json:"totalholdingvalue"`
			TotalInvValue     float64 `json:"totalinvvalue"`
			TotalPnL          float64 `json:"totalprofitandloss"`
			TotalPnLPercent   float64 `json:"totalpnlpercentage"`
		} `json:"statistics"`
	} `json:"data"`
}

func (c *Client) Funds() (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey": c.apiKey,
	}
	return c.makeRequest(http.MethodPost, "funds", payload)
}

func (c *Client) OrderBook() (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey": c.apiKey,
	}
	return c.makeRequest(http.MethodPost, "orderbook", payload)
}

func (c *Client) TradeBook() (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey": c.apiKey,
	}
	return c.makeRequest(http.MethodPost, "tradebook", payload)
}

func (c *Client) PositionBook() (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey": c.apiKey,
	}
	return c.makeRequest(http.MethodPost, "positionbook", payload)
}

func (c *Client) Holdings() (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"apikey": c.apiKey,
	}
	return c.makeRequest(http.MethodPost, "holdings", payload)
}

// MarginPosition represents a position for margin calculation
type MarginPosition struct {
	Symbol    string `json:"symbol"`
	Exchange  string `json:"exchange"`
	Action    string `json:"action"`
	Product   string `json:"product"`
	PriceType string `json:"pricetype"`
	Quantity  string `json:"quantity"`
}

// Margin calculates margin requirements for positions
func (c *Client) Margin(positions []MarginPosition) (map[string]interface{}, error) {
	// Convert positions to interface slice
	positionsPayload := make([]map[string]interface{}, len(positions))
	for i, pos := range positions {
		positionsPayload[i] = map[string]interface{}{
			"symbol":    pos.Symbol,
			"exchange":  pos.Exchange,
			"action":    pos.Action,
			"product":   pos.Product,
			"pricetype": pos.PriceType,
			"quantity":  pos.Quantity,
		}
	}

	payload := map[string]interface{}{
		"apikey":    c.apiKey,
		"positions": positionsPayload,
	}
	return c.makeRequest(http.MethodPost, "margin", payload)
}
