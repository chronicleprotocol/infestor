package origin

import (
	"fmt"
	"strings"
)

type BitStamp struct{}

func (b BitStamp) BuildMock(e []ExchangeMock) ([]byte, error) {
	return CombineMocks(e, b.build)
}

func (b BitStamp) build(e ExchangeMock) ([]byte, error) {
	yaml := `
- request:
    method: GET
    path: '/api/v2/ticker/%s'
  response:
    status: %d
    headers:
      Content-Type: [application/json]
    body: |-
      {
        "last": "%f",
        "timestamp": "%d",
        "bid": "%f",
        "volume": "%f",
        "ask": "%f"
      }`

	return []byte(fmt.Sprintf(
		yaml,
		strings.ToLower(e.Symbol.Format("%s%s")),
		e.StatusCode,
		e.Price,
		e.Timestamp.Unix(),
		e.Bid,
		e.Volume,
		e.Ask,
	)), nil
}
