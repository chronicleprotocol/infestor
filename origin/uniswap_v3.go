package origin

import (
	"fmt"
	"math/big"

	"github.com/chronicleprotocol/infestor/smocker"
	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/hexutil"
	"github.com/defiweb/go-eth/types"
)

type UniswapV3 struct {
	EthRPC
}

func (b UniswapV3) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	mocks := make([]*smocker.Mock, 0)

	superMocks, err := b.EthRPC.BuildMocks(e)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, superMocks...)

	m, err := CombineMocks(e, b.buildSlot0)
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

func (b UniswapV3) buildSlot0(e ExchangeMock) (*smocker.Mock, error) {
	blockNumber, err := e.Custom["blockNumber"].(int)
	if !err {
		return nil, fmt.Errorf("not found block number")
	}
	funcData, ok := e.Custom["slot0"].([]FunctionData)
	if !ok || len(funcData) < 1 {
		return nil, fmt.Errorf("not found function data for slot0")
	}

	var calls []MultiCall
	var data []any
	for i := 0; i < len(funcData); i++ {
		slot0Data, _ := slot0.EncodeArgs()
		calls = append(calls, MultiCall{
			Target: funcData[i].Address,
			Data:   slot0Data,
		})
		sqrtPriceX96 := funcData[i].Return[0].(*big.Int)
		tick := funcData[i].Return[1].(*big.Int)
		observationIndex := funcData[i].Return[2].(*big.Int)
		observationCardinality := funcData[i].Return[3].(*big.Int)
		observationCardinalityNext := funcData[i].Return[4].(*big.Int)
		feeProtocol := funcData[i].Return[5].(int)
		unlocked := funcData[i].Return[6].(bool)
		slot0Bytes := abi.MustEncodeValues(slot0.Outputs(),
			sqrtPriceX96,
			tick,
			observationIndex,
			observationCardinality,
			observationCardinalityNext,
			feeProtocol,
			unlocked)
		data = append(data, slot0Bytes)
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

func (b UniswapV3) buildToken0(e ExchangeMock) (*smocker.Mock, error) {
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
