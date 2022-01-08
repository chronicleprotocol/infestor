package origin

import (
	"fmt"
	"strings"
)

// To setup price for huobi we take it's bid value.
// So if you don't set `WithBid()` we will use `Price` field, but if you will set `Bid`
// - `Price` value will be ignored

type Huobi struct{}

func (h Huobi) BuildMock(e []ExchangeMock) ([]byte, error) {
	return CombineMocks(e, h.build)
}

func (h Huobi) build(e ExchangeMock) ([]byte, error) {
	yaml := `
- request:
    method: GET
    path: '/market/detail/merged'
    query_params:
      symbol: '%s'
  response:
    status: %d
    headers:
      Content-Type: [application/json]
    body: >
      {
        "ch": "market.%s.detail.merged",
        "status": "ok",
        "ts": %d,
        "tick": {
          "vol": %f,
          "bid": [
            %f,
            0.3618
          ],
          "ask": [
            %f,
            1.947
          ]
        }
      }`
	symbol := strings.ToLower(e.Symbol.Format("%s%s"))
	price := e.Price
	if e.Bid != 0 {
		price = e.Bid
	}
	return []byte(fmt.Sprintf(yaml, symbol, e.StatusCode, symbol, e.Timestamp.UnixMilli(), e.Volume, price, e.Ask)), nil
}
