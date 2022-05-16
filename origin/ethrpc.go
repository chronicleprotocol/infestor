package origin

// Simulate an ETHRPC node. Add any general ETH/RPC related mocks here.

import (
	"github.com/chronicleprotocol/infestor/smocker"
)

type EthRPC struct{}

func (b EthRPC) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	return CombineMocks(e, b.build)
}

func (b EthRPC) build(e ExchangeMock) (*smocker.Mock, error) {
	m := smocker.NewSubstringMatcher("eth_blockNumber")

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
			Body: `{"jsonrpc":"2.0","id":1,"result":"0xe125c8"}`,
		},
	}, nil
}
