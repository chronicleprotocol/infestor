package origin

import "fmt"

type Coinbase struct{}

func (c Coinbase) BuildMock(e []ExchangeMock) ([]byte, error) {
	return CombineMocks(e, c.build)
}

func (c Coinbase) build(e ExchangeMock) ([]byte, error) {
	format := "2006-01-02T15:04:05.999999Z"

	yaml := `
- request:
    method: GET
    path: '/products/%s/ticker'
  response:
    status: %d
    headers:
      Content-Type: [application/json]
    body: |-
      {
        "price": "%f",
        "time": "%s",
        "bid": "%f",
        "ask": "%f",
        "volume": "%f"
      }`
	return []byte(
		fmt.Sprintf(
			yaml,
			e.Symbol.Format("%s-%s"),
			e.StatusCode,
			e.Price,
			e.Timestamp.Format(format),
			e.Bid,
			e.Ask,
			e.Volume,
		),
	), nil
}
