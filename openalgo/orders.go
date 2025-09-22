package openalgo

import (
	"fmt"
)

// PlaceOrder places a new order
func (c *Client) PlaceOrder(strategy, symbol, action, exchange, priceType, product string, quantity interface{}, optionalParams ...map[string]interface{}) (map[string]interface{}, error) {
	// Set defaults
	if strategy == "" {
		strategy = "Python"
	}
	if priceType == "" {
		priceType = "MARKET"
	}
	if product == "" {
		product = "MIS"
	}

	payload := map[string]interface{}{
		"apikey":    c.apiKey,
		"strategy":  strategy,
		"symbol":    symbol,
		"action":    action,
		"exchange":  exchange,
		"pricetype": priceType,
		"product":   product,
	}

	// Convert quantity to string
	switch v := quantity.(type) {
	case string:
		payload["quantity"] = v
	case int:
		payload["quantity"] = fmt.Sprintf("%d", v)
	case float64:
		payload["quantity"] = fmt.Sprintf("%.0f", v)
	default:
		payload["quantity"] = "1"
	}

	// Add optional parameters
	if len(optionalParams) > 0 {
		params := optionalParams[0]
		for key, value := range params {
			if value != nil {
				// Convert all values to strings
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
	}

	return c.makeRequest("POST", "placeorder", payload)
}

// PlaceSmartOrder places a smart order considering position size
func (c *Client) PlaceSmartOrder(strategy, symbol, action, exchange, priceType, product string, quantity interface{}, positionSize interface{}, optionalParams ...map[string]interface{}) (map[string]interface{}, error) {
	// Set defaults
	if strategy == "" {
		strategy = "Python"
	}
	if priceType == "" {
		priceType = "MARKET"
	}
	if product == "" {
		product = "MIS"
	}

	payload := map[string]interface{}{
		"apikey":    c.apiKey,
		"strategy":  strategy,
		"symbol":    symbol,
		"action":    action,
		"exchange":  exchange,
		"pricetype": priceType,
		"product":   product,
	}

	// Convert quantity to string
	switch v := quantity.(type) {
	case string:
		payload["quantity"] = v
	case int:
		payload["quantity"] = fmt.Sprintf("%d", v)
	case float64:
		payload["quantity"] = fmt.Sprintf("%.0f", v)
	default:
		payload["quantity"] = "1"
	}

	// Convert position_size to string
	switch v := positionSize.(type) {
	case string:
		payload["position_size"] = v
	case int:
		payload["position_size"] = fmt.Sprintf("%d", v)
	case float64:
		payload["position_size"] = fmt.Sprintf("%.0f", v)
	default:
		return nil, fmt.Errorf("position_size is required")
	}

	// Add optional parameters
	if len(optionalParams) > 0 {
		params := optionalParams[0]
		for key, value := range params {
			if value != nil {
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
	}

	return c.makeRequest("POST", "placesmartorder", payload)
}

// BasketOrder places multiple orders at once
func (c *Client) BasketOrder(strategy string, orders []map[string]interface{}) (map[string]interface{}, error) {
	if strategy == "" {
		strategy = "Python"
	}

	// Process orders to ensure all numeric values are strings
	processedOrders := make([]map[string]interface{}, len(orders))
	for i, order := range orders {
		processedOrder := make(map[string]interface{})
		for key, value := range order {
			switch v := value.(type) {
			case int:
				processedOrder[key] = fmt.Sprintf("%d", v)
			case float64:
				processedOrder[key] = fmt.Sprintf("%g", v)
			default:
				processedOrder[key] = v
			}
		}
		processedOrders[i] = processedOrder
	}

	payload := map[string]interface{}{
		"apikey":   c.apiKey,
		"strategy": strategy,
		"orders":   processedOrders,
	}

	return c.makeRequest("POST", "basketorder", payload)
}

// SplitOrder splits a large order into smaller orders
func (c *Client) SplitOrder(strategy, symbol, exchange, action string, quantity, splitSize interface{}, priceType, product string, optionalParams ...map[string]interface{}) (map[string]interface{}, error) {
	// Set defaults
	if strategy == "" {
		strategy = "Python"
	}
	if priceType == "" {
		priceType = "MARKET"
	}
	if product == "" {
		product = "MIS"
	}

	payload := map[string]interface{}{
		"apikey":    c.apiKey,
		"strategy":  strategy,
		"symbol":    symbol,
		"action":    action,
		"exchange":  exchange,
		"pricetype": priceType,
		"product":   product,
	}

	// Convert quantity to string
	switch v := quantity.(type) {
	case string:
		payload["quantity"] = v
	case int:
		payload["quantity"] = fmt.Sprintf("%d", v)
	case float64:
		payload["quantity"] = fmt.Sprintf("%.0f", v)
	default:
		return nil, fmt.Errorf("quantity is required")
	}

	// Convert splitsize to string
	switch v := splitSize.(type) {
	case string:
		payload["splitsize"] = v
	case int:
		payload["splitsize"] = fmt.Sprintf("%d", v)
	case float64:
		payload["splitsize"] = fmt.Sprintf("%.0f", v)
	default:
		return nil, fmt.Errorf("splitsize is required")
	}

	// Add optional parameters
	if len(optionalParams) > 0 {
		params := optionalParams[0]
		for key, value := range params {
			if value != nil {
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
	}

	return c.makeRequest("POST", "splitorder", payload)
}

// ModifyOrder modifies an existing order
func (c *Client) ModifyOrder(orderID, strategy, symbol, action, exchange, priceType, product string, quantity interface{}, price, disclosedQuantity, triggerPrice string) (map[string]interface{}, error) {
	// Set defaults
	if strategy == "" {
		strategy = "Python"
	}
	if priceType == "" {
		priceType = "LIMIT"
	}
	if disclosedQuantity == "" {
		disclosedQuantity = "0"
	}
	if triggerPrice == "" {
		triggerPrice = "0"
	}

	payload := map[string]interface{}{
		"apikey":              c.apiKey,
		"orderid":             orderID,
		"strategy":            strategy,
		"symbol":              symbol,
		"action":              action,
		"exchange":            exchange,
		"pricetype":           priceType,
		"product":             product,
		"price":               price,
		"disclosed_quantity":  disclosedQuantity,
		"trigger_price":       triggerPrice,
	}

	// Convert quantity to string
	switch v := quantity.(type) {
	case string:
		payload["quantity"] = v
	case int:
		payload["quantity"] = fmt.Sprintf("%d", v)
	case float64:
		payload["quantity"] = fmt.Sprintf("%.0f", v)
	default:
		return nil, fmt.Errorf("quantity is required")
	}

	return c.makeRequest("POST", "modifyorder", payload)
}

// CancelOrder cancels an existing order
func (c *Client) CancelOrder(orderID, strategy string) (map[string]interface{}, error) {
	if strategy == "" {
		strategy = "Python"
	}

	payload := map[string]interface{}{
		"apikey":   c.apiKey,
		"orderid":  orderID,
		"strategy": strategy,
	}

	return c.makeRequest("POST", "cancelorder", payload)
}

// CancelAllOrder cancels all orders for a strategy
func (c *Client) CancelAllOrder(strategy string) (map[string]interface{}, error) {
	if strategy == "" {
		strategy = "Python"
	}

	payload := map[string]interface{}{
		"apikey":   c.apiKey,
		"strategy": strategy,
	}

	return c.makeRequest("POST", "cancelallorder", payload)
}

// ClosePosition closes all open positions for a strategy
func (c *Client) ClosePosition(strategy string) (map[string]interface{}, error) {
	if strategy == "" {
		strategy = "Python"
	}

	payload := map[string]interface{}{
		"apikey":   c.apiKey,
		"strategy": strategy,
	}

	return c.makeRequest("POST", "closeposition", payload)
}

// OrderStatus gets the status of an order
func (c *Client) OrderStatus(orderID, strategy string) (map[string]interface{}, error) {
	if strategy == "" {
		strategy = "Python"
	}

	payload := map[string]interface{}{
		"apikey":   c.apiKey,
		"strategy": strategy,
		"orderid":  orderID,
	}

	return c.makeRequest("POST", "orderstatus", payload)
}

// OpenPosition gets the open position for a symbol
func (c *Client) OpenPosition(strategy, symbol, exchange, product string) (map[string]interface{}, error) {
	if strategy == "" {
		strategy = "Python"
	}

	payload := map[string]interface{}{
		"apikey":   c.apiKey,
		"strategy": strategy,
		"symbol":   symbol,
		"exchange": exchange,
		"product":  product,
	}

	return c.makeRequest("POST", "openposition", payload)
}