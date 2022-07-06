package origin

import (
	"fmt"

	"github.com/chronicleprotocol/infestor/smocker"
)

type CryptoCompare struct{}

func (c CryptoCompare) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	return CombineMocks(e, c.build)
}

func (c CryptoCompare) build(e ExchangeMock) (*smocker.Mock, error) {
	body := `{"%s": %f}`

	return &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.ShouldEqual("GET"),
			Path:   smocker.ShouldEqual("/data/price"),
			QueryParams: map[string]smocker.StringMatcherSlice{
				"fsym": []smocker.StringMatcher{
					smocker.ShouldEqual(e.Symbol.Base),
				},
				"tsyms": []smocker.StringMatcher{
					smocker.ShouldEqual(e.Symbol.Quote),
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
			Body: fmt.Sprintf(body, e.Symbol.Quote, e.Price),
		},
	}, nil
}
