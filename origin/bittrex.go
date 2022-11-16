package origin

import (
	"fmt"

	"github.com/chronicleprotocol/infestor/smocker"
)

type BitTrex struct{}

func (b BitTrex) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	return CombineMocks(e, b.build)
}

func (b BitTrex) build(e ExchangeMock) (*smocker.Mock, error) {
	symbol := fmt.Sprintf("%s-%s", e.Symbol.Base, e.Symbol.Quote)
	fmt.Printf("james symbol: %s\n\n\n", symbol)
	body := `{
  "symbol": "%s",
  "lastTradeRate": "%f",
  "bidRate": "%f",
  "askRate": "%f"
}`

	return &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.ShouldEqual("GET"),
			Path:   smocker.ShouldEqual(fmt.Sprintf("/v3/markets/%s/ticker", symbol)),
		},
		Response: &smocker.MockResponse{
			Status: e.StatusCode,
			Headers: map[string]smocker.StringSlice{
				"Content-Type": []string{
					"application/json",
				},
			},
			Body: fmt.Sprintf(body, symbol, e.Price, e.Bid, e.Ask),
		},
	}, nil
}
