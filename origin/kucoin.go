package origin

import (
	"fmt"

	"github.com/chronicleprotocol/infestor/smocker"
)

type KuCoin struct{}

func (k KuCoin) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	return CombineMocks(e, k.build)
}

func (k KuCoin) build(e ExchangeMock) (*smocker.Mock, error) {
	symbol := e.Symbol.Format("%s-%s")
	body := `{
	 "code": "200000",
	 "data": {
		 "time": %d,
		 "price": "%f",
		 "bestBid": "%f",
		 "bestAsk": "%f"
	 }
 }`

	return &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.NewStringMatcher("GET"),
			Path:   smocker.NewStringMatcher("/api/v1/market/orderbook/level1"),
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
			Body: fmt.Sprintf(body, e.Timestamp.UnixMilli(), e.Price, e.Bid, e.Ask),
		},
	}, nil
}
