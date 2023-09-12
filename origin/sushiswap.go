package origin

import (
	"fmt"
	"github.com/chronicleprotocol/infestor/smocker"
	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/hexutil"
	"github.com/defiweb/go-eth/types"
	"math/big"
)

type Sushiswap struct {
	EthRPC
}

func (b Sushiswap) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	mocks := make([]*smocker.Mock, 0)

	superMocks, err := b.EthRPC.BuildMocks(e)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, superMocks...)

	m, err := CombineMocks(e, b.buildGetReserves)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	m, err = CombineMocks(e, b.buildToken0)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	return mocks, nil
}

func (b Sushiswap) buildGetReserves(e ExchangeMock) (*smocker.Mock, error) {
	blockNumber, err := e.Custom["blockNumber"].(int)
	if !err {
		return nil, fmt.Errorf("not found block number")
	}
	pool, err := e.Custom[e.Symbol.String()].(types.Address)
	if !err {
		return nil, fmt.Errorf("not found pool address")
	}
	funcData, ok := e.Custom["getReserves"].([]FunctionData)
	if !ok || len(funcData) < 1 {
		return nil, fmt.Errorf("not found function data for getReserves")
	}

	data, _ := getReserves.EncodeArgs()
	calls := []MultiCall{
		{
			Target: pool,
			Data:   data,
		},
	}
	args, _ := encodeMultiCallArgs(calls)
	reserve0 := funcData[0].Return[0].(*big.Int)
	reserve1 := funcData[0].Return[1].(*big.Int)
	blockTimestamp := funcData[0].Return[2].(*big.Int)
	resp, _ := encodeMultiCallResponse(int64(blockNumber),
		[]any{abi.MustEncodeValues(getReserves.Outputs(), reserve0, reserve1, blockTimestamp)})

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

func (b Sushiswap) buildToken0(e ExchangeMock) (*smocker.Mock, error) {
	blockNumber, err := e.Custom["blockNumber"].(int)
	if !err {
		return nil, fmt.Errorf("not found block number")
	}
	pool, err := e.Custom[e.Symbol.String()].(types.Address)
	if !err {
		return nil, fmt.Errorf("not found pool address")
	}
	token0FuncData, ok := e.Custom["token0"].([]FunctionData)
	if !ok || len(token0FuncData) < 1 {
		return nil, fmt.Errorf("not found function data for token0")
	}
	token1FuncData, ok := e.Custom["token1"].([]FunctionData)
	if !ok || len(token1FuncData) < 1 {
		return nil, fmt.Errorf("not found function data for token1")
	}

	token0Args, _ := token0Abi.EncodeArgs()
	token1Args, _ := token1Abi.EncodeArgs()
	calls := []MultiCall{
		{
			Target: pool,
			Data:   token0Args,
		},
		{
			Target: pool,
			Data:   token1Args,
		},
	}
	args, _ := encodeMultiCallArgs(calls)
	token0 := token0FuncData[0].Return[0].(types.Address)
	token1 := token1FuncData[0].Return[0].(types.Address)
	resp, _ := encodeMultiCallResponse(int64(blockNumber), []any{
		types.Bytes(token0.Bytes()).PadLeft(32),
		types.Bytes(token1.Bytes()).PadLeft(32),
	})

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
