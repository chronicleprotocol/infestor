package origin

import (
	"fmt"

	"github.com/chronicleprotocol/infestor/smocker"
)

type Bithumb struct{}

func (b Bithumb) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	return CombineMocks(e, b.build)
}

func (b Bithumb) build(e ExchangeMock) (*smocker.Mock, error) {
	symbol := e.Symbol.Format("%s-%s")
	body := `
{
	"data": [
		{
			"vol": "%f",
			"c": "%f",
			"s": "%s"
		}
	],
	"code": "0",
	"msg": "success",
	"timestamp": %d,
	"startTime": null
}
`

	return &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.NewStringMatcher("GET"),
			Path:   smocker.NewStringMatcher("/openapi/v1/spot/ticker"),
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
			Body: fmt.Sprintf(body, e.Volume, e.Price, symbol, e.Timestamp.UnixMilli()),
		},
	}, nil
}
