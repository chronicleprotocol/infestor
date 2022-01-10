package origin

import (
	"fmt"

	"github.com/chronicleprotocol/infestor/smocker"
)

type BitStamp struct{}

func (b BitStamp) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	return CombineMocks(e, b.build)
}

func (b BitStamp) build(e ExchangeMock) (*smocker.Mock, error) {
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
			Path:   smocker.NewStringMatcher(fmt.Sprintf("/api/v2/ticker/%s", e.Symbol.Format("%s%s"))),
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
