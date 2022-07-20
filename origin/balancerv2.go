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

	m, err := CombineMocks(e, b.buildGetPriceRateCache)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	m, err = CombineMocks(e, b.buildGetLatest)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	return mocks, nil
}

func (b BalancerV2) buildGetLatest(e ExchangeMock) (*smocker.Mock, error) {
	price, ok := e.Custom["price"]
	if !ok {
		return nil, fmt.Errorf("`price` custom field is requierd for balancerV2")
	}

	// getLatest(uint256)
	m := smocker.ShouldContainSubstring("0xb10be739")

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
			Body: fmt.Sprintf(rpcJSONResult, price),
		},
	}, nil
}

func (b BalancerV2) buildGetPriceRateCache(e ExchangeMock) (*smocker.Mock, error) {
	rate, ok := e.Custom["rate"]
	if !ok {
		return nil, fmt.Errorf("`rate` custom field is requierd for balancerV2")
	}

	// getPriceRateCache(uint256,uint256,uint256)
	m := smocker.ShouldContainSubstring("0xb867ee5a")

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
			Body: fmt.Sprintf(rpcJSONResult, rate),
		},
	}, nil
}
