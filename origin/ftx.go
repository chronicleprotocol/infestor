package origin

import (
	"fmt"

	"github.com/chronicleprotocol/infestor/smocker"
)

type Ftx struct{}

func (f Ftx) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	return CombineMocks(e, f.build)
}

func (f Ftx) build(e ExchangeMock) (*smocker.Mock, error) {
	body := `{
	       "success": true,
	       "result": {
	         "name": "%s",
	         "last": %f,
	         "bid": %f,
	         "ask": %f,
	         "price": %f,
	         "type": "spot",
	         "baseCurrency": "%s",
	         "quoteCurrency": "%s",
	         "volumeUsd24h": %f
	       }
	     }`

	return &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.NewStringMatcher("GET"),
			Path:   smocker.NewStringMatcher(fmt.Sprintf("/api/markets/%s/%s", e.Symbol.Base, e.Symbol.Quote)),
		},
		Response: &smocker.MockResponse{
			Status: e.StatusCode,
			Headers: map[string]smocker.StringSlice{
				"Content-Type": []string{
					"application/json",
				},
			},
			Body: fmt.Sprintf(
				body,
				e.Symbol.String(),
				e.Price,
				e.Bid,
				e.Ask,
				e.Price,
				e.Symbol.Base,
				e.Symbol.Quote,
				e.Volume),
		},
	}, nil
}
