package origin

import (
	"fmt"

	"github.com/chronicleprotocol/infestor/smocker"
)

type Balancer struct{}

func (b Balancer) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	return CombineMocks(e, b.build)
}

func (b Balancer) build(e ExchangeMock) (*smocker.Mock, error) {
	contract, ok := e.Custom["contract"]
	if !ok {
		return nil, fmt.Errorf("`contract` custom field is requierd for balancer")
	}

	bodyStr := `
{
	"data": {
		"tokenPrices": [
			{
				"price": "%.8f",
				"symbol": "%s"
			}
		]
	}
}`

	mock := &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.ShouldEqual("POST"),
			Path:   smocker.ShouldEqual("/subgraphs/name/balancer-labs/balancer"),
			Body: &smocker.BodyMatcher{
				BodyJSON: map[string]smocker.StringMatcher{
					"variables.id": smocker.ShouldEqual(contract),
				},
			},
		},
		Response: &smocker.MockResponse{
			Status: e.StatusCode,
			Headers: map[string]smocker.StringSlice{
				"Content-Type": []string{
					"application/json",
				},
			},
			Body: fmt.Sprintf(bodyStr, e.Price, e.Symbol.Base),
		},
	}
	return mock, nil
}
