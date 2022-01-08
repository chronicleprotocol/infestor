package origin

import (
	"fmt"
	"strings"
)

type Gemini struct{}

func (g Gemini) BuildMock(e []ExchangeMock) ([]byte, error) {
	return CombineMocks(e, g.build)
}

func (g Gemini) build(e ExchangeMock) ([]byte, error) {
	yaml := `
- request:
    method: GET
    path: '/v1/pubticker/%s'
  response:
    status: %d
    headers:
      Content-Type: application/json
    body: |-
      {
        "bid": "%f",
        "ask": "%f",
        "volume": {
          "ETH": "%f",
          "BTC": "30.42238188324",
          "timestamp": %d
        },
        "last": "%f"
      }`

	symbol := strings.ToLower(e.Symbol.Format("%s%s"))
	return []byte(fmt.Sprintf(
		yaml,
		symbol,
		e.StatusCode,
		e.Bid,
		e.Ask,
		e.Volume,
		e.Timestamp.UnixMilli(),
		e.Price,
	)), nil
}
