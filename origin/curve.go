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
	m, err := CombineMocks(e, b.buildNoChecksum)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	m, err = CombineMocks(e, b.buildChecksum)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	return mocks, nil
}

func (b Curve) buildNoChecksum(e ExchangeMock) (*smocker.Mock, error) {
	return b.build(e, smocker.NewSubstringMatcher("0xdc24316b9ae028f1497c275eb9192a3ea0f67022"))
}

func (b Curve) buildChecksum(e ExchangeMock) (*smocker.Mock, error) {
	return b.build(e, smocker.NewSubstringMatcher("0xDC24316b9AE028F1497c275EB9192a3Ea0f67022"))
}

func (b Curve) build(e ExchangeMock, m smocker.StringMatcher) (*smocker.Mock, error) {
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
