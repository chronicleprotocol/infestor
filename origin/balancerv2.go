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
	m, err := CombineMocks(e, b.buildStethWeth)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	m, err = CombineMocks(e, b.buildRethWeth)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	m, err = CombineMocks(e, b.buildRethWethRef)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	return mocks, nil
}

// "Ref:RETH/WETH": "0xae78736Cd615f374D3085123A210448E74Fc6393",
func (b BalancerV2) buildRethWethRef(e ExchangeMock) (*smocker.Mock, error) {
	m := smocker.ShouldContainSubstring("0xae78736Cd615f374D3085123A210448E74Fc6393")

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

// "RETH/WETH": "0x1E19CF2D73a72Ef1332C882F20534B6519Be0276",
func (b BalancerV2) buildRethWeth(e ExchangeMock) (*smocker.Mock, error) {
	m := smocker.ShouldContainSubstring("0x1E19CF2D73a72Ef1332C882F20534B6519Be0276")

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

// "STETH/WETH": "0x32296969ef14eb0c6d29669c550d4a0449130230",
func (b BalancerV2) buildStethWeth(e ExchangeMock) (*smocker.Mock, error) {
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
