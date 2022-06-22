package origin

import (
	"fmt"

	"github.com/chronicleprotocol/infestor/smocker"
)

type Kraken struct{}

func (k Kraken) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	mocksOne, err := CombineMocks(e, k.build)
	if err != nil {
		return nil, err
	}
	other, err := CombineMocks(e, k.buildWithSlash)
	if err != nil {
		return nil, err
	}
	return append(mocksOne, other...), nil
}

func (k Kraken) buildWithSlash(e ExchangeMock) (*smocker.Mock, error) {
	symbol := e.Symbol.Format("%s/%s")
	body := `{
 "error": [],
 "result": {
	 "%s": {
		 "a": [
			 "%f",
			 "10",
			 "10.000"
		 ],
		 "b": [
			 "%f",
			 "2",
			 "2.000"
		 ],
		 "c": [
			 "%f",
			 "0.22651150"
		 ],
		 "v": [
			 "%f",
			 "5803.57144830"
		 ]
	 }
 }
}`

	return &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.ShouldEqual("GET"),
			Path:   smocker.ShouldEqual("/0/public/Ticker"),
			QueryParams: map[string]smocker.StringMatcherSlice{
				"pair": []smocker.StringMatcher{
					smocker.ShouldEqual(symbol),
				},
			},
		},
		Response: &smocker.MockResponse{
			Status: e.StatusCode,
			Headers: map[string]smocker.StringSlice{
				"Content-Type": []string{
					"application/json",
				},
			},
			Body: fmt.Sprintf(body, symbol, e.Ask, e.Bid, e.Price, e.Volume),
		},
	}, nil
}

func (k Kraken) build(e ExchangeMock) (*smocker.Mock, error) {
	symbol := e.Symbol.Format("%s%s")
	body := `{
 "error": [],
 "result": {
	 "%s": {
		 "a": [
			 "%f",
			 "10",
			 "10.000"
		 ],
		 "b": [
			 "%f",
			 "2",
			 "2.000"
		 ],
		 "c": [
			 "%f",
			 "0.22651150"
		 ],
		 "v": [
			 "%f",
			 "5803.57144830"
		 ]
	 }
 }
}`

	return &smocker.Mock{
		Request: smocker.MockRequest{
			Method: smocker.ShouldEqual("GET"),
			Path:   smocker.ShouldEqual("/0/public/Ticker"),
			QueryParams: map[string]smocker.StringMatcherSlice{
				"pair": []smocker.StringMatcher{
					smocker.ShouldEqual(symbol),
				},
			},
		},
		Response: &smocker.MockResponse{
			Status: e.StatusCode,
			Headers: map[string]smocker.StringSlice{
				"Content-Type": []string{
					"application/json",
				},
			},
			Body: fmt.Sprintf(body, symbol, e.Ask, e.Bid, e.Price, e.Volume),
		},
	}, nil
}
