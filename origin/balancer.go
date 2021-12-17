package origin

import "fmt"

type Balancer struct{}

func (b Balancer) BuildMock(e Exchange) ([]byte, error) {
	yaml := `
- request:
    method: POST
    path: '/subgraphs/name/balancer-labs/balancer'
    body:
      variables.id: %s
  response:
    status: %d
    headers:
      Content-Type: application/json
    body: |-
      {
          "data": {
              "tokenPrices": [
                  {
                      "poolLiquidity": "11224",
                      "price": "%f",
                      "symbol": "%s"
                  }
              ]
          }
      }`
	return []byte(fmt.Sprintf(yaml, e.Symbol.Format("%s%s"), e.StatusCode, e.Price, e.Symbol.Base)), nil
}
