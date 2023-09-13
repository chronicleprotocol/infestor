package origin

import (
	"fmt"
	"math/big"

	"github.com/chronicleprotocol/infestor/smocker"
	"github.com/defiweb/go-eth/hexutil"
	"github.com/defiweb/go-eth/types"
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
	funcData, ok := e.Custom["previewRedeem"].([]FunctionData)
	if !ok || len(funcData) < 1 {
		return nil, fmt.Errorf("not found function data for previewRedeem")
	}

	var calls []MultiCall
	var data []any
	for i := 0; i < len(funcData); i++ {
		previewRedeemData, _ := previewRedeem.EncodeArgs(funcData[i].Args[0].(*big.Int))
		calls = append(calls, MultiCall{
			Target: funcData[i].Address,
			Data:   previewRedeemData,
		})
		rate := funcData[i].Return[0].(*big.Int)
		data = append(data, types.Bytes(rate.Bytes()).PadLeft(32))
	}
	args, _ := encodeMultiCallArgs(calls)
	resp, _ := encodeMultiCallResponse(int64(blockNumber), data)

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
			Body: fmt.Sprintf(RPCJSONResult, hexutil.BytesToHex(resp)),
		},
	}, nil
}
