package origin

// Simulate an ETHRPC node returning the price for WstETH
// https://etherscan.io/address/0x7f39C581F595B53c5cb19bD0b3f8dA6c935E2Ca0#code

import (
	"fmt"

	"github.com/chronicleprotocol/infestor/smocker"
)

type WSTETH struct{}

func (b WSTETH) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	mocks := make([]*smocker.Mock, 0)
	m, err := CombineMocks(e, b.buildTokensPerStEth)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	m, err = CombineMocks(e, b.buildstEthPerToken)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	return mocks, nil
}

func (b WSTETH) buildTokensPerStEth(e ExchangeMock) (*smocker.Mock, error) {
	// tokensPerStEth
	m := smocker.ShouldContainSubstring("0x9576a0c8")

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
			Body: fmt.Sprintf(RpcJSONResult, "0x0000000000000000000000000000000000000000000000000cf6c97c561e07bc"),
		},
	}, nil
}

func (b WSTETH) buildstEthPerToken(e ExchangeMock) (*smocker.Mock, error) {
	// stEthPerToken
	m := smocker.ShouldContainSubstring("0x035faf82")

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
			Body: fmt.Sprintf(RpcJSONResult, "0x0000000000000000000000000000000000000000000000000edb20f642b06506"),
		},
	}, nil
}
