package origin

import (
	"fmt"

	"github.com/chronicleprotocol/infestor/smocker"
)

type Okex struct{}

func (o Okex) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	return CombineMocks(e, o.build)
}

func (o Okex) build(e ExchangeMock) (*smocker.Mock, error) {
	symbol := e.Symbol.Format("%s-%s")
	body := `{
	"best_ask": "%f",
	"best_bid": "%f",
	"instrument_id": "%s",
	"product_id": "%s",
	"last": "%f",
	"ask": "%f",
	"bid": "%f",
	"base_volume_24h": "%f",
	"timestamp": "%s"
}`

	return &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.ShouldEqual("GET"),
			Path:   smocker.ShouldEqual(fmt.Sprintf("/api/spot/v3/instruments/%s/ticker", symbol)),
		},
		Response: &smocker.MockResponse{
			Status: e.StatusCode,
			Headers: map[string]smocker.StringSlice{
				"Content-Type": []string{
					"application/json",
				},
			},
			Body: fmt.Sprintf(body,
				e.Ask,
				e.Bid,
				symbol,
				symbol,
				e.Price,
				e.Ask,
				e.Bid,
				e.Volume,
				e.Timestamp.String()),
		},
	}, nil
}
