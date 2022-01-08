package origin

import (
	"fmt"
)

type CryptoCompare struct{}

func (c CryptoCompare) BuildMock(e []ExchangeMock) ([]byte, error) {
	return CombineMocks(e, c.build)
}

func (c CryptoCompare) build(e ExchangeMock) ([]byte, error) {
	// market = QUOTE-BASE
	yaml := `
- request:
    method: GET
    path: '/data/price'
    query_params:
      fsym: %s
      tsyms: %s
  response:
    status: %d
    headers:
      Content-Type: [application/json]
    body: |-
      {
        "%s": %f
      }`

	return []byte(fmt.Sprintf(
		yaml, e.Symbol.Base, e.Symbol.Quote, e.StatusCode, e.Symbol.Quote, e.Price,
	)), nil
}
