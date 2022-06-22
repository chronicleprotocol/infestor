package origin

import (
	"fmt"

	"github.com/chronicleprotocol/infestor/smocker"
)

type Coinbase struct{}

func (c Coinbase) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	return CombineMocks(e, c.build)
}

func (c Coinbase) build(e ExchangeMock) (*smocker.Mock, error) {
	format := "2006-01-02T15:04:05.999999Z"

	body := `{
	"price": "%f",
	"time": "%s",
	"bid": "%f",
	"ask": "%f",
	"volume": "%f"
}`

	return &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.ShouldEqual("GET"),
			Path:   smocker.ShouldEqual(fmt.Sprintf("/products/%s/ticker", e.Symbol.Format("%s-%s"))),
		},
		Response: &smocker.MockResponse{
			Status: e.StatusCode,
			Headers: map[string]smocker.StringSlice{
				"Content-Type": []string{
					"application/json",
				},
			},
			Body: fmt.Sprintf(body, e.Price, e.Timestamp.Format(format), e.Bid, e.Ask, e.Volume),
		},
	}, nil
}
