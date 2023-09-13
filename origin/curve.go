package origin

// Simulate an ETHRPC node returning the price for STETH/ETH on Curve: get_dy(int128,int128,uint256)
// https://etherscan.io/address/0xdc24316b9ae028f1497c275eb9192a3ea0f67022#code

import (
	"fmt"
	"math/big"

	"github.com/chronicleprotocol/infestor/smocker"
	"github.com/defiweb/go-eth/hexutil"
	"github.com/defiweb/go-eth/types"
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

	m, err = CombineMocks(e, b.buildGetDy1)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	m, err = CombineMocks(e, b.buildGetDy2)
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
	funcData, ok := e.Custom["coins"].([]FunctionData)
	if !ok || len(funcData) < 1 {
		return nil, fmt.Errorf("not found function data for coins")
	}

	var calls []MultiCall
	var data []any
	for _, funcDataItem := range funcData {
		coinsArg, _ := coins.EncodeArgs(funcDataItem.Args[0].(int))

		calls = append(calls, MultiCall{
			Target: funcDataItem.Address,
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
			Body: fmt.Sprintf(RPCJSONResult, hexutil.BytesToHex(resp)),
		},
	}, nil
}

func (b Curve) buildGetDy1(e ExchangeMock) (*smocker.Mock, error) {
	blockNumber, err := e.Custom["blockNumber"].(int) // Should use same block number with EthRPC exchange
	if !err {
		return nil, fmt.Errorf("not found block number")
	}
	funcData, ok := e.Custom["get_dy1"].([]FunctionData)
	if !ok || len(funcData) < 1 {
		return nil, nil
	}

	var calls []MultiCall
	var data []any
	for i := 0; i < len(funcData); i++ {
		getDy1Data, _ := getDy1.EncodeArgs(
			funcData[i].Args[0].(int), funcData[i].Args[1].(int), funcData[i].Args[2].(*big.Int))
		calls = append(calls, MultiCall{
			Target: funcData[i].Address,
			Data:   getDy1Data,
		})
		price := funcData[i].Return[0].(*big.Int)
		data = append(data, types.Bytes(price.Bytes()).PadLeft(32))
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
			Body: fmt.Sprintf(RPCJSONResult, hexutil.BytesToHex(resp)),
		},
	}, nil
}

func (b Curve) buildGetDy2(e ExchangeMock) (*smocker.Mock, error) {
	blockNumber, err := e.Custom["blockNumber"].(int) // Should use same block number with EthRPC exchange
	if !err {
		return nil, fmt.Errorf("not found block number")
	}
	funcData, ok := e.Custom["get_dy2"].([]FunctionData)
	if !ok || len(funcData) < 1 {
		return nil, nil
	}

	var calls []MultiCall
	var data []any
	for i := 0; i < len(funcData); i++ {
		getDy2Data, _ := getDy2.EncodeArgs(
			funcData[i].Args[0].(int), funcData[i].Args[1].(int), funcData[i].Args[2].(*big.Int))
		calls = append(calls, MultiCall{
			Target: funcData[i].Address,
			Data:   getDy2Data,
		})
		price := funcData[i].Return[0].(*big.Int)
		data = append(data, types.Bytes(price.Bytes()).PadLeft(32))
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
			Body: fmt.Sprintf(RPCJSONResult, hexutil.BytesToHex(resp)),
		},
	}, nil
}
