package origin

import (
	"fmt"
)

type Bitfinex struct{}

func (b Bitfinex) BuildMock(e Exchange) ([]byte, error) {
	yaml := `
- request:
    method: GET
    path: '/v2/ticker/t%s'
  response:
    status: %d
    headers:
      Content-Type: [application/json]
    body: |-
      [
        %f,
        90.17754546000003,
        %f,
        77.37476201,
        0.001204,
        0.0139,
        %f,
        %f,
        0.088377,
        0.08629
      ]`

	return []byte(fmt.Sprintf(yaml, e.Symbol.Format("%s%s"), e.StatusCode, e.Bid, e.Ask, e.Price, e.Volume)), nil
}
