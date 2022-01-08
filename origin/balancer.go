package origin

import (
	"fmt"
)

type Balancer struct{}

func (b Balancer) BuildMock(e []ExchangeMock) ([]byte, error) {
	return CombineMocks(e, b.build)
}

func (b Balancer) build(e ExchangeMock) ([]byte, error) {
	contract, ok := e.Custom["contract"]
	if !ok {
		return nil, fmt.Errorf("`contract` custom field is requierd for balancer")
	}

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
                      "price": "%.8f",
                      "symbol": "%s"
                  }
              ]
          }
      }`
	return []byte(fmt.Sprintf(yaml, contract, e.StatusCode, e.Price, e.Symbol.Base)), nil
}
