package origin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/chronicleprotocol/infestor/smocker"
)

type Poloniex struct{}

func (p Poloniex) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	status := http.StatusOK
	if len(e) > 0 {
		status = e[0].StatusCode
	}
	return []*smocker.Mock{
		{
			Request: smocker.MockRequest{
				Method: smocker.ShouldEqual("GET"),
				Path:   smocker.ShouldEqual("/public"),
				QueryParams: map[string]smocker.StringMatcherSlice{
					"command": []smocker.StringMatcher{
						smocker.ShouldEqual("returnTicker"),
					},
				},
			},
			Response: &smocker.MockResponse{
				Status: status,
				Headers: map[string]smocker.StringSlice{
					"Content-Type": []string{
						"application/json",
					},
				},
				Body: fmt.Sprintf("{%s}", p.build(e)),
			},
		},
	}, nil
}

func (p Poloniex) build(mocks []ExchangeMock) string {
	var result []string
	yaml := `
        "%s": {
          "last": "%.8f",
          "lowestAsk": "%.8f",
          "highestBid": "%.8f",
          "baseVolume": "%.8f"
        }`
	for _, e := range p.filter(mocks) {
		result = append(
			result,
			fmt.Sprintf(
				yaml,
				fmt.Sprintf("%s_%s", e.Symbol.Quote, e.Symbol.Base),
				e.Price,
				e.Ask,
				e.Bid,
				e.Volume,
			),
		)
	}
	return strings.Join(result, ",")
}

// removes duplicated pairs from list, otherwise it will cause issues
func (p Poloniex) filter(mocks []ExchangeMock) map[string]ExchangeMock {
	filtered := map[string]ExchangeMock{}
	for _, e := range mocks {
		filtered[e.Symbol.String()] = e
	}
	return filtered
}
