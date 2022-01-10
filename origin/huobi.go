package origin

import (
	"fmt"
	"strings"

	"github.com/chronicleprotocol/infestor/smocker"
)

// To setup price for huobi we take it's bid value.
// So if you don't set `WithBid()` we will use `Price` field, but if you will set `Bid`
// - `Price` value will be ignored

type Huobi struct{}

func (h Huobi) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	return CombineMocks(e, h.build)
}

func (h Huobi) build(e ExchangeMock) (*smocker.Mock, error) {
	symbol := strings.ToLower(e.Symbol.Format("%s%s"))
	price := e.Price
	if e.Bid != 0 {
		price = e.Bid
	}
	body := `{
	 "status": "ok",
	 "ts": %d,
	 "tick": {
		 "vol": %f,
		 "bid": [
			 %f,
			 0.3618
		 ],
		 "ask": [
			 %f,
			 1.947
		 ]
	 }
 }`

	return &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.NewStringMatcher("GET"),
			Path:   smocker.NewStringMatcher("/market/detail/merged"),
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
			Body: fmt.Sprintf(body, e.Timestamp.UnixMilli(), e.Volume, price, e.Ask),
		},
	}, nil
}
