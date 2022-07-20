package origin

import (
	"fmt"

	"github.com/chronicleprotocol/infestor/smocker"
)

type RocketPool struct{}

func (b RocketPool) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	mocks := make([]*smocker.Mock, 0)
	m, err := CombineMocks(e, b.buildGetExchangeRate)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	m, err = CombineMocks(e, b.buildGetRethValue)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	m, err = CombineMocks(e, b.buildGetCircuitBreakerValue)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	m, err = CombineMocks(e, b.buildGetCircuitBreakerDivisor)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	return mocks, nil
}

func (b RocketPool) buildGetExchangeRate(e ExchangeMock) (*smocker.Mock, error) {
	// getExchangeRate
	m := smocker.ShouldContainSubstring("0xe6aa216c")

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
			Body: fmt.Sprintf(rpcJSONResult, "0x0000000000000000000000000000000000000000000000000de0b6b3a7640000"),
		},
	}, nil
}

func (b RocketPool) buildGetRethValue(e ExchangeMock) (*smocker.Mock, error) {
	// getRethValue(uint256)
	m := smocker.ShouldContainSubstring("0x4346f03e")

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
			Body: fmt.Sprintf(rpcJSONResult, "0x0000000000000000000000000000000000000000000000000de0b6b3a7640000"),
		},
	}, nil
}

func (b RocketPool) buildGetCircuitBreakerValue(e ExchangeMock) (*smocker.Mock, error) {
	// read()
	m := smocker.ShouldContainSubstring("0x57de26a4")

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
			// Set to 10%
			Body: fmt.Sprintf(rpcJSONResult, "0x0000000000000000000000000000000000000000000000000000000000002710"),
		},
	}, nil
}

func (b RocketPool) buildGetCircuitBreakerDivisor(e ExchangeMock) (*smocker.Mock, error) {
	// divisor()
	m := smocker.ShouldContainSubstring("0x1f2dc5ef")

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
			// 10^5
			Body: fmt.Sprintf(rpcJSONResult, "0x00000000000000000000000000000000000000000000000000000000000186a0"),
		},
	}, nil
}
