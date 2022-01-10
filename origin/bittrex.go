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
	symbol := fmt.Sprintf("%s-%s", e.Symbol.Quote, e.Symbol.Base)
	body := `{
	 "success": true,
	 "message": "",
	 "result": {
		 "Bid": %f,
		 "Ask": %f,
		 "Last": %f
	 }
 }`

	return &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.NewStringMatcher("GET"),
			Path:   smocker.NewStringMatcher("/api/v1.1/public/getticker"),
			QueryParams: map[string]smocker.StringMatcherSlice{
				"market": []smocker.StringMatcher{
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
			Body: fmt.Sprintf(body, e.Bid, e.Ask, e.Price),
		},
	}, nil
}
