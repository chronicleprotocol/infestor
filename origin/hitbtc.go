package origin

import (
	"fmt"

	"github.com/chronicleprotocol/infestor/smocker"
)

type HitBTC struct{}

func (h HitBTC) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	return CombineMocks(e, h.build)
}

func (h HitBTC) build(e ExchangeMock) (*smocker.Mock, error) {
	symbol := e.Symbol.Format("%s%s")
	body := `{
	"symbol": "%s",
	"ask": "%f",
	"bid": "%f",
	"last": "%f",
	"volume": "%f",
	"timestamp": "%s"
}`

	return &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.ShouldEqual("GET"),
			Path:   smocker.ShouldEqual(fmt.Sprintf("/api/2/public/ticker/%s", symbol)),
		},
		Response: &smocker.MockResponse{
			Status: e.StatusCode,
			Headers: map[string]smocker.StringSlice{
				"Content-Type": []string{
					"application/json",
				},
			},
			Body: fmt.Sprintf(body, symbol,
				e.Ask,
				e.Bid,
				e.Price,
				e.Volume,
				e.Timestamp.String()),
		},
	}, nil
}
