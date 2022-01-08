package origin

import (
	"fmt"
)

type Okex struct{}

func (o Okex) BuildMock(e []ExchangeMock) ([]byte, error) {
	return CombineMocks(e, o.build)
}

func (o Okex) build(e ExchangeMock) ([]byte, error) {
	yaml := `
- request:
    method: GET
    path: '/api/spot/v3/instruments/%s/ticker'
  response:
    status: %d
    headers:
      Content-Type: application/json
    body: |-
      {
        "best_ask": "%f",
        "best_bid": "%f",
        "instrument_id": "%s",
        "product_id": "%s",
        "last": "%f",
        "ask": "%f",
        "bid": "%f",
        "base_volume_24h": "%f",
        "timestamp": "%s"
      }`

	symbol := e.Symbol.Format("%s-%s")
	return []byte(fmt.Sprintf(
		yaml,
		symbol,
		e.StatusCode,
		e.Ask,
		e.Bid,
		symbol,
		symbol,
		e.Price,
		e.Ask,
		e.Bid,
		e.Volume,
		e.Timestamp.String(),
	)), nil
}
