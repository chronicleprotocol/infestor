package origin

// Simulate an ETHRPC node returning the price for STETH/ETH on BalancerV2
// https://etherscan.io/address/0x32296969ef14eb0c6d29669c550d4a0449130230#code

import (
	"fmt"
	"github.com/chronicleprotocol/infestor/smocker"
	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/hexutil"
	"github.com/defiweb/go-eth/types"
	"math/big"
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

	m, err = CombineMocks(e, b.buildGetPriceRateCache)
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
		return nil, fmt.Errorf("not added ethrpc exchange")
	}
	pool, ok := e.Custom["pool"].(types.Address)
	if !ok {
		return nil, fmt.Errorf("not found pool address")
	}
	funcData, ok := e.Custom["getLatest"].([]FunctionData)
	if !ok || len(funcData) < 1 {
		return nil, fmt.Errorf("not found function data for getLatest")
	}

	data, _ := getLatest.EncodeArgs(funcData[0].Args[0].(byte))

	calls := []MultiCall{
		{
			Target: pool,
			Data:   data,
		},
	}
	args, _ := encodeMultiCallArgs(calls)
	price := funcData[0].Return[0].(*big.Int)
	resp, _ := encodeMultiCallResponse(int64(blockNumber), []any{types.Bytes(price.Bytes()).PadLeft(32)})

	fmt.Println("args", hexutil.BytesToHex(args))
	fmt.Println("resp", hexutil.BytesToHex(resp))

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

func (b BalancerV2) buildGetPriceRateCache(e ExchangeMock) (*smocker.Mock, error) {
	// cast sig "getPriceRateCache(address)(uint256,uint256,uint256)" == 0xb867ee5a
	//                                     rate uint256, duration uint256, expires uint256
	pool, err := e.Custom["pool"].(types.Address)
	if !err {
		return nil, fmt.Errorf("not found pool address")
	}
	funcData, ok := e.Custom["getPriceRateCache"].([]FunctionData)
	if !ok || len(funcData) < 1 {
		return nil, nil
	}

	data, _ := getPriceRateCache.EncodeArgs(funcData[0].Args[0].(types.Address))

	calls := []MultiCall{
		{
			Target: pool,
			Data:   data,
		},
	}
	args, _ := encodeMultiCallArgs(calls)
	rate := funcData[0].Return[0].(*big.Int)
	duration := funcData[0].Return[1].(*big.Int)
	expires := funcData[0].Return[2].(*big.Int)
	resp, _ := encodeMultiCallResponse(100,
		[]any{abi.MustEncodeValues(getPriceRateCache.Outputs(),
			types.Bytes(rate.Bytes()).PadLeft(32),     // rate
			types.Bytes(duration.Bytes()).PadLeft(32), // duration
			types.Bytes(expires.Bytes()).PadLeft(32),  // expires
		)})

	fmt.Println("args", hexutil.BytesToHex(args))
	fmt.Println("resp", hexutil.BytesToHex(resp))

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
