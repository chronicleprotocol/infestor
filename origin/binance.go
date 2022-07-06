package origin

import (
	"fmt"

	"github.com/chronicleprotocol/infestor/smocker"
)

type Binance struct{}

func (b Binance) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	mocksOne, err := CombineMocks(e, b.build)
	if err != nil {
		return nil, err
	}
	tickers, err := CombineMocks(e, b.buildWholeDay)
	if err != nil {
		return nil, err
	}
	return append(mocksOne, tickers...), nil
}

func (b Binance) buildWholeDay(e ExchangeMock) (*smocker.Mock, error) {
	symbol := e.Symbol.Format("%s%s")
	body := `[{"symbol":"%s","lastPrice":"%.8f","bidPrice":"%.8f","askPrice":"%.8f","volume":"%.8f","closeTime":%d}]`

	return &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.ShouldEqual("GET"),
			Path:   smocker.ShouldEqual("/api/v3/ticker/24hr"),
		},
		Response: &smocker.MockResponse{
			Status: e.StatusCode,
			Headers: map[string]smocker.StringSlice{
				"Content-Type": []string{
					"application/json",
				},
			},
			Body: fmt.Sprintf(body, symbol, e.Price, e.Bid, e.Ask, e.Volume, e.Timestamp.UnixMilli()),
		},
	}, nil
}

func (b Binance) build(e ExchangeMock) (*smocker.Mock, error) {
	symbol := e.Symbol.Format("%s%s")
	body := `{"symbol": "%s","price": "%.8f"}`

	return &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.ShouldEqual("GET"),
			Path:   smocker.ShouldEqual("/api/v3/ticker/price"),
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
			Body: fmt.Sprintf(body, symbol, e.Price),
		},
	}, nil
}
