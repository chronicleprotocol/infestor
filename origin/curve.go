package origin

// Simulate an ETHRPC node returning the price for STETH/ETH on Curve: get_dy(int128,int128,uint256)
// https://etherscan.io/address/0xdc24316b9ae028f1497c275eb9192a3ea0f67022#code

import (
	"fmt"
	"github.com/chronicleprotocol/infestor/smocker"
	"github.com/defiweb/go-eth/hexutil"
	"github.com/defiweb/go-eth/types"
	"math/big"
)

type Curve struct {
	EthRPC
}

func (b Curve) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	mocks := make([]*smocker.Mock, 0)

	superMocks, err := b.EthRPC.BuildMocks(e)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, superMocks...)

	m, err := CombineMocks(e, b.buildCoins)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	m, err = CombineMocks(e, b.buildGetDy)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	return mocks, nil
}

func (b Curve) buildCoins(e ExchangeMock) (*smocker.Mock, error) {
	blockNumber, err := e.Custom["blockNumber"].(int)
	if !err {
		return nil, fmt.Errorf("not found block number")
	}
	pool, err := e.Custom[e.Symbol.String()].(types.Address)
	if !err {
		return nil, fmt.Errorf("not found pool address")
	}
	funcData, ok := e.Custom["coins"].([]FunctionData)
	if !ok || len(funcData) < 1 {
		return nil, fmt.Errorf("not found function data for coins")
	}

	var calls []MultiCall
	var data []any
	for _, funcDataItem := range funcData {
		coinsArg, _ := coins.EncodeArgs(funcDataItem.Args[0].(int))

		calls = append(calls, MultiCall{
			Target: pool,
			Data:   coinsArg,
		})
		data = append(data, types.Bytes(funcDataItem.Return[0].(types.Address).Bytes()).PadLeft(32))
	}

	args, _ := encodeMultiCallArgs(calls)
	resp, _ := encodeMultiCallResponse(int64(blockNumber), data)

	m := smocker.ShouldContainSubstring(hexutil.BytesToHex(args))

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
			Body: fmt.Sprintf(RpcJSONResult, hexutil.BytesToHex(resp)),
		},
	}, nil
}

func (b Curve) buildGetDy(e ExchangeMock) (*smocker.Mock, error) {
	blockNumber, err := e.Custom["blockNumber"].(int) // Should use same block number with EthRPC exchange
	if !err {
		return nil, fmt.Errorf("not found block number")
	}
	pool, err := e.Custom[e.Symbol.String()].(types.Address)
	if !err {
		return nil, fmt.Errorf("not found pool address")
	}
	funcData, ok := e.Custom["get_dy"].([]FunctionData)
	if !ok || len(funcData) < 1 {
		return nil, fmt.Errorf("not found function data for getDy")
	}

	data, _ := getDy1.EncodeArgs(funcData[0].Args[0].(int), funcData[0].Args[1].(int), funcData[0].Args[2].(*big.Int))
	calls := []MultiCall{
		{
			Target: pool,
			Data:   data,
		},
	}
	args, _ := encodeMultiCallArgs(calls)
	price := funcData[0].Return[0].(*big.Int)
	resp, _ := encodeMultiCallResponse(int64(blockNumber), []any{types.Bytes(price.Bytes()).PadLeft(32)})

	m := smocker.ShouldContainSubstring(hexutil.BytesToHex(args))

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
			Body: fmt.Sprintf(RpcJSONResult, hexutil.BytesToHex(resp)),
		},
	}, nil
}
