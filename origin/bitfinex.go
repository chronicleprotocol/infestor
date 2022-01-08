package origin

import (
	"fmt"
)

// NOTE: For symbols of 4 chars you have to write `SYMBOL:` otherwise API request to smocker will fail.
// Example: AVAX/USD should be written in mock as `AVAX:/USD`

type Bitfinex struct{}

func (b Bitfinex) BuildMock(e []ExchangeMock) ([]byte, error) {
	return CombineMocks(e, b.build)
}

func (b Bitfinex) build(e ExchangeMock) ([]byte, error) {
	yaml := `
- request:
    method: GET
    path: '/v2/ticker/t%s'
  response:
    status: %d
    headers:
      Content-Type: [application/json]
    body: |-
      [
        %f,
        90.17754546000003,
        %f,
        77.37476201,
        0.001204,
        0.0139,
        %f,
        %f,
        0.088377,
        0.08629
      ]`

	return []byte(fmt.Sprintf(yaml, e.Symbol.Format("%s%s"), e.StatusCode, e.Bid, e.Ask, e.Price, e.Volume)), nil
}
