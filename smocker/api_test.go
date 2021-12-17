package smocker

import (
	"context"
	"testing"
)

func TestResetAPICall(t *testing.T) {
	api := NewAPI("http://localhost", 8081)
	err := api.Reset(context.Background())
	if err != nil {
		t.Errorf("Error resetting API call: %s", err)
	}
}

func TestAddMock(t *testing.T) {
	mock := Mock{
		Body: []byte(`- request:
    method: POST
    path: '/subgraphs/name/balancer-labs/balancer'
    body:
      variables.id: 0xba100000625a3754423978a60c9317c58a424e3d
  response:
    status: 200
    headers:
      Content-Type: application/json
    body: |-
      {
          "data": {
              "tokenPrices": [
                  {
                      "poolLiquidity": "11108555",
                      "price": "15",
                      "symbol": "BAL"
                  }
              ]
          }
      }`),
	}
	api := NewAPI("http://localhost", 8081)
	err := api.AddMock(context.Background(), mock)
	if err != nil {
		t.Errorf("Error adding mock: %s", err)
	}
}
