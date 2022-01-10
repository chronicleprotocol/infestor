package origin

import (
	"fmt"

	"github.com/chronicleprotocol/infestor/smocker"
)

type GateIO struct{}

func (g GateIO) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	return CombineMocks(e, g.build)
}

func (g GateIO) build(e ExchangeMock) (*smocker.Mock, error) {
	symbol := e.Symbol.Format("%s_%s")
	body := `[
	{
		"currency_pair": "%s",
		"last": "%f",
		"lowest_ask": "%f",
		"highest_bid": "%f",
		"base_volume": "%f"
	}
]`

	return &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.NewStringMatcher("GET"),
			Path:   smocker.NewStringMatcher("/api/v4/spot/tickers"),
			QueryParams: map[string]smocker.StringMatcherSlice{
				"currency_pair": []smocker.StringMatcher{
					smocker.NewStringMatcher(symbol),
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
			Body: fmt.Sprintf(body, symbol, e.Price, e.Ask, e.Bid, e.Volume),
		},
	}, nil
}
