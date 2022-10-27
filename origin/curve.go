package origin

// Simulate an ETHRPC node returning the price for STETH/ETH on Curve: get_dy(int128,int128,uint256)
// https://etherscan.io/address/0xdc24316b9ae028f1497c275eb9192a3ea0f67022#code

import (
	"fmt"

	"github.com/chronicleprotocol/infestor/smocker"
)

type Curve struct{}

func (b Curve) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	mocks := make([]*smocker.Mock, 0)

	n := smocker.ShouldContainSubstring("net_version")
	mocks = append(mocks, &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.ShouldEqual("POST"),
			Path:   smocker.ShouldEqual("/"),
			Body: &smocker.BodyMatcher{
				BodyString: &n,
			},
		},
		Response: &smocker.MockResponse{
			Status: 200,
			Headers: map[string]smocker.StringSlice{
				"Content-Type": []string{
					"application/json",
				},
			},
			Body: "{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"1\"}",
		},
	})

	m, err := CombineMocks(e, b.build)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	return mocks, nil
}

func (b Curve) build(e ExchangeMock) (*smocker.Mock, error) {
	m := smocker.ShouldContainSubstring(e.Custom["match"])
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
			// Multicall response:
			Body: fmt.Sprintf(rpcJSONResult, e.Custom["response"]),
		},
	}, nil
}
