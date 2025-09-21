# OpenAlgo Go SDK

Go SDK for OpenAlgo trading platform.

## Installation

```bash
go get github.com/marketcalls/openalgo-go
```

## Getting Started

First, import the OpenAlgo package and initialize it with your API key:

```go
import "github.com/marketcalls/openalgo-go/openalgo"

// Replace 'your_api_key_here' with your actual API key
// Specify the host URL with your hosted domain or ngrok domain
// If running locally in windows then use the default host value
client := openalgo.NewClient("your_api_key_here", "http://127.0.0.1:5000")
```

## Check OpenAlgo Version

```go
import "github.com/marketcalls/openalgo-go/openalgo"
fmt.Printf("Version: %s\n", openalgo.Version)
```

## API Functions

### Order Management

- `PlaceOrder` - Place a new order
- `PlaceSmartOrder` - Place a smart order considering position size
- `BasketOrder` - Place multiple orders at once
- `SplitOrder` - Split a large order into smaller chunks
- `ModifyOrder` - Modify an existing order
- `CancelOrder` - Cancel a specific order
- `CancelAllOrder` - Cancel all pending orders
- `ClosePosition` - Close all open positions
- `OrderStatus` - Get status of a specific order
- `OpenPosition` - Get open position for a symbol

### Market Data

- `Quotes` - Get real-time quotes
- `Depth` - Get market depth (order book)
- `History` - Get historical data
- `Intervals` - Get available time intervals
- `Symbol` - Get symbol details
- `Search` - Search for symbols
- `Expiry` - Get expiry dates for derivatives

### Account Information

- `Funds` - Get account funds
- `OrderBook` - Get all orders
- `TradeBook` - Get all trades
- `PositionBook` - Get all positions
- `Holdings` - Get holdings

### Analyzer

- `AnalyzerStatus` - Get analyzer status
- `AnalyzerToggle` - Toggle analyzer mode

### WebSocket Streaming

- `Connect` - Connect to WebSocket
- `Disconnect` - Disconnect from WebSocket
- `SubscribeLTP` - Subscribe to LTP updates
- `UnsubscribeLTP` - Unsubscribe from LTP
- `SubscribeQuote` - Subscribe to quote updates
- `UnsubscribeQuote` - Unsubscribe from quotes
- `SubscribeDepth` - Subscribe to market depth
- `UnsubscribeDepth` - Unsubscribe from depth

## Examples

See the `examples/main.go` file for detailed usage examples of all functions.

## Function Parameters

All functions match the Python SDK exactly with the same mandatory and optional parameters. Optional parameters are passed as a map[string]interface{} in Go.

### PlaceOrder Parameters

**Mandatory:**
- strategy (string)
- symbol (string)
- action (string) - BUY/SELL
- exchange (string) - NSE/BSE/NFO/MCX/CDS
- price_type (string) - MARKET/LIMIT/SL/SL-M
- product (string) - MIS/CNC/NRML
- quantity (string/int/float64)

**Optional:**
- price (float64) - Required for LIMIT orders
- trigger_price (float64) - Required for SL orders
- disclosed_quantity (string)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.