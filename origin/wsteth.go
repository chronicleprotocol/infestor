package origin

// Simulate an ETHRPC node returning the price for WstETH
// https://etherscan.io/address/0x7f39C581F595B53c5cb19bD0b3f8dA6c935E2Ca0#code

import (
	"fmt"
	"github.com/defiweb/go-eth/hexutil"
	"github.com/defiweb/go-eth/types"
	"math/big"

	"github.com/chronicleprotocol/infestor/smocker"
)

type WSTETH struct {
	EthRPC
}

func (b WSTETH) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	mocks := make([]*smocker.Mock, 0)

	superMocks, err := b.EthRPC.BuildMocks(e)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, superMocks...)

	m, err := CombineMocks(e, b.buildSTEthPerToken)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	return mocks, nil
}

func (b WSTETH) buildSTEthPerToken(e ExchangeMock) (*smocker.Mock, error) {
	blockNumber, err := e.Custom["blockNumber"].(int)
	if !err {
		return nil, fmt.Errorf("not found block number")
	}
	pool, err := e.Custom[e.Symbol.String()].(types.Address)
	if !err {
		return nil, fmt.Errorf("not found pool address")
	}
	funcData, ok := e.Custom["stEthPerToken"].([]FunctionData)
	if !ok || len(funcData) < 1 {
		return nil, fmt.Errorf("not found function data for stEthPerToken")
	}

	data, _ := stEthPerToken.EncodeArgs()
	calls := []MultiCall{
		{
			Target: pool,
			Data:   data,
		},
	}
	args, _ := encodeMultiCallArgs(calls)
	rate := funcData[0].Return[0].(*big.Int)
	resp, _ := encodeMultiCallResponse(int64(blockNumber), []any{types.Bytes(rate.Bytes()).PadLeft(32)})

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
