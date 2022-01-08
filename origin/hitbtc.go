package origin

import (
	"fmt"
)

type HitBTC struct{}

func (h HitBTC) BuildMock(e []ExchangeMock) ([]byte, error) {
	return CombineMocks(e, h.build)
}

func (h HitBTC) build(e ExchangeMock) ([]byte, error) {
	yaml := `
- request:
    method: GET
    path: '/api/2/public/ticker/%s'
  response:
    status: %d
    headers:
      Content-Type: application/json
    body: |-
      {
        "symbol": "%s",
        "ask": "%f",
        "bid": "%f",
        "last": "%f",
        "volume": "%f",
        "timestamp": "%s"
      }`

	symbol := e.Symbol.Format("%s%s")
	return []byte(fmt.Sprintf(
		yaml,
		symbol,
		e.StatusCode,
		symbol,
		e.Ask,
		e.Bid,
		e.Price,
		e.Volume,
		e.Timestamp.String(),
	)), nil
}
