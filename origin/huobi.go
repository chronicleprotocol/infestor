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
	mocksOne, err := CombineMocks(e, h.buildForOne)
	if err != nil {
		return nil, err
	}
	tickers, err := CombineMocks(e, h.buildTickers)
	if err != nil {
		return nil, err
	}
	return append(mocksOne, tickers...), nil
}

func (h Huobi) buildTickers(e ExchangeMock) (*smocker.Mock, error) {
	symbol := strings.ToLower(e.Symbol.Format("%s%s"))
	price := e.Price
	if e.Bid != 0 {
		price = e.Bid
	}
	body := `{
		"status": "ok",
		"ts": %d,
		"data": [
			{
				"symbol":"%s",
				"open":%f,     
				"high":%f,     
				"low":%f,        
				"close":%f,     
				"amount":36551302.17544405,
				"vol":%f,
				"count":1709,
				"bid":%f,
				"bidSize":54300.341,
				"ask":%f,
				"askSize":1923.4879
				}
		]
	}`

	return &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.ShouldEqual("GET"),
			Path:   smocker.ShouldEqual("/market/tickers"),
		},
		Response: &smocker.MockResponse{
			Status: e.StatusCode,
			Headers: map[string]smocker.StringSlice{
				"Content-Type": []string{
					"application/json",
				},
			},
			Body: fmt.Sprintf(body, e.Timestamp.UnixMilli(), symbol, price, price, price, price, e.Volume, price, e.Ask),
		},
	}, nil
}

func (h Huobi) buildForOne(e ExchangeMock) (*smocker.Mock, error) {
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
			Method: smocker.ShouldEqual("GET"),
			Path:   smocker.ShouldEqual("/market/detail/merged"),
			QueryParams: map[string]smocker.StringMatcherSlice{
				"symbol": []smocker.StringMatcher{
					smocker.ShouldEqual(symbol),
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
