package origin

import (
	"fmt"
	"strings"
)

type BitStamp struct{}

func (b BitStamp) BuildMock(e Exchange) ([]byte, error) {
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
					"high": "0.08234246",
					"last": "%f",
					"timestamp": "%d",
					"bid": "%f",
					"vwap": "0.08262116",
					"volume": "%f",
					"low": "0.07913291",
					"ask": "%f",
					"open": "0.08234246"
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
