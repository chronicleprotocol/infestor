package origin

import (
	"fmt"
)

type Huobi struct{}

func (b Huobi) BuildMock(e Exchange) ([]byte, error) {
	yaml := `
- request:
    method: GET
    path: '/market/detail/merged'
    query_params:
      symbol: '%s'
  dynamic_response:
    engine: go_template
    script: >
      headers:
        Content-Type: [application/json]
      body: >
        {
          "ch": "market.%s.detail.merged",
          "status": "ok",
          "ts": %d,
          "tick": {
              "id": 239151223299,
              "version": 239151223299,
              "open": 0.086276,
              "close": 0.08741,
              "low": 0.086266,
              "high": 0.088377,
              "amount": 10284.8934,
              "vol": %f,
              "count": 41879,
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
	symbol := e.Symbol.Format("%s%s")
	return []byte(fmt.Sprintf(yaml, symbol, symbol, e.Timestamp.UnixMilli(), e.Volume, e.Bid, e.Ask)), nil
}
