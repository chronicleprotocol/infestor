package origin

import (
	"fmt"
	"github.com/chronicleprotocol/infestor/smocker"
	"github.com/defiweb/go-eth/hexutil"
	"github.com/defiweb/go-eth/types"
	"math/big"
)

type SDAI struct {
	EthRPC
}

func (b SDAI) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	mocks := make([]*smocker.Mock, 0)

	superMocks, err := b.EthRPC.BuildMocks(e)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, superMocks...)

	m, err := CombineMocks(e, b.buildPreviewRedeem)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	return mocks, nil
}

func (b SDAI) buildPreviewRedeem(e ExchangeMock) (*smocker.Mock, error) {
	blockNumber, err := e.Custom["blockNumber"].(int)
	if !err {
		return nil, fmt.Errorf("not found block number")
	}
	sdai, err := e.Custom[e.Symbol.String()].(types.Address)
	if !err {
		return nil, fmt.Errorf("not found sdai address")
	}
	funcData, ok := e.Custom["previewRedeem"].([]FunctionData)
	if !ok || len(funcData) < 1 {
		return nil, fmt.Errorf("not found function data for previewRedeem")
	}

	data, _ := previewRedeem.EncodeArgs(funcData[0].Args[0].(*big.Int))
	calls := []MultiCall{
		{
			Target: sdai,
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
