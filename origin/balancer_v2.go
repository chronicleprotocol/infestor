package origin

// Simulate an ETHRPC node returning the price for STETH/ETH on BalancerV2
// https://etherscan.io/address/0x32296969ef14eb0c6d29669c550d4a0449130230#code

import (
	"fmt"
	"math/big"

	"github.com/chronicleprotocol/infestor/smocker"
	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/hexutil"
	"github.com/defiweb/go-eth/types"
)

type BalancerV2 struct {
	EthRPC
}

func (b BalancerV2) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	mocks := make([]*smocker.Mock, 0)

	superMocks, err := b.EthRPC.BuildMocks(e)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, superMocks...)

	m, err := CombineMocks(e, b.buildGetLatest)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	m, err = CombineMocks(e, b.buildWithGetPriceRateCache)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	return mocks, nil
}

func (b BalancerV2) buildGetLatest(e ExchangeMock) (*smocker.Mock, error) {
	// cast sig "getLatest(uint8)(uint256)" == 0xb10be739
	blockNumber, ok := e.Custom["blockNumber"].(int) // Should use same block number with EthRPC exchange
	if !ok {
		return nil, fmt.Errorf("not found block number")
	}
	funcData, ok := e.Custom["getLatest"].([]FunctionData)
	if !ok || len(funcData) < 1 {
		return nil, fmt.Errorf("not found function data for getLatest")
	}

	var calls []MultiCall
	for i := 0; i < len(funcData); i++ {
		data, _ := getLatest.EncodeArgs(funcData[i].Args[0].(byte))
		calls = append(calls, MultiCall{
			Target: funcData[i].Address,
			Data:   data,
		})
	}
	args, _ := encodeMultiCallArgs(calls)
	var data []any
	for i := 0; i < len(funcData); i++ {
		price := funcData[i].Return[0].(*big.Int)
		data = append(data, types.Bytes(price.Bytes()).PadLeft(32))
	}
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

func (b BalancerV2) buildWithGetPriceRateCache(e ExchangeMock) (*smocker.Mock, error) {
	// cast sig "getPriceRateCache(address)(uint256,uint256,uint256)" == 0xb867ee5a
	//                                     rate uint256, duration uint256, expires uint256
	blockNumber, ok := e.Custom["blockNumber"].(int) // Should use same block number with EthRPC exchange
	if !ok {
		return nil, fmt.Errorf("not found block number")
	}
	getLatestFuncData, ok := e.Custom["getLatest"].([]FunctionData)
	if !ok || len(getLatestFuncData) < 1 {
		return nil, fmt.Errorf("not found function data for getLatest")
	}
	getPriceRateCacheFuncData, ok := e.Custom["getPriceRateCache"].([]FunctionData)
	if !ok || len(getPriceRateCacheFuncData) < 1 {
		return nil, nil
	}
	if len(getLatestFuncData) != len(getPriceRateCacheFuncData) {
		return nil, fmt.Errorf("not match function data for getLatest and getPriceRateCache")
	}

	var calls []MultiCall
	var data []any
	for i := 0; i < len(getLatestFuncData); i++ {
		getLatestData, _ := getLatest.EncodeArgs(getLatestFuncData[i].Args[0].(byte))
		getPriceRateCacheData, _ := getPriceRateCache.EncodeArgs(
			getPriceRateCacheFuncData[i].Args[0].(types.Address))
		calls = append(calls, MultiCall{
			Target: getLatestFuncData[i].Address,
			Data:   getLatestData,
		}, MultiCall{
			Target: getPriceRateCacheFuncData[i].Address,
			Data:   getPriceRateCacheData,
		})
		price := getLatestFuncData[i].Return[0].(*big.Int)
		rate := getPriceRateCacheFuncData[i].Return[0].(*big.Int)
		duration := getPriceRateCacheFuncData[i].Return[1].(*big.Int)
		expires := getPriceRateCacheFuncData[i].Return[2].(*big.Int)
		data = append(data, types.Bytes(price.Bytes()).PadLeft(32),
			abi.MustEncodeValues(getPriceRateCache.Outputs(), rate, duration, expires))
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
