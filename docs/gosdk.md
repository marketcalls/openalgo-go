# Go SDK for OpenAlgo

To install the OpenAlgo Go library, use go get:

```bash
go get github.com/marketcalls/openalgo-go
```

### Get the OpenAlgo apikey

Make Sure that your OpenAlgo Application is running. Login to OpenAlgo Application with valid credentials and get the OpenAlgo apikey

For detailed function parameters refer to the [API Documentation](https://docs.openalgo.in/api-documentation/v1)

### Getting Started with OpenAlgo

First, import the OpenAlgo package and initialize it with your API key:

```go
package main

import (
    "github.com/marketcalls/openalgo-go/openalgo"
)

// Replace 'your_api_key_here' with your actual API key
// Specify the host URL with your hosted domain or ngrok domain.
// If running locally then use the default host value.
// Parameters: apiKey, host, version, websocketURL (optional)
client := openalgo.NewClient(
    "your_api_key_here",
    "http://127.0.0.1:5000",
    "v1",                    // API version
    "ws://127.0.0.1:8765",   // WebSocket URL (optional)
)

```

### Check OpenAlgo Version

```go
import (
    "fmt"
    "github.com/marketcalls/openalgo-go/openalgo"
)

fmt.Printf("OpenAlgo Go SDK Version: %s\n", openalgo.Version)
```

### Complete List of Implemented Functions

#### Order Management
- `PlaceOrder` - Place a new order
- `PlaceSmartOrder` - Place a smart order with position sizing
- `BasketOrder` - Place multiple orders at once
- `SplitOrder` - Split large orders into smaller chunks
- `ModifyOrder` - Modify an existing order
- `CancelOrder` - Cancel a specific order
- `CancelAllOrder` - Cancel all pending orders
- `ClosePosition` - Close all open positions
- `OrderStatus` - Get order status
- `OpenPosition` - Get open position for a symbol

#### Market Data
- `Quotes` - Get real-time quotes
- `Depth` - Get market depth
- `History` - Get historical data
- `Intervals` - Get available time intervals
- `Symbol` - Get symbol information
- `Search` - Search for symbols
- `Expiry` - Get expiry dates

#### Account Information
- `Funds` - Get account funds
- `OrderBook` - Get all orders
- `TradeBook` - Get all trades
- `PositionBook` - Get all positions
- `Holdings` - Get holdings

#### Utility
- `Ping` - Check API connectivity
- `AnalyzerStatus` - Get analyzer status
- `AnalyzerToggle` - Toggle analyzer mode (requires boolean: true to enable, false to disable)

#### WebSocket Streaming
- `Connect` - Connect to WebSocket
- `Disconnect` - Disconnect from WebSocket
- `SubscribeLTP` - Subscribe to LTP updates
- `UnsubscribeLTP` - Unsubscribe from LTP
- `SubscribeQuote` - Subscribe to quote updates
- `UnsubscribeQuote` - Unsubscribe from quotes
- `SubscribeDepth` - Subscribe to market depth
- `UnsubscribeDepth` - Unsubscribe from depth

### Examples

Please refer to the documentation on [order constants](https://docs.openalgo.in/api-documentation/v1/order-constants), and consult the API reference for details on optional parameters

### Ping Example

Check API connectivity:

```go
response, err := client.Ping()
if err != nil {
    log.Printf("Error: %v", err)
}
fmt.Printf("%v\n", response)
```

**Ping Response**

```json
{"status": "success", "message": "pong"}
```

### PlaceOrder Example

To place a new market order:

```go
response, err := client.PlaceOrder(
    "GO Strategy", // strategy
    "NHPC",      // symbol
    "BUY",       // action
    "NSE",       // exchange
    "MARKET",    // price_type
    "MIS",       // product
    1,           // quantity
)
if err != nil {
    log.Printf("Error: %v", err)
}
fmt.Printf("%v\n", response)
```

**Place Market Order Response**

```json
{"orderid": "250408000989443", "status": "success"}
```

To place a new limit order:

```go
response, err := client.PlaceOrder(
    "GO Strategy", // strategy
    "YESBANK",   // symbol
    "BUY",       // action
    "NSE",       // exchange
    "LIMIT",     // price_type
    "MIS",       // product
    "1",         // quantity
    map[string]interface{}{
        "price":              16.0,
        "trigger_price":      0.0,
        "disclosed_quantity": "0",
    },
)
if err != nil {
    log.Printf("Error: %v", err)
}
fmt.Printf("%v\n", response)
```

**Place Limit Order Response**

```json
{"orderid": "250408001003813", "status": "success"}
```

### PlaceSmartOrder Example

To place a smart order considering the current position size:

```go
response, err := client.PlaceSmartOrder(
    "GO Strategy", // strategy
    "TATAMOTORS",  // symbol
    "SELL",        // action
    "NSE",         // exchange
    "MARKET",      // price_type
    "MIS",         // product
    1,             // quantity
    5,             // position_size
)
if err != nil {
    log.Printf("Error: %v", err)
}
fmt.Printf("%v\n", response)
```

**Place Smart Market Order Response**

```json
{"orderid": "250408000997543", "status": "success"}
```

### BasketOrder Example

To place a new basket order:

```go
basketOrders := []map[string]interface{}{
    {
        "symbol":    "BHEL",
        "exchange":  "NSE",
        "action":    "BUY",
        "quantity":  1,
        "pricetype": "MARKET",
        "product":   "MIS",
    },
    {
        "symbol":    "ZOMATO",
        "exchange":  "NSE",
        "action":    "SELL",
        "quantity":  1,
        "pricetype": "MARKET",
        "product":   "MIS",
    },
}
response, err := client.BasketOrder("GO Strategy", basketOrders)
if err != nil {
    log.Printf("Error: %v", err)
}
fmt.Printf("%v\n", response)
```

**Basket Order Response**

```json
{
  "status": "success",
  "results": [
    {
      "symbol": "BHEL",
      "status": "success",
      "orderid": "250408000999544"
    },
    {
      "symbol": "ZOMATO",
      "status": "success",
      "orderid": "250408000997545"
    }
  ]
}
```

### SplitOrder Example

To place a new split order:

```go
response, err := client.SplitOrder(
    "GO Strategy", // strategy
    "YESBANK",   // symbol
    "NSE",       // exchange
    "SELL",      // action
    105,         // quantity
    20,          // split_size
    "MARKET",    // price_type
    "MIS",       // product
)
if err != nil {
    log.Printf("Error: %v", err)
}
fmt.Printf("%v\n", response)
```

**SplitOrder Response**

```json
{
  "status": "success",
  "split_size": 20,
  "total_quantity": 105,
  "results": [
    {
      "order_num": 1,
      "orderid": "250408001021467",
      "quantity": 20,
      "status": "success"
    },
    {
      "order_num": 2,
      "orderid": "250408001021459",
      "quantity": 20,
      "status": "success"
    },
    {
      "order_num": 3,
      "orderid": "250408001021466",
      "quantity": 20,
      "status": "success"
    },
    {
      "order_num": 4,
      "orderid": "250408001021470",
      "quantity": 20,
      "status": "success"
    },
    {
      "order_num": 5,
      "orderid": "250408001021471",
      "quantity": 20,
      "status": "success"
    },
    {
      "order_num": 6,
      "orderid": "250408001021472",
      "quantity": 5,
      "status": "success"
    }
  ]
}
```

### ModifyOrder Example

To modify an existing order:

```go
response, err := client.ModifyOrder(
    "250408001002736", // order_id
    "GO Strategy",     // strategy
    "YESBANK",         // symbol
    "BUY",             // action
    "NSE",             // exchange
    "LIMIT",           // price_type
    "CNC",             // product
    1,                 // quantity
    "16.5",            // price
    "0",               // disclosed_quantity
    "0",               // trigger_price
)
if err != nil {
    log.Printf("Error: %v", err)
}
fmt.Printf("%v\n", response)
```

**Modify Order Response**

```json
{"orderid": "250408001002736", "status": "success"}
```

### CancelOrder Example

To cancel an existing order:

```go
response, err := client.CancelOrder(
    "250408001002736", // order_id
    "GO Strategy",     // strategy
)
if err != nil {
    log.Printf("Error: %v", err)
}
fmt.Printf("%v\n", response)
```

**Cancelorder Response**

```json
{"orderid": "250408001002736", "status": "success"}
```

### CancelAllOrder Example

To cancel all open orders and trigger pending orders:

```go
response, err := client.CancelAllOrder("GO Strategy")
if err != nil {
    log.Printf("Error: %v", err)
}
fmt.Printf("%v\n", response)
```

**Cancelallorder Response**

```json
{
  "status": "success",
  "message": "Canceled 5 orders. Failed to cancel 0 orders.",
  "canceled_orders": [
    "250408001042620",
    "250408001042667",
    "250408001042642",
    "250408001043015",
    "250408001043386"
  ],
  "failed_cancellations": []
}
```

### ClosePosition Example

To close all open positions across various exchanges:

```go
response, err := client.ClosePosition("GO Strategy")
if err != nil {
    log.Printf("Error: %v", err)
}
fmt.Printf("%v\n", response)
```

**ClosePosition Response**

```json
{"message": "All Open Positions Squared Off", "status": "success"}
```

### OrderStatus Example

To Get the Current OrderStatus:

```go
response, err := client.OrderStatus(
    "250408000989443", // order_id
    "GO Strategy",     // strategy
)
if err != nil {
    log.Printf("Error: %v", err)
}
fmt.Printf("%v\n", response)
```

**Orderstatus Response**

```json
{
  "status": "success",
  "data": {
    "orderid": "250408000989443",
    "symbol": "RELIANCE",
    "exchange": "NSE",
    "action": "BUY",
    "order_status": "complete",
    "quantity": "1",
    "price": 1186.0,
    "pricetype": "MARKET",
    "trigger_price": 0.0,
    "product": "MIS",
    "timestamp": "08-Apr-2025 13:58:03"
  }
}
```

### OpenPosition Example

To Get the Current OpenPosition:

```go
response, err := client.OpenPosition(
    "GO Strategy",   // strategy
    "YESBANK",       // symbol
    "NSE",           // exchange
    "MIS",           // product
)
if err != nil {
    log.Printf("Error: %v", err)
}
fmt.Printf("%v\n", response)
```

**OpenPosition Response**

```json
{"quantity": "-10", "status": "success"}
```

### Quotes Example

```go
response, err := client.Quotes("RELIANCE", "NSE")
if err != nil {
    log.Printf("Error: %v", err)
}
fmt.Printf("%v\n", response)
```

**Quotes response**

```json
{
  "status": "success",
  "data": {
    "open": 1172.0,
    "high": 1196.6,
    "low": 1163.3,
    "ltp": 1187.75,
    "ask": 1188.0,
    "bid": 1187.85,
    "prev_close": 1165.7,
    "volume": 14414545
  }
}
```

### Depth Example

```go
response, err := client.Depth("SBIN", "NSE")
if err != nil {
    log.Printf("Error: %v", err)
}
fmt.Printf("%v\n", response)
```

**Depth Response**

```json
{
  "status": "success",
  "data": {
    "open": 760.0,
    "high": 774.0,
    "low": 758.15,
    "ltp": 769.6,
    "ltq": 205,
    "prev_close": 746.9,
    "volume": 9362799,
    "oi": 161265750,
    "totalbuyqty": 591351,
    "totalsellqty": 835701,
    "asks": [
      {
        "price": 769.6,
        "quantity": 767
      },
      {
        "price": 769.65,
        "quantity": 115
      },
      {
        "price": 769.7,
        "quantity": 162
      },
      {
        "price": 769.75,
        "quantity": 1121
      },
      {
        "price": 769.8,
        "quantity": 430
      }
    ],
    "bids": [
      {
        "price": 769.4,
        "quantity": 886
      },
      {
        "price": 769.35,
        "quantity": 212
      },
      {
        "price": 769.3,
        "quantity": 351
      },
      {
        "price": 769.25,
        "quantity": 343
      },
      {
        "price": 769.2,
        "quantity": 399
      }
    ]
  }
}
```

### History Example

```go
response, err := client.History(
    "SBIN",       // symbol
    "NSE",        // exchange
    "5m",         // interval
    "2025-04-01", // start_date
    "2025-04-08", // end_date
)
if err != nil {
    log.Printf("Error: %v", err)
}
fmt.Printf("%v\n", response)
```

**History Response**

```json
                            close    high     low    open  volume
timestamp
2025-04-01 09:15:00+05:30  772.50  774.00  763.20  766.50  318625
2025-04-01 09:20:00+05:30  773.20  774.95  772.10  772.45  197189
2025-04-01 09:25:00+05:30  775.15  775.60  772.60  773.20  227544
2025-04-01 09:30:00+05:30  777.35  777.50  774.85  775.15  134596
2025-04-01 09:35:00+05:30  778.00  778.00  776.25  777.50  145385
...                           ...     ...     ...     ...     ...
2025-04-08 14:00:00+05:30  768.25  770.70  767.85  768.50  142478
2025-04-08 14:05:00+05:30  769.10  769.80  766.60  768.15  128283
2025-04-08 14:10:00+05:30  769.05  769.85  768.40  769.10  119084
2025-04-08 14:15:00+05:30  770.05  770.50  769.05  769.05  158299
2025-04-08 14:20:00+05:30  769.95  770.50  769.40  770.05  125485

[437 rows x 5 columns]
```

### Intervals Example

```go
response, err := client.Intervals()
if err != nil {
    log.Printf("Error: %v", err)
}
fmt.Printf("%v\n", response)
```

**Intervals response**

```json
{
  "status": "success",
  "data": {
    "months": [],
    "weeks": [],
    "days": ["D"],
    "hours": ["1h"],
    "minutes": ["10m", "15m", "1m", "30m", "3m", "5m"],
    "seconds": []
  }
}
```

### Symbol Example

```go
response, err := client.Symbol("RELIANCE", "NSE")
if err != nil {
    log.Printf("Error: %v", err)
}
fmt.Printf("%v\n", response)
```

**Symbols Response**

```json
{
  "status": "success",
  "data": {
    "id": 979,
    "name": "RELIANCE",
    "symbol": "RELIANCE",
    "brsymbol": "RELIANCE-EQ",
    "exchange": "NSE",
    "brexchange": "NSE",
    "instrumenttype": "",
    "expiry": "",
    "strike": -0.01,
    "lotsize": 1,
    "tick_size": 0.05,
    "token": "2885"
  }
}
```

### Search Example

```go
response, err := client.Search("NIFTY 25000 JUL CE", "NFO")
if err != nil {
    log.Printf("Error: %v", err)
}
fmt.Printf("%v\n", response)
```

**Search Response**

```json
{
  "data": [
    {
      "brexchange": "NFO",
      "brsymbol": "NIFTY17JUL2525000CE",
      "exchange": "NFO",
      "expiry": "17-JUL-25",
      "instrumenttype": "OPTIDX",
      "lotsize": 75,
      "name": "NIFTY",
      "strike": 25000,
      "symbol": "NIFTY17JUL2525000CE",
      "tick_size": 0.05,
      "token": "47275"
    },
    {
      "brexchange": "NFO",
      "brsymbol": "FINNIFTY31JUL2525000CE",
      "exchange": "NFO",
      "expiry": "31-JUL-25",
      "instrumenttype": "OPTIDX",
      "lotsize": 65,
      "name": "FINNIFTY",
      "strike": 25000,
      "symbol": "FINNIFTY31JUL2525000CE",
      "tick_size": 0.05,
      "token": "54763"
    },
    {
      "brexchange": "NFO",
      "brsymbol": "NIFTY24JUL2525000CE",
      "exchange": "NFO",
      "expiry": "24-JUL-25",
      "instrumenttype": "OPTIDX",
      "lotsize": 75,
      "name": "NIFTY",
      "strike": 25000,
      "symbol": "NIFTY24JUL2525000CE",
      "tick_size": 0.05,
      "token": "49487"
    }
  ],
  "message": "Found 6 matching symbols",
  "status": "success"
}
```

### Expiry Example

```go
response, err := client.Expiry(
    "NIFTY",   // symbol
    "NFO",     // exchange
    "options", // instrument_type
)
if err != nil {
    log.Printf("Error: %v", err)
}
fmt.Printf("%v\n", response)
```

**Expiry Response**

```json
{
  "data": [
    "10-JUL-25",
    "17-JUL-25",
    "24-JUL-25",
    "31-JUL-25",
    "07-AUG-25",
    "28-AUG-25",
    "25-SEP-25",
    "24-DEC-25",
    "26-MAR-26",
    "25-JUN-26",
    "31-DEC-26",
    "24-JUN-27",
    "30-DEC-27",
    "29-JUN-28",
    "28-DEC-28",
    "28-JUN-29",
    "27-DEC-29",
    "25-JUN-30"
  ],
  "message": "Found 18 expiry dates for NIFTY options in NFO",
  "status": "success"
}
```

### Funds Example

```go
response, err := client.Funds()
if err != nil {
    log.Printf("Error: %v", err)
}
fmt.Printf("%v\n", response)
```

**Funds Response**

```json
{
  "status": "success",
  "data": {
    "availablecash": "320.66",
    "collateral": "0.00",
    "m2mrealized": "3.27",
    "m2munrealized": "-7.88",
    "utiliseddebits": "679.34"
  }
}
```

### OrderBook Example

```go
response, err := client.OrderBook()
if err != nil {
    log.Printf("Error: %v", err)
}
fmt.Printf("%v\n", response)
```

```json
{
  "status": "success",
  "data": {
    "orders": [
      {
        "action": "BUY",
        "symbol": "RELIANCE",
        "exchange": "NSE",
        "orderid": "250408000989443",
        "product": "MIS",
        "quantity": "1",
        "price": 1186.0,
        "pricetype": "MARKET",
        "order_status": "complete",
        "trigger_price": 0.0,
        "timestamp": "08-Apr-2025 13:58:03"
      },
      {
        "action": "BUY",
        "symbol": "YESBANK",
        "exchange": "NSE",
        "orderid": "250408001002736",
        "product": "MIS",
        "quantity": "1",
        "price": 16.5,
        "pricetype": "LIMIT",
        "order_status": "cancelled",
        "trigger_price": 0.0,
        "timestamp": "08-Apr-2025 14:13:45"
      }
    ],
    "statistics": {
      "total_buy_orders": 2.0,
      "total_sell_orders": 0.0,
      "total_completed_orders": 1.0,
      "total_open_orders": 0.0,
      "total_rejected_orders": 0.0
    }
  }
}
```

### TradeBook Example

```go
response, err := client.TradeBook()
if err != nil {
    log.Printf("Error: %v", err)
}
fmt.Printf("%v\n", response)
```

**TradeBook Response**

```json
{
  "status": "success",
  "data": [
    {
      "action": "BUY",
      "symbol": "RELIANCE",
      "exchange": "NSE",
      "orderid": "250408000989443",
      "product": "MIS",
      "quantity": 0.0,
      "average_price": 1180.1,
      "timestamp": "13:58:03",
      "trade_value": 1180.1
    },
    {
      "action": "SELL",
      "symbol": "NHPC",
      "exchange": "NSE",
      "orderid": "250408001086129",
      "product": "MIS",
      "quantity": 0.0,
      "average_price": 83.74,
      "timestamp": "14:28:49",
      "trade_value": 83.74
    }
  ]
}
```

### PositionBook Example

```go
response, err := client.PositionBook()
if err != nil {
    log.Printf("Error: %v", err)
}
fmt.Printf("%v\n", response)
```

**PositionBook Response**

```json
{
  "status": "success",
  "data": [
    {
      "symbol": "NHPC",
      "exchange": "NSE",
      "product": "MIS",
      "quantity": "-1",
      "average_price": "83.74",
      "ltp": "83.72",
      "pnl": "0.02"
    },
    {
      "symbol": "RELIANCE",
      "exchange": "NSE",
      "product": "MIS",
      "quantity": "0",
      "average_price": "0.0",
      "ltp": "1189.9",
      "pnl": "5.90"
    },
    {
      "symbol": "YESBANK",
      "exchange": "NSE",
      "product": "MIS",
      "quantity": "-104",
      "average_price": "17.2",
      "ltp": "17.31",
      "pnl": "-10.44"
    }
  ]
}
```

### Holdings Example

```go
response, err := client.Holdings()
if err != nil {
    log.Printf("Error: %v", err)
}
fmt.Printf("%v\n", response)
```

**Holdings Response**

```json
{
  "status": "success",
  "data": {
    "holdings": [
      {
        "symbol": "RELIANCE",
        "exchange": "NSE",
        "product": "CNC",
        "quantity": 1,
        "pnl": -149.0,
        "pnlpercent": -11.1
      },
      {
        "symbol": "TATASTEEL",
        "exchange": "NSE",
        "product": "CNC",
        "quantity": 1,
        "pnl": -15.0,
        "pnlpercent": -10.41
      },
      {
        "symbol": "CANBK",
        "exchange": "NSE",
        "product": "CNC",
        "quantity": 5,
        "pnl": -69.0,
        "pnlpercent": -13.43
      }
    ],
    "statistics": {
      "totalholdingvalue": 1768.0,
      "totalinvvalue": 2001.0,
      "totalprofitandloss": -233.15,
      "totalpnlpercentage": -11.65
    }
  }
}
```

### Analyzer Status Example

```go
response, err := client.AnalyzerStatus()
if err != nil {
    log.Printf("Error: %v", err)
}
fmt.Printf("%v\n", response)
```

**Analyzer Status Response**

```json
{
  "data": {
    "analyze_mode": true,
    "mode": "analyze",
    "total_logs": 2
  },
  "status": "success"
}
```

### Analyzer Toggle Example

```go
// Switch to analyze mode (simulated responses)
response, err := client.AnalyzerToggle(true)
if err != nil {
    log.Printf("Error: %v", err)
}
fmt.Printf("%v\n", response)
```

**Analyzer Toggle Response**

```json
{
  "data": {
    "analyze_mode": true,
    "message": "Analyzer mode switched to analyze",
    "mode": "analyze",
    "total_logs": 2
  },
  "status": "success"
}
```

### LTP Data (Streaming Websocket)

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/marketcalls/openalgo-go/openalgo"
)

func main() {
    // Initialize OpenAlgo client
    client := openalgo.NewClient(
        "your_api_key",           // Replace with your actual OpenAlgo API key
        "http://127.0.0.1:5000",  // REST API host
        "v1",                     // API version
        "ws://127.0.0.1:8765",    // WebSocket host (optional)
    )

    // Connect to WebSocket
    if err := client.Connect(); err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer client.Disconnect()

    // Define instruments to subscribe for LTP
    instruments := []openalgo.Instrument{
        {Exchange: "NSE", Symbol: "RELIANCE"},
        {Exchange: "NSE", Symbol: "INFY"},
    }

    // Callback function for LTP updates
    onLTP := func(data interface{}) {
        fmt.Println("LTP Update Received:")
        fmt.Println(data)
    }

    // Subscribe
    if err := client.SubscribeLTP(instruments, onLTP); err != nil {
        log.Printf("Error subscribing to LTP: %v", err)
    }

    // Run for a few seconds to receive data
    time.Sleep(10 * time.Second)

    // Unsubscribe
    client.UnsubscribeLTP(instruments)
}
```

### Quotes (Streaming Websocket)

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/marketcalls/openalgo-go/openalgo"
)

func main() {
    // Initialize OpenAlgo client
    client := openalgo.NewClient(
        "your_api_key",           // Replace with your actual OpenAlgo API key
        "http://127.0.0.1:5000",  // REST API host
        "v1",                     // API version
        "ws://127.0.0.1:8765",    // WebSocket host (optional)
    )

    // Connect
    if err := client.Connect(); err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer client.Disconnect()

    // Instruments list
    instruments := []openalgo.Instrument{
        {Exchange: "NSE", Symbol: "RELIANCE"},
        {Exchange: "NSE", Symbol: "INFY"},
    }

    // Callback for Quote updates
    onQuote := func(data interface{}) {
        fmt.Println("Quote Update Received:")
        fmt.Println(data)
    }

    // Subscribe to quote stream
    if err := client.SubscribeQuote(instruments, onQuote); err != nil {
        log.Printf("Error subscribing to quotes: %v", err)
    }

    // Keep the script running to receive data
    time.Sleep(10 * time.Second)

    // Unsubscribe
    client.UnsubscribeQuote(instruments)
}
```

### Depth (Streaming Websocket)

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/marketcalls/openalgo-go/openalgo"
)

func main() {
    // Initialize OpenAlgo client
    client := openalgo.NewClient(
        "your_api_key",           // Replace with your actual OpenAlgo API key
        "http://127.0.0.1:5000",  // REST API host
        "v1",                     // API version
        "ws://127.0.0.1:8765",    // WebSocket host (optional)
    )

    // Connect
    if err := client.Connect(); err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer client.Disconnect()

    // Instruments list for depth
    instruments := []openalgo.Instrument{
        {Exchange: "NSE", Symbol: "RELIANCE"},
        {Exchange: "NSE", Symbol: "INFY"},
    }

    // Callback for market depth updates
    onDepth := func(data interface{}) {
        fmt.Println("Market Depth Update Received:")
        fmt.Println(data)
    }

    // Subscribe to depth stream
    if err := client.SubscribeDepth(instruments, onDepth); err != nil {
        log.Printf("Error subscribing to depth: %v", err)
    }

    // Run for a few seconds to collect data
    time.Sleep(10 * time.Second)

    // Unsubscribe
    client.UnsubscribeDepth(instruments)
}
```