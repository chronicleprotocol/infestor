package origin

import (
	"fmt"

	"github.com/chronicleprotocol/infestor/smocker"
)

type Upbit struct{}

func (u Upbit) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	return CombineMocks(e, u.build)
}

func (u Upbit) build(e ExchangeMock) (*smocker.Mock, error) {
	symbol := fmt.Sprintf("%s-%s", e.Symbol.Quote, e.Symbol.Base)
	body := `[
	{
		"market": "%s",
		"trade_date": "%s",
		"trade_time": "%s",
		"trade_timestamp": %d,
		"trade_price": %f,
		"trade_volume": %f,
		"timestamp": %d
	}
]`

	return &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.ShouldEqual("GET"),
			Path:   smocker.ShouldEqual("/v1/ticker"),
			QueryParams: map[string]smocker.StringMatcherSlice{
				"markets": []smocker.StringMatcher{
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
			Body: fmt.Sprintf(body,
				symbol,
				e.Timestamp.Format("20060102"),
				e.Timestamp.Format("150405"),
				e.Timestamp.UnixMilli(),
				e.Price,
				e.Volume,
				e.Timestamp.UnixMilli()),
		},
	}, nil
}
