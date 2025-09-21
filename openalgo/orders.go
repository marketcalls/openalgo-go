package openalgo

import (
	"fmt"
)

type PlaceOrderRequest struct {
	Strategy           string  `json:"strategy"`
	Symbol             string  `json:"symbol"`
	Action             string  `json:"action"`
	Exchange           string  `json:"exchange"`
	PriceType          string  `json:"pricetype"`
	Product            string  `json:"product"`
	Quantity           string  `json:"quantity"`
	Price              float64 `json:"price,omitempty"`
	TriggerPrice       float64 `json:"trigger_price,omitempty"`
	DisclosedQuantity  string  `json:"disclosed_quantity,omitempty"`
}

type PlaceSmartOrderRequest struct {
	Strategy      string  `json:"strategy"`
	Symbol        string  `json:"symbol"`
	Action        string  `json:"action"`
	Exchange      string  `json:"exchange"`
	PriceType     string  `json:"pricetype"`
	Product       string  `json:"product"`
	Quantity      string  `json:"quantity"`
	PositionSize  int     `json:"position_size"`
	Price         float64 `json:"price,omitempty"`
	TriggerPrice  float64 `json:"trigger_price,omitempty"`
	DisclosedQuantity string `json:"disclosed_quantity,omitempty"`
}

type BasketOrderItem struct {
	Symbol        string  `json:"symbol"`
	Exchange      string  `json:"exchange"`
	Action        string  `json:"action"`
	Quantity      int     `json:"quantity"`
	PriceType     string  `json:"pricetype"`
	Product       string  `json:"product"`
	Price         float64 `json:"price,omitempty"`
	TriggerPrice  float64 `json:"trigger_price,omitempty"`
}

type BasketOrderRequest struct {
	Orders []BasketOrderItem `json:"orders"`
}

type SplitOrderRequest struct {
	Symbol       string  `json:"symbol"`
	Exchange     string  `json:"exchange"`
	Action       string  `json:"action"`
	Quantity     int     `json:"quantity"`
	SplitSize    int     `json:"splitsize"`
	PriceType    string  `json:"pricetype"`
	Product      string  `json:"product"`
	Price        float64 `json:"price,omitempty"`
	TriggerPrice float64 `json:"trigger_price,omitempty"`
}

type ModifyOrderRequest struct {
	OrderID           string  `json:"orderid"`
	Strategy          string  `json:"strategy"`
	Symbol            string  `json:"symbol"`
	Action            string  `json:"action"`
	Exchange          string  `json:"exchange"`
	PriceType         string  `json:"pricetype"`
	Product           string  `json:"product"`
	Quantity          string  `json:"quantity"`
	Price             float64 `json:"price,omitempty"`
	TriggerPrice      float64 `json:"trigger_price,omitempty"`
	DisclosedQuantity string  `json:"disclosed_quantity,omitempty"`
}

type CancelOrderRequest struct {
	OrderID  string `json:"orderid"`
	Strategy string `json:"strategy"`
}

type CancelAllOrderRequest struct {
	Strategy string `json:"strategy"`
}

type ClosePositionRequest struct {
	Strategy string `json:"strategy"`
}

type OrderStatusRequest struct {
	OrderID  string `json:"orderid"`
	Strategy string `json:"strategy"`
}

type OpenPositionRequest struct {
	Strategy string `json:"strategy"`
	Symbol   string `json:"symbol"`
	Exchange string `json:"exchange"`
	Product  string `json:"product"`
}

func (c *Client) PlaceOrder(strategy, symbol, action, exchange, priceType, product string, quantity interface{}, optionalParams ...map[string]interface{}) (map[string]interface{}, error) {
	req := PlaceOrderRequest{
		Strategy:  strategy,
		Symbol:    symbol,
		Action:    action,
		Exchange:  exchange,
		PriceType: priceType,
		Product:   product,
	}

	switch v := quantity.(type) {
	case string:
		req.Quantity = v
	case int:
		req.Quantity = fmt.Sprintf("%d", v)
	case float64:
		req.Quantity = fmt.Sprintf("%.0f", v)
	default:
		return nil, fmt.Errorf("quantity must be string, int, or float64")
	}

	if len(optionalParams) > 0 {
		params := optionalParams[0]
		if price, ok := params["price"].(float64); ok {
			req.Price = price
		} else if price, ok := params["price"].(string); ok {
			var p float64
			fmt.Sscanf(price, "%f", &p)
			req.Price = p
		}

		if triggerPrice, ok := params["trigger_price"].(float64); ok {
			req.TriggerPrice = triggerPrice
		} else if triggerPrice, ok := params["trigger_price"].(string); ok {
			var tp float64
			fmt.Sscanf(triggerPrice, "%f", &tp)
			req.TriggerPrice = tp
		}

		if disclosedQty, ok := params["disclosed_quantity"].(string); ok {
			req.DisclosedQuantity = disclosedQty
		}
	}

	return c.makeRequest("POST", "/api/v1/placeorder", req)
}

func (c *Client) PlaceSmartOrder(strategy, symbol, action, exchange, priceType, product string, quantity interface{}, positionSize int, optionalParams ...map[string]interface{}) (map[string]interface{}, error) {
	req := PlaceSmartOrderRequest{
		Strategy:     strategy,
		Symbol:       symbol,
		Action:       action,
		Exchange:     exchange,
		PriceType:    priceType,
		Product:      product,
		PositionSize: positionSize,
	}

	switch v := quantity.(type) {
	case string:
		req.Quantity = v
	case int:
		req.Quantity = fmt.Sprintf("%d", v)
	case float64:
		req.Quantity = fmt.Sprintf("%.0f", v)
	default:
		return nil, fmt.Errorf("quantity must be string, int, or float64")
	}

	if len(optionalParams) > 0 {
		params := optionalParams[0]
		if price, ok := params["price"].(float64); ok {
			req.Price = price
		}
		if triggerPrice, ok := params["trigger_price"].(float64); ok {
			req.TriggerPrice = triggerPrice
		}
		if disclosedQty, ok := params["disclosed_quantity"].(string); ok {
			req.DisclosedQuantity = disclosedQty
		}
	}

	return c.makeRequest("POST", "/api/v1/placesmartorder", req)
}

func (c *Client) BasketOrder(orders []BasketOrderItem) (map[string]interface{}, error) {
	req := BasketOrderRequest{
		Orders: orders,
	}
	return c.makeRequest("POST", "/api/v1/basketorder", req)
}

func (c *Client) SplitOrder(symbol, exchange, action string, quantity, splitSize int, priceType, product string, optionalParams ...map[string]interface{}) (map[string]interface{}, error) {
	req := SplitOrderRequest{
		Symbol:    symbol,
		Exchange:  exchange,
		Action:    action,
		Quantity:  quantity,
		SplitSize: splitSize,
		PriceType: priceType,
		Product:   product,
	}

	if len(optionalParams) > 0 {
		params := optionalParams[0]
		if price, ok := params["price"].(float64); ok {
			req.Price = price
		}
		if triggerPrice, ok := params["trigger_price"].(float64); ok {
			req.TriggerPrice = triggerPrice
		}
	}

	return c.makeRequest("POST", "/api/v1/splitorder", req)
}

func (c *Client) ModifyOrder(orderID, strategy, symbol, action, exchange, priceType, product string, quantity interface{}, optionalParams ...map[string]interface{}) (map[string]interface{}, error) {
	req := ModifyOrderRequest{
		OrderID:   orderID,
		Strategy:  strategy,
		Symbol:    symbol,
		Action:    action,
		Exchange:  exchange,
		PriceType: priceType,
		Product:   product,
	}

	switch v := quantity.(type) {
	case string:
		req.Quantity = v
	case int:
		req.Quantity = fmt.Sprintf("%d", v)
	case float64:
		req.Quantity = fmt.Sprintf("%.0f", v)
	default:
		return nil, fmt.Errorf("quantity must be string, int, or float64")
	}

	if len(optionalParams) > 0 {
		params := optionalParams[0]
		if price, ok := params["price"].(float64); ok {
			req.Price = price
		}
		if triggerPrice, ok := params["trigger_price"].(float64); ok {
			req.TriggerPrice = triggerPrice
		}
		if disclosedQty, ok := params["disclosed_quantity"].(string); ok {
			req.DisclosedQuantity = disclosedQty
		}
	}

	return c.makeRequest("POST", "/api/v1/modifyorder", req)
}

func (c *Client) CancelOrder(orderID, strategy string) (map[string]interface{}, error) {
	req := CancelOrderRequest{
		OrderID:  orderID,
		Strategy: strategy,
	}
	return c.makeRequest("POST", "/api/v1/cancelorder", req)
}

func (c *Client) CancelAllOrder(strategy string) (map[string]interface{}, error) {
	req := CancelAllOrderRequest{
		Strategy: strategy,
	}
	return c.makeRequest("POST", "/api/v1/cancelallorder", req)
}

func (c *Client) ClosePosition(strategy string) (map[string]interface{}, error) {
	req := ClosePositionRequest{
		Strategy: strategy,
	}
	return c.makeRequest("POST", "/api/v1/closeposition", req)
}

func (c *Client) OrderStatus(orderID, strategy string) (map[string]interface{}, error) {
	req := OrderStatusRequest{
		OrderID:  orderID,
		Strategy: strategy,
	}
	return c.makeRequest("POST", "/api/v1/orderstatus", req)
}

func (c *Client) OpenPosition(strategy, symbol, exchange, product string) (map[string]interface{}, error) {
	req := OpenPositionRequest{
		Strategy: strategy,
		Symbol:   symbol,
		Exchange: exchange,
		Product:  product,
	}
	return c.makeRequest("POST", "/api/v1/openposition", req)
}