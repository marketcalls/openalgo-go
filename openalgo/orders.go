package openalgo

import (
	"fmt"
	"net/http"
)

// PlaceOrder places a new order
func (c *Client) PlaceOrder(strategy, symbol, action, exchange, priceType, product string, quantity interface{}, optionalParams ...map[string]interface{}) (map[string]interface{}, error) {
	// Set defaults
	if strategy == "" {
		strategy = "GO Strategy"
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

	return c.makeRequest(http.MethodPost, "placeorder", payload)
}

// PlaceSmartOrder places a smart order considering position size
func (c *Client) PlaceSmartOrder(strategy, symbol, action, exchange, priceType, product string, quantity interface{}, positionSize interface{}, optionalParams ...map[string]interface{}) (map[string]interface{}, error) {
	// Set defaults
	if strategy == "" {
		strategy = "GO Strategy"
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

	return c.makeRequest(http.MethodPost, "placesmartorder", payload)
}

// BasketOrder places multiple orders at once
func (c *Client) BasketOrder(strategy string, orders []map[string]interface{}) (map[string]interface{}, error) {
	if strategy == "" {
		strategy = "GO Strategy"
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

	return c.makeRequest(http.MethodPost, "basketorder", payload)
}

// SplitOrder splits a large order into smaller orders
func (c *Client) SplitOrder(strategy, symbol, exchange, action string, quantity, splitSize interface{}, priceType, product string, optionalParams ...map[string]interface{}) (map[string]interface{}, error) {
	// Set defaults
	if strategy == "" {
		strategy = "GO Strategy"
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

	return c.makeRequest(http.MethodPost, "splitorder", payload)
}

// ModifyOrder modifies an existing order
func (c *Client) ModifyOrder(orderID, strategy, symbol, action, exchange, priceType, product string, quantity interface{}, price, disclosedQuantity, triggerPrice string) (map[string]interface{}, error) {
	// Set defaults
	if strategy == "" {
		strategy = "GO Strategy"
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
		"apikey":             c.apiKey,
		"orderid":            orderID,
		"strategy":           strategy,
		"symbol":             symbol,
		"action":             action,
		"exchange":           exchange,
		"pricetype":          priceType,
		"product":            product,
		"price":              price,
		"disclosed_quantity": disclosedQuantity,
		"trigger_price":      triggerPrice,
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

	return c.makeRequest(http.MethodPost, "modifyorder", payload)
}

// CancelOrder cancels an existing order
func (c *Client) CancelOrder(orderID, strategy string) (map[string]interface{}, error) {
	if strategy == "" {
		strategy = "GO Strategy"
	}

	payload := map[string]interface{}{
		"apikey":   c.apiKey,
		"orderid":  orderID,
		"strategy": strategy,
	}

	return c.makeRequest(http.MethodPost, "cancelorder", payload)
}

// CancelAllOrder cancels all orders for a strategy
func (c *Client) CancelAllOrder(strategy string) (map[string]interface{}, error) {
	if strategy == "" {
		strategy = "GO Strategy"
	}

	payload := map[string]interface{}{
		"apikey":   c.apiKey,
		"strategy": strategy,
	}

	return c.makeRequest(http.MethodPost, "cancelallorder", payload)
}

// ClosePosition closes all open positions for a strategy
func (c *Client) ClosePosition(strategy string) (map[string]interface{}, error) {
	if strategy == "" {
		strategy = "GO Strategy"
	}

	payload := map[string]interface{}{
		"apikey":   c.apiKey,
		"strategy": strategy,
	}

	return c.makeRequest(http.MethodPost, "closeposition", payload)
}

// OrderStatus gets the status of an order
func (c *Client) OrderStatus(orderID, strategy string) (map[string]interface{}, error) {
	if strategy == "" {
		strategy = "GO Strategy"
	}

	payload := map[string]interface{}{
		"apikey":   c.apiKey,
		"strategy": strategy,
		"orderid":  orderID,
	}

	return c.makeRequest(http.MethodPost, "orderstatus", payload)
}

// OpenPosition gets the open position for a symbol
func (c *Client) OpenPosition(strategy, symbol, exchange, product string) (map[string]interface{}, error) {
	if strategy == "" {
		strategy = "GO Strategy"
	}

	payload := map[string]interface{}{
		"apikey":   c.apiKey,
		"strategy": strategy,
		"symbol":   symbol,
		"exchange": exchange,
		"product":  product,
	}

	return c.makeRequest(http.MethodPost, "openposition", payload)
}

// OptionsOrder places an options order with automatic strike selection
func (c *Client) OptionsOrder(strategy, underlying, exchange, expiryDate, offset, optionType, action string, quantity interface{}, priceType, product string, splitSize interface{}) (map[string]interface{}, error) {
	if strategy == "" {
		strategy = "GO Strategy"
	}
	if priceType == "" {
		priceType = "MARKET"
	}
	if product == "" {
		product = "NRML"
	}

	payload := map[string]interface{}{
		"apikey":      c.apiKey,
		"strategy":    strategy,
		"underlying":  underlying,
		"exchange":    exchange,
		"expiry_date": expiryDate,
		"offset":      offset,
		"option_type": optionType,
		"action":      action,
		"pricetype":   priceType,
		"product":     product,
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

	// Convert splitsize to string
	switch v := splitSize.(type) {
	case string:
		payload["splitsize"] = v
	case int:
		payload["splitsize"] = fmt.Sprintf("%d", v)
	case float64:
		payload["splitsize"] = fmt.Sprintf("%.0f", v)
	default:
		payload["splitsize"] = "0"
	}

	return c.makeRequest(http.MethodPost, "optionsorder", payload)
}

// OptionsLeg represents a single leg in a multi-leg options order
type OptionsLeg struct {
	Offset     string `json:"offset"`
	OptionType string `json:"option_type"`
	Action     string `json:"action"`
	Quantity   string `json:"quantity"`
	ExpiryDate string `json:"expiry_date,omitempty"`
}

// OptionsMultiOrder places a multi-leg options order (e.g., Iron Condor, Spreads)
func (c *Client) OptionsMultiOrder(strategy, underlying, exchange, expiryDate string, legs []OptionsLeg) (map[string]interface{}, error) {
	if strategy == "" {
		strategy = "GO Strategy"
	}

	// Convert legs to interface slice
	legsPayload := make([]map[string]interface{}, len(legs))
	for i, leg := range legs {
		legMap := map[string]interface{}{
			"offset":      leg.Offset,
			"option_type": leg.OptionType,
			"action":      leg.Action,
			"quantity":    leg.Quantity,
		}
		if leg.ExpiryDate != "" {
			legMap["expiry_date"] = leg.ExpiryDate
		}
		legsPayload[i] = legMap
	}

	payload := map[string]interface{}{
		"apikey":     c.apiKey,
		"strategy":   strategy,
		"underlying": underlying,
		"exchange":   exchange,
		"legs":       legsPayload,
	}

	if expiryDate != "" {
		payload["expiry_date"] = expiryDate
	}

	return c.makeRequest(http.MethodPost, "optionsmultiorder", payload)
}
