package origin

import (
	"fmt"
)

type Upbit struct{}

func (u Upbit) BuildMock(e []ExchangeMock) ([]byte, error) {
	return CombineMocks(e, u.build)
}

func (u Upbit) build(e ExchangeMock) ([]byte, error) {
	yaml := `
- request:
    method: GET
    path: '/v1/ticker'
    query_params:
      markets: %s
  response:
    status: %d
    headers:
      Content-Type: application/json
    body: |-
      [
        {
          "market": "%s",
          "trade_date": "%s",
          "trade_time": "%s",
          "trade_timestamp": %d,
          "trade_price": %f,
          "trade_volume": %f,
          "timestamp": %d
        }
      ]`

	symbol := fmt.Sprintf("%s-%s", e.Symbol.Quote, e.Symbol.Base)
	return []byte(fmt.Sprintf(
		yaml,
		symbol,
		e.StatusCode,
		symbol,
		e.Timestamp.Format("20060102"),
		e.Timestamp.Format("150405"),
		e.Timestamp.UnixMilli(),
		e.Price,
		e.Volume,
		e.Timestamp.UnixMilli(),
	)), nil
}
