package origin

import (
	"fmt"

	"github.com/chronicleprotocol/infestor/smocker"
)

type Binance struct{}

func (b Binance) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	return CombineMocks(e, b.build)
}

func (b Binance) build(e ExchangeMock) (*smocker.Mock, error) {
	symbol := e.Symbol.Format("%s%s")
	body := `{"symbol": "%s","price": "%.8f"}`

	return &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.NewStringMatcher("GET"),
			Path:   smocker.NewStringMatcher("/api/v3/ticker/price"),
			QueryParams: map[string]smocker.StringMatcherSlice{
				"symbol": []smocker.StringMatcher{
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
			Body: fmt.Sprintf(body, symbol, e.Price),
		},
	}, nil
}
