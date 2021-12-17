package origin

import (
	"fmt"
)

type HitBTC struct{}

func (h HitBTC) BuildMock(e Exchange) ([]byte, error) {
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
				"low": "0.082113",
				"high": "0.084802",
				"open": "0.082663",
				"volume": "%f",
				"volumeQuote": "2054.2598378881",
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
