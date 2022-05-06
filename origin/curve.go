package origin

// Simulate an ETHRPC node returning the price for STETH/ETH on Curve
// https://etherscan.io/address/0xdc24316b9ae028f1497c275eb9192a3ea0f67022#code

import (
	"fmt"

	"github.com/chronicleprotocol/infestor/smocker"
)

type Curve struct{}

func (b Curve) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	return CombineMocks(e, b.build)
}

func (b Curve) build(e ExchangeMock) (*smocker.Mock, error) {
	m := smocker.NewSubstringMatcher("0xdc24316b9ae028f1497c275eb9192a3ea0f67022")

	return &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.NewStringMatcher("POST"),
			Path:   smocker.NewStringMatcher("/"),
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
			Body: fmt.Sprintf(rpcJSONResult, "0x0000000000000000000000000000000000000000000000000dcfcc3d4023b410"),
		},
	}, nil
}
