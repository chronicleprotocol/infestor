package origin

import (
	"fmt"
)

type BitTrex struct{}

func (b BitTrex) BuildMock(e Exchange) ([]byte, error) {
	yaml := `
- request:
		method: GET
		path: '/api/v1.1/public/getticker'
		query_params:
			market: %s
	response:
		status: %d
		headers:
			Content-Type: [application/json]
		body: |-
			{
				"success": true,
				"message": "",
				"result": {
					"Bid": %f,
					"Ask": %f,
					"Last": %f
				}
			}`

	return []byte(fmt.Sprintf(
		yaml, fmt.Sprintf("%s-%s", e.Symbol.Quote, e.Symbol.Base), e.StatusCode, e.Bid, e.Ask, e.Price,
	)), nil
}
