package origin

import "fmt"

type Coinbase struct{}

func (b Coinbase) BuildMock(e Exchange) ([]byte, error) {
	format := "2006-01-02T15:04:05.999999Z"

	yaml := `
- request:
    method: GET
    path: '/products/%s/ticker'
  dynamic_response:
    engine: go_template
    script: >
      headers:
        Content-Type: [application/json]
      body: >
        {
          "trade_id": 24292500,
          "price": "%f",
          "size": "0.16783975",
          "time": "%s",
          "bid": "%f",
          "ask": "%f",
          "volume": "%f"
        }`
	return []byte(
		fmt.Sprintf(yaml, e.Symbol.Format("%s%s"), e.Price, e.Timestamp.Format(format), e.Bid, e.Ask, e.Volume),
	), nil
}
