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
- `OptionsOrder` - Place an options order with automatic strike selection (optional `strike_int`, `price`, `trigger_price`, `disclosed_quantity`)
- `OptionsMultiOrder` - Place a multi-leg options order (Iron Condor, spreads, etc.)

### Market Data

- `Quotes` - Get real-time quotes
- `Depth` - Get market depth (order book)
- `History` - Get historical data (optional `source`: `"api"` default or `"db"`)
- `Intervals` - Get available time intervals
- `Symbol` - Get symbol details
- `Search` - Search for symbols
- `Expiry` - Get expiry dates for derivatives
- `Instruments` - Download instruments for an exchange, or **all exchanges** when called with `""`
- `MultiQuotes` - Get quotes for multiple symbols in one call
- `OptionChain` - Get the full option chain for an underlying/expiry
- `OptionSymbol` - Resolve an option symbol from underlying + strike offset (optional `strike_int`)
- `OptionGreeks` - Calculate Delta/Gamma/Theta/Vega/Rho and IV (all params beyond `symbol`/`exchange` optional: `interest_rate`, `forward_price`, `underlying_symbol`, `underlying_exchange`, `expiry_time`)
- `SyntheticFuture` - Calculate synthetic future price from ATM options

### Account Information

- `Funds` - Get account funds
- `OrderBook` - Get all orders
- `TradeBook` - Get all trades
- `PositionBook` - Get all positions
- `Holdings` - Get holdings
- `Margin` - Calculate margin requirements for a basket of positions (optional `Price`/`TriggerPrice` per position, default `"0"`)

### Analyzer

- `AnalyzerStatus` - Get analyzer status
- `AnalyzerToggle` - Toggle analyzer mode

### Utilities

- `Ping` - Check API connectivity
- `Holidays` - Get market holidays for a year; omit the argument to default to the current year
- `Timings` - Get exchange trading timings for a date; omit the argument to default to today

### Notifications

- `Telegram` / `TelegramWithPriority` - Send a Telegram notification
- `WhatsApp` - Send a WhatsApp message (text/image/document) via the OpenAlgo paired device

### Strategy Webhook

- `NewStrategy` - Create a standalone TradingView-style webhook client (host URL + webhook ID, no API key needed)
- `StrategyOrder` - Send a BUY/SELL signal (with optional position size) to the strategy's webhook

### WebSocket Streaming

- `Connect` - Connect to WebSocket
- `Disconnect` - Disconnect from WebSocket
- `SubscribeLTP` - Subscribe to LTP updates
- `UnsubscribeLTP` - Unsubscribe from LTP
- `SubscribeQuote` - Subscribe to quote updates
- `UnsubscribeQuote` - Unsubscribe from quotes
- `SubscribeDepth` - Subscribe to market depth
- `UnsubscribeDepth` - Unsubscribe from depth
- `GetLTP(exchange, symbol string)` - Read the latest cached LTP snapshot(s); pass `""` for either argument to skip that filter
- `GetQuotes(exchange, symbol string)` - Read the latest cached Quote snapshot(s), same filtering rules
- `GetDepth(exchange, symbol string)` - Read the latest cached Market Depth snapshot(s) (5-level `buyBook`/`sellBook`), same filtering rules

  These three read from a local cache that is populated automatically as WebSocket
  market data messages arrive (mirroring the Python SDK's `ltp_data` / `quotes_data`
  / `depth_data`) - no callback wiring required if you just want the latest snapshot.

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

### OptionGreeks Parameters

```go
greeks, err := client.OptionGreeks("NIFTY02DEC2526000CE", "NFO", map[string]interface{}{
    "interest_rate":       6.5,   // optional, defaults to 0 server-side
    "forward_price":       26350, // optional, skips underlying LTP fetch
    "underlying_symbol":   "NIFTY30DEC25FUT",
    "underlying_exchange": "NFO",
    "expiry_time":         "19:00", // optional, needed for MCX
})
```

All parameters after `symbol`/`exchange` are optional - omit the map entirely for the simplest case: `client.OptionGreeks("NIFTY02DEC2526000CE", "NFO")`.

### OptionsOrder / OptionSymbol Optional Parameters

```go
// strike_int is DEPRECATED (matches the Python SDK) but still supported
symbolInfo, err := client.OptionSymbol("NIFTY", "NSE_INDEX", "28OCT25", "ATM", "CE", map[string]interface{}{
    "strike_int": 50,
})

orderResp, err := client.OptionsOrder(
    "test_strategy", "NIFTY", "NSE_INDEX", "28NOV24", "ATM", "CE", "BUY",
    75, "LIMIT", "MIS", 0,
    map[string]interface{}{
        "strike_int": 50,
        "price":      "50.0",
    },
)
```

### Margin Parameters

`Price` and `TriggerPrice` on `MarginPosition` are optional and default to `"0"` (set `Price` for LIMIT orders, `TriggerPrice` for SL/SL-M orders):

```go
resp, err := client.Margin([]openalgo.MarginPosition{
    {
        Symbol: "NIFTY30DEC2526000CE", Exchange: "NFO", Action: "SELL",
        Product: "NRML", PriceType: "LIMIT", Quantity: "75",
        Price: "150.75",
    },
})
```

### History and Instruments

```go
// source is optional: "api" (default) or "db" (OpenAlgo DuckDB/Historify)
hist, err := client.History("SBIN", "NSE", "5m", "2025-04-01", "2025-04-08", "db")

// Instruments("") downloads and combines instruments from every supported exchange
allInstruments, err := client.Instruments("")
nseOnly, err := client.Instruments("NSE")
```

### Holidays and Timings

```go
// Both default client-side when omitted: current year / today
holidays, err := client.Holidays()      // current year
holidays2025, err := client.Holidays(2025)

timings, err := client.Timings()        // today
timingsOn, err := client.Timings("2025-12-25")
```

### WhatsApp Example

```go
// Send to self (the paired device's own number)
client.WhatsApp("Build #482 deployed. P&L: +1.2%")

// Send to a single number with an image
client.WhatsApp("NIFTY end-of-day chart", map[string]interface{}{
    "to":    "919876543210",
    "image": "/srv/charts/nifty_eod.png",
})

// Fire-and-forget broadcast (up to 5 numbers) for time-critical alerts
client.WhatsApp("Stop-loss hit on BANKNIFTY!", map[string]interface{}{
    "to":                []string{"919876543210", "919812345678"},
    "wait_for_delivery": false,
})
```

### Strategy Webhook Example

`Strategy` is a standalone client (separate from `Client`) for posting TradingView-style
signals to an OpenAlgo strategy webhook - it only needs the server host and the
strategy's webhook ID, not an API key:

```go
strategy := openalgo.NewStrategy("http://127.0.0.1:5000", "your-webhook-id")

// The strategy mode (LONG_ONLY, SHORT_ONLY, BOTH) is configured in OpenAlgo.
// position_size is required for BOTH mode; pass nil to omit it otherwise.
resp, err := strategy.StrategyOrder("RELIANCE", "BUY", nil)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Strategy order response: %+v\n", resp)
```

### WebSocket Snapshot Getters Example

```go
client.Connect()
client.SubscribeLTP([]openalgo.Instrument{{Exchange: "NSE", Symbol: "RELIANCE"}}, nil)
client.SubscribeQuote([]openalgo.Instrument{{Exchange: "NSE", Symbol: "RELIANCE"}}, nil)
client.SubscribeDepth([]openalgo.Instrument{{Exchange: "NSE", Symbol: "RELIANCE"}}, nil)

// Later, read the latest cached snapshots without wiring up a callback:
ltpSnapshot := client.GetLTP("NSE", "RELIANCE") // {"ltp": {"NSE": {"RELIANCE": {"timestamp": ..., "ltp": ...}}}}
allLTP := client.GetLTP("", "")                 // every cached LTP symbol

quoteSnapshot := client.GetQuotes("NSE", "RELIANCE") // {"quote": {"NSE": {"RELIANCE": {"open": ..., "high": ..., "low": ..., "close": ..., "ltp": ..., "volume": ...}}}}
allQuotes := client.GetQuotes("", "")                // every cached quote symbol

depthSnapshot := client.GetDepth("NSE", "RELIANCE") // {"depth": {"NSE": {"RELIANCE": {"ltp": ..., "buyBook": {"1": {...}, ...}, "sellBook": {"1": {...}, ...}}}}}
allDepth := client.GetDepth("", "")                 // every cached depth symbol
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.