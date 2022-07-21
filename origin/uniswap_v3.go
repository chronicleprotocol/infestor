package origin

import (
	"fmt"

	"github.com/chronicleprotocol/infestor/smocker"
)

type UniswapV3 struct{}

func (b UniswapV3) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	mocks := make([]*smocker.Mock, 0)
	m, err := CombineMocks(e, b.buildETHUSD)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)
	return mocks, nil
}

func (b UniswapV3) buildETHUSD(e ExchangeMock) (*smocker.Mock, error) {
	m := smocker.ShouldContainSubstring("0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640")
	price := 1588.556212216532373794459636134139
	if e.Price != 0 {
		price = e.Price
	}

	return &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.ShouldEqual("POST"),
			Path:   smocker.ShouldEqual("/subgraphs/name/uniswap/uniswap-v3"),
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
			Body: fmt.Sprintf(`{
  "data": {
    "pools": [
      {
        "token0Price": "%f",
        "token1Price": "0.0001111111111111111111111111111111111"
      }
    ]
  }
}`, price),
		},
	}, nil
}
