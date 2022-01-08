package origin

import (
	"fmt"
)

type KuCoin struct{}

func (k KuCoin) BuildMock(e []ExchangeMock) ([]byte, error) {
	return CombineMocks(e, k.build)
}

func (k KuCoin) build(e ExchangeMock) ([]byte, error) {
	yaml := `
- request:
    method: GET
    path: '/api/v1/market/orderbook/level1'
    query_params:
      symbol: %s
  response:
    status: %d
    headers:
      Content-Type: application/json
    body: |-
      {
        "code": "200000",
        "data": {
          "time": %d,
          "sequence": "1615098154456",
          "price": "%f",
          "size": "0.0036768",
          "bestBid": "%f",
          "bestBidSize": "7.4758085",
          "bestAsk": "%f",
          "bestAskSize": "5.5416409"
        }
      }`

	symbol := e.Symbol.Format("%s-%s")
	return []byte(fmt.Sprintf(
		yaml,
		symbol,
		e.StatusCode,
		e.Timestamp.UnixMilli(),
		e.Price,
		e.Bid,
		e.Ask,
	)), nil
}
