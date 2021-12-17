package origin

import "fmt"

type GateIO struct{}

func (g GateIO) BuildMock(e Exchange) ([]byte, error) {
	yaml := `
- request:
    method: GET
    path: '/api/v4/spot/tickers'
		query_params:
			currency_pair: %s
  response:
    status: %d
    headers:
      Content-Type: application/json
    body: |-
      [
    		{
					"currency_pair": "%s",
					"last": "%f",
					"lowest_ask": "%f",
					"highest_bid": "%f",
					"change_percentage": "-0.62",
					"base_volume": "%f",
					"quote_volume": "36.88146210359657",
					"high_24h": "0.084826",
					"low_24h": "0.082033"
				}
			]`

	symbol := e.Symbol.Format("%s_%s")
	return []byte(fmt.Sprintf(
		yaml,
		symbol,
		e.StatusCode,
		symbol,
		e.Price,
		e.Ask,
		e.Bid,
		e.Volume,
	)), nil
}
