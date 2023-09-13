package origin

import (
	"fmt"
	"math/big"

	"github.com/chronicleprotocol/infestor/smocker"
	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/hexutil"
	"github.com/defiweb/go-eth/types"
)

type UniswapV2 struct {
	EthRPC
}

func (b UniswapV2) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
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

func (b UniswapV2) buildGetReserves(e ExchangeMock) (*smocker.Mock, error) {
	blockNumber, err := e.Custom["blockNumber"].(int)
	if !err {
		return nil, fmt.Errorf("not found block number")
	}
	funcData, ok := e.Custom["getReserves"].([]FunctionData)
	if !ok || len(funcData) < 1 {
		return nil, fmt.Errorf("not found function data for getReserves")
	}

	var calls []MultiCall
	var data []any
	for i := 0; i < len(funcData); i++ {
		getReservesData, _ := getReserves.EncodeArgs()
		calls = append(calls, MultiCall{
			Target: funcData[i].Address,
			Data:   getReservesData,
		})
		reserve0 := funcData[i].Return[0].(*big.Int)
		reserve1 := funcData[i].Return[1].(*big.Int)
		blockTimestamp := funcData[i].Return[2].(*big.Int)
		data = append(data, abi.MustEncodeValues(getReserves.Outputs(), reserve0, reserve1, blockTimestamp))
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

func (b UniswapV2) buildToken0(e ExchangeMock) (*smocker.Mock, error) {
	blockNumber, err := e.Custom["blockNumber"].(int)
	if !err {
		return nil, fmt.Errorf("not found block number")
	}
	token0FuncData, ok := e.Custom["token0"].([]FunctionData)
	if !ok || len(token0FuncData) < 1 {
		return nil, fmt.Errorf("not found function data for token0")
	}
	token1FuncData, ok := e.Custom["token1"].([]FunctionData)
	if !ok || len(token1FuncData) < 1 {
		return nil, fmt.Errorf("not found function data for token1")
	}
	if len(token0FuncData) != len(token1FuncData) {
		return nil, fmt.Errorf("not found function data for token0 and token1")
	}

	token0Data, _ := token0Abi.EncodeArgs()
	token1Data, _ := token1Abi.EncodeArgs()
	var calls []MultiCall
	var data []any
	for i := 0; i < len(token0FuncData); i++ {
		calls = append(calls, MultiCall{
			Target: token0FuncData[i].Address,
			Data:   token0Data,
		}, MultiCall{
			Target: token1FuncData[i].Address,
			Data:   token1Data,
		})
		token0 := token0FuncData[i].Return[0].(types.Address)
		token1 := token1FuncData[i].Return[0].(types.Address)
		data = append(data, types.Bytes(token0.Bytes()).PadLeft(32),
			types.Bytes(token1.Bytes()).PadLeft(32))
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
