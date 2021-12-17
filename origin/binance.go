package origin

import (
	"fmt"
)

type Binance struct{}

func (b Binance) BuildMock(e Exchange) ([]byte, error) {
	yaml := `
- request:
    method: GET
    path: '/api/v3/ticker/price'
    query_params:
      symbol: "%s"
  response:
	status: %d
    headers:
	  Content-Type: [application/json]
    body: |-
	  {
	    "symbol": "%s",
	    "price": "%f"
	  }`
	symbol := e.Symbol.Format("%s%s")
	return []byte(fmt.Sprintf(yaml, symbol, e.StatusCode, symbol, e.Price)), nil
}
