package origin

import (
	"fmt"
)

type Okex struct{}

func (o Okex) BuildMock(e Exchange) ([]byte, error) {
	yaml := `
- request:
    method: GET
    path: '/api/spot/v3/instruments/%s/ticker'
  response:
    status: %d
    headers:
      Content-Type: application/json
    body: |-
      {
				"best_ask": "%f",
				"best_bid": "%f",
				"instrument_id": "%s",
				"open_utc0": "0.08309",
				"open_utc8": "0.08406",
				"product_id": "%s",
				"last": "%f",
				"last_qty": "0.350759",
				"ask": "%f",
				"best_ask_size": "2.981467",
				"bid": "%f",
				"best_bid_size": "5.389945",
				"open_24h": "0.08262",
				"high_24h": "0.08481",
				"low_24h": "0.08209",
				"base_volume_24h": "%f",
				"timestamp": "%s",
				"quote_volume_24h": "575.58243"
			}`

	symbol := e.Symbol.Format("%s-%s")
	return []byte(fmt.Sprintf(
		yaml,
		symbol,
		e.StatusCode,
		e.Ask,
		e.Bid,
		symbol,
		symbol,
		e.Price,
		e.Ask,
		e.Bid,
		e.Volume,
		e.Timestamp.String(),
	)), nil
}
