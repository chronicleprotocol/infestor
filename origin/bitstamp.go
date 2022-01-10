package origin

import (
	"fmt"
	"strings"

	"github.com/chronicleprotocol/infestor/smocker"
)

type BitStamp struct{}

func (b BitStamp) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	return CombineMocks(e, b.build)
}

func (b BitStamp) build(e ExchangeMock) (*smocker.Mock, error) {
	symbol := strings.ToLower(e.Symbol.Format("%s%s"))
	body := `{
	"last": "%f",
	"timestamp": "%d",
	"bid": "%f",
	"volume": "%f",
	"ask": "%f"
}`

	return &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.NewStringMatcher("GET"),
			Path:   smocker.NewStringMatcher(fmt.Sprintf("/api/v2/ticker/%s", symbol)),
		},
		Response: &smocker.MockResponse{
			Status: e.StatusCode,
			Headers: map[string]smocker.StringSlice{
				"Content-Type": []string{
					"application/json",
				},
			},
			Body: fmt.Sprintf(body, e.Price, e.Timestamp.Unix(), e.Bid, e.Volume, e.Ask),
		},
	}, nil
}
