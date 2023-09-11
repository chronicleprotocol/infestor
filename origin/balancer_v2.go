package origin

// Simulate an ETHRPC node returning the price for STETH/ETH on BalancerV2
// https://etherscan.io/address/0x32296969ef14eb0c6d29669c550d4a0449130230#code

import (
	"fmt"

	"github.com/chronicleprotocol/infestor/smocker"
)

type BalancerV2 struct{}

func (b BalancerV2) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
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

	m, err := CombineMocks(e, b.buildGetLatest)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m)

	m, err := CombineMocks(e, b.build)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	return mocks, nil
}

func (b BalancerV2) build(e ExchangeMock) (*smocker.Mock, error) {
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
