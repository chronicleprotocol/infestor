package origin

import (
	"fmt"
	"strings"

	"github.com/chronicleprotocol/infestor/smocker"
)

type Gemini struct{}

func (g Gemini) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	return CombineMocks(e, g.build)
}

func (g Gemini) build(e ExchangeMock) (*smocker.Mock, error) {
	symbol := strings.ToLower(e.Symbol.Format("%s%s"))
	body := `{
	"bid": "%f",
	"ask": "%f",
	"volume": {
		"%s": "%f",
		"%s": "%f",
		"timestamp": %d
	},
	"last": "%f"
}`

	return &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.NewStringMatcher("GET"),
			Path:   smocker.NewStringMatcher(fmt.Sprintf("/v1/pubticker/%s", symbol)),
		},
		Response: &smocker.MockResponse{
			Status: e.StatusCode,
			Headers: map[string]smocker.StringSlice{
				"Content-Type": []string{
					"application/json",
				},
			},
			Body: fmt.Sprintf(body, e.Bid,
				e.Ask,
				e.Symbol.Base,
				e.Volume,
				e.Symbol.Quote,
				e.Volume,
				e.Timestamp.UnixMilli(),
				e.Price),
		},
	}, nil
}
