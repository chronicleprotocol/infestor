package origin

import (
	"fmt"

	"github.com/chronicleprotocol/infestor/smocker"
)

// NOTE: For symbols of 4 chars you have to write `SYMBOL:` otherwise API request to smocker will fail.
// Example: AVAX/USD should be written in mock as `AVAX:/USD`

type Bitfinex struct{}

func (b Bitfinex) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	mocksOne, err := CombineMocks(e, b.build)
	if err != nil {
		return nil, err
	}
	other, err := CombineMocks(e, b.buildList)
	if err != nil {
		return nil, err
	}
	return append(mocksOne, other...), nil
}

func (b Bitfinex) buildList(e ExchangeMock) (*smocker.Mock, error) {
	symbol := e.Symbol.Format("t%s%s")

	body := `[[
	"%s",
	%f,
	90.17754546000003,
	%f,
	77.37476201,
	0.001204,
	0.0139,
	%f,
	%f,
	0.088377,
	0.08629
]]`

	return &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.ShouldEqual("GET"),
			Path:   smocker.ShouldEqual("/v2/tickers?symbols=ALL"),
			QueryParams: map[string]smocker.StringMatcherSlice{
				"symbols": []smocker.StringMatcher{
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
			Body: fmt.Sprintf(body, symbol, e.Bid, e.Ask, e.Price, e.Volume),
		},
	}, nil
}

func (b Bitfinex) build(e ExchangeMock) (*smocker.Mock, error) {
	body := `[
	%f,
	90.17754546000003,
	%f,
	77.37476201,
	0.001204,
	0.0139,
	%f,
	%f,
	0.088377,
	0.08629
]`

	return &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.ShouldEqual("GET"),
			Path:   smocker.ShouldEqual(fmt.Sprintf("/v2/ticker/t%s", e.Symbol.Format("%s%s"))),
		},
		Response: &smocker.MockResponse{
			Status: e.StatusCode,
			Headers: map[string]smocker.StringSlice{
				"Content-Type": []string{
					"application/json",
				},
			},
			Body: fmt.Sprintf(body, e.Bid, e.Ask, e.Price, e.Volume),
		},
	}, nil
}
