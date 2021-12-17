# Infestor
Library for writing integration tests based for different exchanges.
It offers a simple way to mock exchange API responses and test your code.

## Example:

```go
package yourpack

import (
	"github.com/chronicleprotocol/infestor"
	"github.com/chronicleprotocol/infestor/smocker"
)

api := smocker.NewApi("http://localhost", 8081)

err := infestor.NewMocksBuilder()
    .Reset()
    .Add(origin.NewExchange("binance").WithSymbol("ETHBTC").WithPrice(1))
    .Add(infestor.NewExchange("kraken").WithSymbol("ETHBTC").WithPrice(2))
    .Deploy(api)

if err != nil {
	// ... Didn't able to send required mocks for smoker
}

// Run your tests with base URL `http://localhost:8080` for your exchanges
// Example for binance: `http://localhost:8080/api/v3/ticker/price?symbol=ETHBTC`

// it will reply with:
// {
//  "symbol": "ETHBTC",
//  "price": "1"
// }
```