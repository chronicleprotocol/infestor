package origin

// Simulate an ETHRPC node returning the price for STETH/ETH on BalancerV2
// https://etherscan.io/address/0x32296969ef14eb0c6d29669c550d4a0449130230#code

import (
	"fmt"

	"github.com/chronicleprotocol/infestor/smocker"
)

type BalancerV2 struct{}

func (b BalancerV2) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	return CombineMocks(e, b.build)
}

func (b BalancerV2) build(e ExchangeMock) (*smocker.Mock, error) {
	m := smocker.ShouldContainSubstring("0x32296969ef14eb0c6d29669c550d4a0449130230")

	return &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.ShouldEqual("POST"),
			Path:   smocker.ShouldEqual("/"),
			Body: &smocker.BodyMatcher{
				BodyString: &m,
			},
		},
		Response: &smocker.MockResponse{
			Status: e.StatusCode,
			Headers: map[string]smocker.StringSlice{
				"Content-Type": []string{
					"application/json",
				},
			},
			Body: fmt.Sprintf(rpcJSONResult, "0x0000000000000000000000000000000000000000000000000dd22d6848e229b8"),
		},
	}, nil
}
