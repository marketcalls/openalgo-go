# OpenAlgo Go SDK

Official Go SDK for [OpenAlgo](https://openalgo.in) - the open source algorithmic trading platform.

## Installation

### Method 1: Install in your project

```bash
# Create a new project directory
mkdir my-trading-app
cd my-trading-app

# Initialize Go module
go mod init my-trading-app

# Install OpenAlgo Go SDK
go get github.com/marketcalls/openalgo-go

# Clean up dependencies
go mod tidy
```

### Method 2: Clone from GitHub

```bash
# Clone the repository
git clone https://github.com/marketcalls/openalgo-go.git
cd openalgo-go

# Install dependencies
go mod download

# Clean up dependencies
go mod tidy

# Run the example
go run example.go
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "github.com/marketcalls/openalgo-go/openalgo"
)

func main() {
    // Initialize the client
    client := openalgo.NewClient(
        "YOUR_API_KEY",          // Your OpenAlgo API key
        "http://127.0.0.1:5000", // OpenAlgo server URL
        "v1",                    // API version
        "ws://127.0.0.1:8765",   // WebSocket URL (optional)
    )

    // Fetch account funds
    funds, err := client.Funds()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Funds: %+v\n", funds)
}
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

## Running the Example

1. Update `example.go` with your API key:
```go
client := openalgo.NewClient(
    "YOUR_API_KEY",          // Replace with your actual API key
    "http://127.0.0.1:5000", // Your OpenAlgo server URL
    "v1",                    // API version
    "ws://127.0.0.1:8765",   // WebSocket URL (optional)
)
```

2. Run the example:
```bash
go run example.go
```

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