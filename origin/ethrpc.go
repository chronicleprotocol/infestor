package origin

// Simulate an ETHRPC node. Add any general ETH/RPC related mocks here.

import (
	"fmt"
	"github.com/chronicleprotocol/infestor/smocker"
	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/hexutil"
	"github.com/defiweb/go-eth/types"
	"math/big"
	"strconv"
)

const RpcJSONResult = `{
  "jsonrpc": "2.0",
  "id": 1,
  "result": "%s"
}`

const RpcCallRequestJSON = `{"method":"eth_call","params":[{"from":"%s","to":"%s","data":"%s"},"%s"],"id":1,"jsonrpc":"2.0"}`

type EthRPC struct{}

func (b EthRPC) BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error) {
	mocks := make([]*smocker.Mock, 0)
	m, err := CombineMocks(e, b.buildChainId)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	m, err = CombineMocks(e, b.buildBlockNumber)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	m, err = CombineMocks(e, b.buildNetVersion)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	m, err = CombineMocks(e, b.buildSymbols)
	if err != nil {
		return nil, err
	}
	mocks = append(mocks, m...)

	return mocks, nil
}

func (b EthRPC) buildChainId(e ExchangeMock) (*smocker.Mock, error) {
	m := smocker.ShouldContainSubstring("eth_chainId")

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
			Body: `{"jsonrpc":"2.0","id":1,"result":"1"}`,
		},
	}, nil
}

func (b EthRPC) buildBlockNumber(e ExchangeMock) (*smocker.Mock, error) {
	blockNumber, err := e.Custom["blockNumber"].(int)
	if !err {
		return nil, fmt.Errorf("not found pool address")
	}

	blockNumberHex := strconv.FormatInt(int64(blockNumber), 16)

	m := smocker.ShouldContainSubstring("eth_blockNumber")

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
			Body: fmt.Sprintf(`{"jsonrpc":"2.0","id":1,"result":"%s"}`, blockNumberHex),
		},
	}, nil
}

func (b EthRPC) buildNetVersion(e ExchangeMock) (*smocker.Mock, error) {
	m := smocker.ShouldContainSubstring("net_version")

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
			Body: `{"jsonrpc":"2.0","id":1,"result":"1"}`,
		},
	}, nil
}

func (b EthRPC) buildSymbols(e ExchangeMock) (*smocker.Mock, error) {
	blockNumber, ok := e.Custom["blockNumber"].(int) // Should use same block number with EthRPC exchange
	if !ok {
		return nil, fmt.Errorf("not found block number")
	}
	tokens, err := e.Custom["tokens"].([]types.Address)
	if !err {
		return nil, nil
	}
	symbols, err := e.Custom["symbols"].([]FunctionData)
	if !err {
		return nil, fmt.Errorf("not found return values for symbols")
	}
	decimals, err := e.Custom["decimals"].([]FunctionData)
	if !err {
		return nil, fmt.Errorf("not found return values for decimals")
	}

	var calls []MultiCall
	for _, token := range tokens {
		symbolArg, _ := getSymbol.EncodeArgs()
		decimalArgs, _ := getDecimals.EncodeArgs()

		calls = append(calls, MultiCall{
			Target: token,
			Data:   symbolArg,
		}, MultiCall{
			Target: token,
			Data:   decimalArgs,
		})
	}

	var data []any
	for i := 0; i < len(symbols); i++ {
		symbol := symbols[i].Return[0].(string)
		symbolAbi := abi.MustParseType("(string memory)")
		symbolMap := make(map[string]string)
		symbolMap["arg0"] = symbol
		decimal := decimals[i].Return[0].(*big.Int)
		decimalBytes := types.Bytes(decimal.Bytes()).PadLeft(32)

		data = append(data, abi.MustEncodeValue(symbolAbi, symbolMap), decimalBytes)
	}
	args, _ := encodeMultiCallArgs(calls)
	resp, _ := encodeMultiCallResponse(int64(blockNumber), data)

	fmt.Println("symbols, args", hexutil.BytesToHex(args))
	fmt.Println("symbols, resp", hexutil.BytesToHex(resp))

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
