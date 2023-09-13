package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/defiweb/go-eth/types"
	"github.com/stretchr/testify/suite"

	"github.com/chronicleprotocol/infestor"
	"github.com/chronicleprotocol/infestor/origin"
	"github.com/chronicleprotocol/infestor/smocker"
)

func parseBody(resp *http.Response, r interface{}) error {
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, r)
}

func TestExchangesE2ESuite(t *testing.T) {
	suite.Run(t, new(ExchangesE2ESuite))
}

type ExchangesE2ESuite struct {
	suite.Suite
	api smocker.API
	url string
}

type jsonrpcMessage struct {
	Version string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  string          `json:"result,omitempty"`
}

func (s *ExchangesE2ESuite) SetupSuite() {
	smockerHost, exist := os.LookupEnv("SMOCKER_HOST")
	s.Require().True(exist, "SMOCKER_HOST env variable have to be set")

	s.api = smocker.API{URL: smockerHost + ":8081"}

	s.url = fmt.Sprintf("%s:8080", smockerHost)
}

func (s *ExchangesE2ESuite) SetupTest() {
	err := s.api.Reset(context.Background())
	s.Require().NoError(err)
}

func (s *ExchangesE2ESuite) TestEthRPC() {
	const blockNumber int = 100
	ex := origin.NewExchange("ethrpc").
		WithCustom("blockNumber", blockNumber).
		WithCustom("tokens",
			[]types.Address{
				types.MustAddressFromHex("0xae7ab96520DE3A18E5e111B5EaAb095312D7fE84"),
				types.MustAddressFromHex("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
			}).
		WithFunctionData("symbols", []origin.FunctionData{
			{
				Args:   []any{},
				Return: []any{"stETH"},
			},
			{
				Args:   []any{},
				Return: []any{"WETH"},
			},
		}).
		WithFunctionData("decimals", []origin.FunctionData{
			{
				Args:   []any{},
				Return: []any{big.NewInt(18)},
			},
			{
				Args:   []any{},
				Return: []any{big.NewInt(18)},
			},
		})
	err := infestor.NewMocksBuilder().Reset().Debug().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/", s.url)
	reqJSON := fmt.Sprintf(origin.RPCCallRequestJSON, "null", "0xeefba1e63905ef1d7acba5a8513c70307c1ce441", "0x252dba42000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001800000000000000000000000000000000000000000000000000000000000000200000000000000000000000000ae7ab96520de3a18e5e111b5eaab095312d7fe840000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000495d89b4100000000000000000000000000000000000000000000000000000000000000000000000000000000ae7ab96520de3a18e5e111b5eaab095312d7fe8400000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000004313ce56700000000000000000000000000000000000000000000000000000000000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc20000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000495d89b4100000000000000000000000000000000000000000000000000000000000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000004313ce56700000000000000000000000000000000000000000000000000000000", "latest")
	jsonStr := []byte(fmt.Sprint(reqJSON))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var response jsonrpcMessage
	err = parseBody(resp, &response)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	// Multicall response should be equal to the expected one
	s.Require().Equal(response.Result, "0x00000000000000000000000000000000000000000000000000000000000000640000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000014000000000000000000000000000000000000000000000000000000000000001c0000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000005737445544800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000004574554480000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000012")

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)
	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
	defer func() { _ = resp.Body.Close() }()
}

func (s *ExchangesE2ESuite) TestBalancerV2() {
	ex := origin.NewExchange("balancerV2").WithSymbol("STETH/ETH").WithPrice(1)
	// No `pool` field should fail
	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().Error(err)

	const blockNumber int = 100
	price := big.NewInt(0.94 * 1e18)
	ex = ex.
		WithCustom("STETH/ETH", types.MustAddressFromHex("0x32296969ef14eb0c6d29669c550d4a0449130230")).
		WithFunctionData("getLatest", []origin.FunctionData{
			{
				Args:   []any{byte(0)},
				Return: []any{price},
			},
		}).
		WithCustom("blockNumber", blockNumber).
		WithPrice(0.94)
	err = infestor.NewMocksBuilder().Debug().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/", s.url)

	// Calling RPC url via multicall contract
	reqJSON := fmt.Sprintf(origin.RPCCallRequestJSON, "null", "0xeefba1e63905ef1d7acba5a8513c70307c1ce441", "0x252dba4200000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000032296969ef14eb0c6d29669c550d4a044913023000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000024b10be739000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", "latest")
	jsonStr := []byte(fmt.Sprint(reqJSON))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var response jsonrpcMessage
	err = parseBody(resp, &response)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	// Multicall response should be equal to the expected one
	s.Require().Equal(response.Result, "0x000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000d0b8d0508de0000")

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)
	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
	defer func() { _ = resp.Body.Close() }()
}

func (s *ExchangesE2ESuite) TestBinance() {
	ex := origin.NewExchange("binance").WithSymbol("ETH/BTC").WithPrice(1)

	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/api/v3/ticker/price?symbol=ETHBTC", s.url)
	resp, err := http.Get(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	binanceResponse := struct {
		Symbol string
		Price  string
	}{}
	err = parseBody(resp, &binanceResponse)

	s.Require().NoError(err)
	s.Require().Equal("ETHBTC", binanceResponse.Symbol)
	s.Require().Equal("1.00000000", binanceResponse.Price)

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)

	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *ExchangesE2ESuite) TestBitfinex() {
	ex := origin.NewExchange("bitfinex").
		WithSymbol("ETH/BTC").
		WithPrice(1).
		WithAsk(2).
		WithBid(3).
		WithVolume(4)

	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/v2/ticker/tETHBTC", s.url)
	resp, err := http.Get(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var response []float64
	err = parseBody(resp, &response)

	s.Require().NoError(err)
	s.Require().Len(response, 10)
	s.Require().Equal(float64(3), response[0])
	s.Require().Equal(float64(2), response[2])
	s.Require().Equal(float64(1), response[6])
	s.Require().Equal(float64(4), response[7])

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)

	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *ExchangesE2ESuite) TestBitStamp() {
	ts := time.Now()
	ex := origin.NewExchange("bitstamp").
		WithSymbol("ETH/BTC").
		WithPrice(1).
		WithVolume(2).
		WithBid(3).
		WithAsk(4).
		WithTime(ts)

	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/api/v2/ticker/ethbtc", s.url)
	resp, err := http.Get(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	response := struct {
		Last      string
		Bid       string
		Ask       string
		Volume    string
		Timestamp string
	}{}

	err = parseBody(resp, &response)

	s.Require().NoError(err)
	s.Require().Equal("1.000000", response.Last)
	s.Require().Equal("2.000000", response.Volume)
	s.Require().Equal("3.000000", response.Bid)
	s.Require().Equal("4.000000", response.Ask)
	s.Require().Equal(fmt.Sprintf("%d", ts.Unix()), response.Timestamp)

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)

	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *ExchangesE2ESuite) TestCoinbase() {
	ts := time.Now()
	format := "2006-01-02T15:04:05.999999Z"

	ex := origin.NewExchange("coinbase").
		WithSymbol("ETH/BTC").
		WithPrice(1).
		WithVolume(2).
		WithBid(3).
		WithAsk(4).
		WithTime(ts)

	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/products/ETH-BTC/ticker", s.url)
	resp, err := http.Get(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	response := struct {
		Price  string
		Bid    string
		Ask    string
		Volume string
		Time   string
	}{}

	err = parseBody(resp, &response)

	s.Require().NoError(err)
	s.Require().Equal("1.000000", response.Price)
	s.Require().Equal("2.000000", response.Volume)
	s.Require().Equal("3.000000", response.Bid)
	s.Require().Equal("4.000000", response.Ask)
	s.Require().Equal(ts.UTC().Format(format), response.Time)

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)

	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *ExchangesE2ESuite) TestCurve() {
	ex := origin.NewExchange("curve").WithSymbol("ETH/STETH").WithPrice(1)
	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().Error(err)

	const blockNumber int = 100
	price := big.NewInt(0.94 * 1e18)
	ex = ex.
		WithCustom("ETH/STETH", types.MustAddressFromHex("0xDC24316b9AE028F1497c275EB9192a3Ea0f67022")).
		WithFunctionData("coins", []origin.FunctionData{
			{
				Args:   []any{0},
				Return: []any{types.MustAddressFromHex("0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee")},
			},
			{
				Args:   []any{1},
				Return: []any{types.MustAddressFromHex("0xae7ab96520de3a18e5e111b5eaab095312d7fe84")},
			},
		}).
		WithCustom("tokens", []types.Address{
			// types.MustAddressFromHex("0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"), // do not call for ETH
			types.MustAddressFromHex("0xae7ab96520de3a18e5e111b5eaab095312d7fe84"),
		}).
		WithFunctionData("symbols", []origin.FunctionData{
			// { // do not call for ETH
			//	Args:   []any{},
			//	Return: []any{"ETH"},
			// },
			{
				Args:   []any{},
				Return: []any{"stETH"},
			},
		}).
		WithFunctionData("decimals", []origin.FunctionData{
			// { // do not call for ETH
			//	Args:   []any{},
			//	Return: []any{big.NewInt(18)},
			// },
			{
				Args:   []any{},
				Return: []any{big.NewInt(18)},
			},
		}).
		WithFunctionData("get_dy", []origin.FunctionData{
			{
				Args:   []any{0, 1, big.NewInt(1e18)},
				Return: []any{price},
			},
		}).
		WithCustom("blockNumber", blockNumber).
		WithPrice(0.94)
	err = infestor.NewMocksBuilder().Debug().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/", s.url)

	// Calling RPC url via multicall contract, coins()
	reqJSON := fmt.Sprintf(origin.RPCCallRequestJSON, "null", "0xeefba1e63905ef1d7acba5a8513c70307c1ce441", "0x252dba4200000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000dc24316b9ae028f1497c275eb9192a3ea0f6702200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000024c6610657000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000dc24316b9ae028f1497c275eb9192a3ea0f6702200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000024c6610657000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000", "latest")
	jsonStr := []byte(fmt.Sprint(reqJSON))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var response jsonrpcMessage
	err = parseBody(resp, &response)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	// Multicall response should be equal to the expected one
	s.Require().Equal(response.Result, "0x000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000020000000000000000000000000eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000ae7ab96520de3a18e5e111b5eaab095312d7fe84")

	// symbol(), decimals()
	reqJSON = fmt.Sprintf(origin.RPCCallRequestJSON, "null", "0xeefba1e63905ef1d7acba5a8513c70307c1ce441", "0x252dba4200000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000ae7ab96520de3a18e5e111b5eaab095312d7fe840000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000495d89b4100000000000000000000000000000000000000000000000000000000000000000000000000000000ae7ab96520de3a18e5e111b5eaab095312d7fe8400000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000004313ce56700000000000000000000000000000000000000000000000000000000", "latest")
	jsonStr = []byte(fmt.Sprint(reqJSON))
	resp, err = http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	err = parseBody(resp, &response)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	s.Require().Equal(response.Result, "0x000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000005737445544800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000012")

	// get_dy()
	reqJSON = fmt.Sprintf(origin.RPCCallRequestJSON, "null", "0xeefba1e63905ef1d7acba5a8513c70307c1ce441", "0x252dba42000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000020000000000000000000000000dc24316b9ae028f1497c275eb9192a3ea0f67022000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000645e0d443f000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000de0b6b3a764000000000000000000000000000000000000000000000000000000000000", "latest")
	jsonStr = []byte(fmt.Sprint(reqJSON))
	resp, err = http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	err = parseBody(resp, &response)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	s.Require().Equal(response.Result, "0x000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000d0b8d0508de0000")

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)
	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
	defer func() { _ = resp.Body.Close() }()
}

func (s *ExchangesE2ESuite) TestDSR() {
	ex := origin.NewExchange("dsr").WithSymbol("DSR/RATE").WithPrice(1)
	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().Error(err)

	const blockNumber int = 100
	rate := new(big.Int).Mul(big.NewInt(102), new(big.Int).Exp(big.NewInt(10), big.NewInt(25), nil))
	ex = ex.
		WithCustom("DSR/RATE", types.MustAddressFromHex("0x197E90f9FAD81970bA7976f33CbD77088E5D7cf7")).
		WithFunctionData("dsr", []origin.FunctionData{
			{
				Args:   []any{},
				Return: []any{rate},
			},
		}).
		WithCustom("blockNumber", blockNumber).
		WithPrice(1.02)
	err = infestor.NewMocksBuilder().Debug().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/", s.url)

	// Calling RPC url via multicall contract, coins()
	reqJSON := fmt.Sprintf(origin.RPCCallRequestJSON, "null", "0xeefba1e63905ef1d7acba5a8513c70307c1ce441", "0x252dba42000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000020000000000000000000000000197e90f9fad81970ba7976f33cbd77088e5d7cf700000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000004487bf08200000000000000000000000000000000000000000000000000000000", "latest")
	jsonStr := []byte(fmt.Sprint(reqJSON))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var response jsonrpcMessage
	err = parseBody(resp, &response)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	// Multicall response should be equal to the expected one
	s.Require().Equal(response.Result, "0x000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000034bb966cbf882cd7c000000")

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)
	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
	defer func() { _ = resp.Body.Close() }()
}

func (s *ExchangesE2ESuite) TestGemini() {
	ex := origin.NewExchange("gemini").
		WithSymbol("ETH/BTC").
		WithPrice(1).
		WithBid(2).
		WithAsk(3).
		WithVolume(4)

	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/v1/pubticker/ethbtc", s.url)
	resp, err := http.Get(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var response struct {
		Last string
		Bid  string
		Ask  string
	}

	err = parseBody(resp, &response)

	s.Require().NoError(err)
	s.Require().Equal("1.000000", response.Last)
	s.Require().Equal("2.000000", response.Bid)
	s.Require().Equal("3.000000", response.Ask)

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)

	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *ExchangesE2ESuite) TestHitBTC() {
	ex := origin.NewExchange("hitbtc").
		WithSymbol("ETH/BTC").
		WithPrice(1).
		WithBid(2).
		WithAsk(3).
		WithVolume(4)

	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/api/2/public/ticker/ETHBTC", s.url)
	resp, err := http.Get(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var response struct {
		Symbol string
		Last   string
		Bid    string
		Ask    string
		Volume string
	}

	err = parseBody(resp, &response)

	s.Require().NoError(err)
	s.Require().Equal("ETHBTC", response.Symbol)
	s.Require().Equal("1.000000", response.Last)
	s.Require().Equal("2.000000", response.Bid)
	s.Require().Equal("3.000000", response.Ask)
	s.Require().Equal("4.000000", response.Volume)

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)

	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *ExchangesE2ESuite) TestHuobi() {
	ts := time.Now()
	ex := origin.NewExchange("huobi").
		WithSymbol("ETH/BTC").
		WithPrice(1).
		WithBid(2).
		WithAsk(3).
		WithVolume(4).
		WithTime(ts)

	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/market/detail/merged?symbol=ethbtc", s.url)
	resp, err := http.Get(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var response struct {
		TS   int64 `json:"ts"`
		Tick struct {
			Vol float64
			Bid []float64
			Ask []float64
		}
	}

	err = parseBody(resp, &response)

	s.Require().NoError(err)
	s.Require().Equal(ts.UnixMilli(), response.TS)
	s.Require().Equal(float64(2), response.Tick.Bid[0])
	s.Require().Equal(float64(3), response.Tick.Ask[0])
	s.Require().Equal(float64(4), response.Tick.Vol)

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)

	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *ExchangesE2ESuite) TestKraken() {
	ex := origin.NewExchange("kraken").
		WithSymbol("ETH/BTC").
		WithPrice(1).
		WithBid(2).
		WithAsk(3).
		WithVolume(4)

	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/0/public/Ticker?pair=ETHBTC", s.url)
	resp, err := http.Get(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var response struct {
		Result map[string]struct {
			A []string
			B []string
			C []string
			V []string
		}
	}

	err = parseBody(resp, &response)

	s.Require().NoError(err)
	s.Require().Len(response.Result, 1)
	s.Require().Equal("1.000000", response.Result["ETHBTC"].C[0])
	s.Require().Equal("2.000000", response.Result["ETHBTC"].B[0])
	s.Require().Equal("3.000000", response.Result["ETHBTC"].A[0])
	s.Require().Equal("4.000000", response.Result["ETHBTC"].V[0])

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)

	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *ExchangesE2ESuite) TestKuCoin() {
	ts := time.Now()
	ex := origin.NewExchange("kucoin").
		WithSymbol("ETH/BTC").
		WithPrice(1).
		WithBid(2).
		WithAsk(3).
		WithTime(ts)

	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/api/v1/market/orderbook/level1?symbol=ETH-BTC", s.url)
	resp, err := http.Get(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var response struct {
		Data struct {
			Time  int64
			Price string
			Ask   string `json:"bestAsk"`
			Bid   string `json:"bestBid"`
		}
	}

	err = parseBody(resp, &response)

	s.Require().NoError(err)
	s.Require().Equal("1.000000", response.Data.Price)
	s.Require().Equal("2.000000", response.Data.Bid)
	s.Require().Equal("3.000000", response.Data.Ask)
	s.Require().Equal(ts.UnixMilli(), response.Data.Time)

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)

	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *ExchangesE2ESuite) TestOkex() {
	ts := time.Now()
	ex := origin.NewExchange("okex").
		WithSymbol("ETH/BTC").
		WithPrice(1).
		WithBid(2).
		WithAsk(3).
		WithVolume(4).
		WithTime(ts)

	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/api/spot/v3/instruments/ETH-BTC/ticker", s.url)
	resp, err := http.Get(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var response struct {
		Last      string
		Ask       string
		Bid       string
		Volume    string `json:"base_volume_24h"`
		Timestamp string
	}

	err = parseBody(resp, &response)

	s.Require().NoError(err)
	s.Require().Equal("1.000000", response.Last)
	s.Require().Equal("2.000000", response.Bid)
	s.Require().Equal("3.000000", response.Ask)
	s.Require().Equal("4.000000", response.Volume)
	s.Require().Equal(ts.String(), response.Timestamp)

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)

	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *ExchangesE2ESuite) TestRocketPool() {
	ex := origin.NewExchange("rocketpool").WithSymbol("RETH/ETH").WithPrice(1)
	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().Error(err)

	const blockNumber int = 100
	rate := big.NewInt(0.94 * 1e18)
	ex = ex.
		WithCustom("RETH/ETH", types.MustAddressFromHex("0xae78736Cd615f374D3085123A210448E74Fc6393")).
		WithFunctionData("getExchangeRate", []origin.FunctionData{
			{
				Args:   []any{},
				Return: []any{rate},
			},
		}).
		WithCustom("blockNumber", blockNumber).
		WithPrice(0.94)
	err = infestor.NewMocksBuilder().Debug().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/", s.url)

	// Calling RPC url via multicall contract, coins()
	reqJSON := fmt.Sprintf(origin.RPCCallRequestJSON, "null", "0xeefba1e63905ef1d7acba5a8513c70307c1ce441", "0x252dba42000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000020000000000000000000000000ae78736cd615f374d3085123a210448e74fc639300000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000004e6aa216c00000000000000000000000000000000000000000000000000000000", "latest")
	jsonStr := []byte(fmt.Sprint(reqJSON))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var response jsonrpcMessage
	err = parseBody(resp, &response)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	// Multicall response should be equal to the expected one
	s.Require().Equal(response.Result, "0x000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000d0b8d0508de0000")

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)
	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
	defer func() { _ = resp.Body.Close() }()
}

func (s *ExchangesE2ESuite) TestSDAI() {
	ex := origin.NewExchange("sdai").WithSymbol("SDAI/DAI").WithPrice(1)
	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().Error(err)

	const blockNumber int = 100
	rate := big.NewInt(1.02 * 1e18)
	ex = ex.
		WithCustom("SDAI/DAI", types.MustAddressFromHex("0x83F20F44975D03b1b09e64809B757c47f942BEeA")).
		WithFunctionData("previewRedeem", []origin.FunctionData{
			{
				Args:   []any{big.NewInt(1e18)},
				Return: []any{rate},
			},
		}).
		WithCustom("blockNumber", blockNumber).
		WithPrice(1.02)
	err = infestor.NewMocksBuilder().Debug().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/", s.url)

	// Calling RPC url via multicall contract, coins()
	reqJSON := fmt.Sprintf(origin.RPCCallRequestJSON, "null", "0xeefba1e63905ef1d7acba5a8513c70307c1ce441", "0x252dba4200000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000083f20f44975d03b1b09e64809b757c47f942beea000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000244cdad5060000000000000000000000000000000000000000000000000de0b6b3a764000000000000000000000000000000000000000000000000000000000000", "latest")
	jsonStr := []byte(fmt.Sprint(reqJSON))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var response jsonrpcMessage
	err = parseBody(resp, &response)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	// Multicall response should be equal to the expected one
	s.Require().Equal(response.Result, "0x000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000e27c49886e60000")

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)
	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
	defer func() { _ = resp.Body.Close() }()
}

func (s *ExchangesE2ESuite) TestSushiswap() {
	ex := origin.NewExchange("sushiswap").WithSymbol("DAI/WETH").WithPrice(1)
	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().Error(err)

	const blockNumber int = 100
	ex = ex.
		WithCustom("DAI/WETH", types.MustAddressFromHex("0xC3D03e4F041Fd4cD388c549Ee2A29a9E5075882f")).
		WithFunctionData("getReserves", []origin.FunctionData{
			{
				Args: []any{},
				Return: []any{
					new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18)),
					new(big.Int).Mul(big.NewInt(200), big.NewInt(1e18)),
					big.NewInt(1692188531),
				},
			},
		}).
		WithFunctionData("token0", []origin.FunctionData{
			{
				Args:   []any{},
				Return: []any{types.MustAddressFromHex("0x6B175474E89094C44Da98b954EedeAC495271d0F")},
			},
		}).
		WithFunctionData("token1", []origin.FunctionData{
			{
				Args:   []any{},
				Return: []any{types.MustAddressFromHex("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")},
			},
		}).
		WithCustom("tokens", []types.Address{
			types.MustAddressFromHex("0x6B175474E89094C44Da98b954EedeAC495271d0F"),
			types.MustAddressFromHex("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
		}).
		WithFunctionData("symbols", []origin.FunctionData{
			{
				Args:   []any{},
				Return: []any{"DAI"},
			},
			{
				Args:   []any{},
				Return: []any{"WETH"},
			},
		}).
		WithFunctionData("decimals", []origin.FunctionData{
			{
				Args:   []any{},
				Return: []any{big.NewInt(18)},
			},
			{
				Args:   []any{},
				Return: []any{big.NewInt(18)},
			},
		}).
		WithCustom("blockNumber", blockNumber).
		WithPrice(2)
	err = infestor.NewMocksBuilder().Debug().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/", s.url)

	// Calling RPC url via multicall contract, getReserves()
	reqJSON := fmt.Sprintf(origin.RPCCallRequestJSON, "null", "0xeefba1e63905ef1d7acba5a8513c70307c1ce441", "0x0x252dba42000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000020000000000000000000000000c3d03e4f041fd4cd388c549ee2a29a9e5075882f000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000040902f1ac00000000000000000000000000000000000000000000000000000000", "latest")
	jsonStr := []byte(fmt.Sprint(reqJSON))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var response jsonrpcMessage
	err = parseBody(resp, &response)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	// Multicall response should be equal to the expected one
	s.Require().Equal(response.Result, "0x000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000056bc75e2d6310000000000000000000000000000000000000000000000000000ad78ebc5ac62000000000000000000000000000000000000000000000000000000000000064dcbf73")

	// token0(), token1()
	reqJSON = fmt.Sprintf(origin.RPCCallRequestJSON, "null", "0xeefba1e63905ef1d7acba5a8513c70307c1ce441", "0x252dba4200000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000c3d03e4f041fd4cd388c549ee2a29a9e5075882f000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000040dfe168100000000000000000000000000000000000000000000000000000000000000000000000000000000c3d03e4f041fd4cd388c549ee2a29a9e5075882f00000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000004d21220a700000000000000000000000000000000000000000000000000000000", "latest")
	jsonStr = []byte(fmt.Sprint(reqJSON))
	resp, err = http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	err = parseBody(resp, &response)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	s.Require().Equal(response.Result, "0x0000000000000000000000000000000000000000000000000000000000000064000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000200000000000000000000000006b175474e89094c44da98b954eedeac495271d0f0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")

	// symbol(), decimals()
	reqJSON = fmt.Sprintf(origin.RPCCallRequestJSON, "null", "0xeefba1e63905ef1d7acba5a8513c70307c1ce441", "0x252dba420000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000018000000000000000000000000000000000000000000000000000000000000002000000000000000000000000006b175474e89094c44da98b954eedeac495271d0f0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000495d89b41000000000000000000000000000000000000000000000000000000000000000000000000000000006b175474e89094c44da98b954eedeac495271d0f00000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000004313ce56700000000000000000000000000000000000000000000000000000000000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc20000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000495d89b4100000000000000000000000000000000000000000000000000000000000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000004313ce56700000000000000000000000000000000000000000000000000000000", "latest")
	jsonStr = []byte(fmt.Sprint(reqJSON))
	resp, err = http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	err = parseBody(resp, &response)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	s.Require().Equal(response.Result, "0x00000000000000000000000000000000000000000000000000000000000000640000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000014000000000000000000000000000000000000000000000000000000000000001c0000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000003444149000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000004574554480000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000012")

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)
	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
	defer func() { _ = resp.Body.Close() }()
}

func (s *ExchangesE2ESuite) TestUniswapV2() {
	ex := origin.NewExchange("uniswapV2").WithSymbol("STETH/WETH").WithPrice(1)
	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().Error(err)

	const blockNumber int = 100
	ex = ex.
		WithCustom("STETH/WETH", types.MustAddressFromHex("0x4028DAAC072e492d34a3Afdbef0ba7e35D8b55C4")).
		WithFunctionData("getReserves", []origin.FunctionData{
			{
				Args: []any{},
				Return: []any{
					new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18)),
					new(big.Int).Mul(big.NewInt(200), big.NewInt(1e18)),
					big.NewInt(1692188531),
				},
			},
		}).
		WithFunctionData("token0", []origin.FunctionData{
			{
				Args:   []any{},
				Return: []any{types.MustAddressFromHex("0xae7ab96520DE3A18E5e111B5EaAb095312D7fE84")},
			},
		}).
		WithFunctionData("token1", []origin.FunctionData{
			{
				Args:   []any{},
				Return: []any{types.MustAddressFromHex("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")},
			},
		}).
		WithCustom("tokens", []types.Address{
			types.MustAddressFromHex("0xae7ab96520DE3A18E5e111B5EaAb095312D7fE84"),
			types.MustAddressFromHex("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
		}).
		WithFunctionData("symbols", []origin.FunctionData{
			{
				Args:   []any{},
				Return: []any{"stETH"},
			},
			{
				Args:   []any{},
				Return: []any{"WETH"},
			},
		}).
		WithFunctionData("decimals", []origin.FunctionData{
			{
				Args:   []any{},
				Return: []any{big.NewInt(18)},
			},
			{
				Args:   []any{},
				Return: []any{big.NewInt(18)},
			},
		}).
		WithCustom("blockNumber", blockNumber).
		WithPrice(2)
	err = infestor.NewMocksBuilder().Debug().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/", s.url)

	// Calling RPC url via multicall contract, getReserves()
	reqJSON := fmt.Sprintf(origin.RPCCallRequestJSON, "null", "0xeefba1e63905ef1d7acba5a8513c70307c1ce441", "0x252dba420000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000200000000000000000000000004028daac072e492d34a3afdbef0ba7e35d8b55c4000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000040902f1ac00000000000000000000000000000000000000000000000000000000", "latest")
	jsonStr := []byte(fmt.Sprint(reqJSON))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var response jsonrpcMessage
	err = parseBody(resp, &response)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	// Multicall response should be equal to the expected one
	s.Require().Equal(response.Result, "0x000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000056bc75e2d6310000000000000000000000000000000000000000000000000000ad78ebc5ac62000000000000000000000000000000000000000000000000000000000000064dcbf73")

	// token0(), token1()
	reqJSON = fmt.Sprintf(origin.RPCCallRequestJSON, "null", "0xeefba1e63905ef1d7acba5a8513c70307c1ce441", "0x252dba4200000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000004028daac072e492d34a3afdbef0ba7e35d8b55c4000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000040dfe1681000000000000000000000000000000000000000000000000000000000000000000000000000000004028daac072e492d34a3afdbef0ba7e35d8b55c400000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000004d21220a700000000000000000000000000000000000000000000000000000000", "latest")
	jsonStr = []byte(fmt.Sprint(reqJSON))
	resp, err = http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	err = parseBody(resp, &response)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	s.Require().Equal(response.Result, "0x000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000020000000000000000000000000ae7ab96520de3a18e5e111b5eaab095312d7fe840000000000000000000000000000000000000000000000000000000000000020000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")

	// symbol(), decimals()
	reqJSON = fmt.Sprintf(origin.RPCCallRequestJSON, "null", "0xeefba1e63905ef1d7acba5a8513c70307c1ce441", "0x252dba42000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001800000000000000000000000000000000000000000000000000000000000000200000000000000000000000000ae7ab96520de3a18e5e111b5eaab095312d7fe840000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000495d89b4100000000000000000000000000000000000000000000000000000000000000000000000000000000ae7ab96520de3a18e5e111b5eaab095312d7fe8400000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000004313ce56700000000000000000000000000000000000000000000000000000000000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc20000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000495d89b4100000000000000000000000000000000000000000000000000000000000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000004313ce56700000000000000000000000000000000000000000000000000000000", "latest")
	jsonStr = []byte(fmt.Sprint(reqJSON))
	resp, err = http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	err = parseBody(resp, &response)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	s.Require().Equal(response.Result, "0x00000000000000000000000000000000000000000000000000000000000000640000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000014000000000000000000000000000000000000000000000000000000000000001c0000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000005737445544800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000004574554480000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000012")

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)
	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
	defer func() { _ = resp.Body.Close() }()
}

func (s *ExchangesE2ESuite) TestUniswapV3() {
	ex := origin.NewExchange("uniswapV3").WithSymbol("WSTETH/WETH").WithPrice(1)
	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().Error(err)

	const blockNumber int = 100
	sqrtPriceX96, _ := new(big.Int).SetString("84554395222218770838379633172", 10)
	ex = ex.
		WithCustom("WSTETH/WETH", types.MustAddressFromHex("0x109830a1AAaD605BbF02a9dFA7B0B92EC2FB7dAa")).
		WithFunctionData("slot0", []origin.FunctionData{
			{
				Args: []any{},
				Return: []any{
					sqrtPriceX96,     // sqrtPriceX96
					big.NewInt(1301), // tick
					big.NewInt(23),   // observationIndex
					big.NewInt(150),  // observationCardinality
					big.NewInt(150),  // observationCardinalityNext
					0,                // feeProtocol
					false,            // unlocked
				},
			},
		}).
		WithFunctionData("token0", []origin.FunctionData{
			{
				Args:   []any{},
				Return: []any{types.MustAddressFromHex("0x7f39C581F595B53c5cb19bD0b3f8dA6c935E2Ca0")},
			},
		}).
		WithFunctionData("token1", []origin.FunctionData{
			{
				Args:   []any{},
				Return: []any{types.MustAddressFromHex("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")},
			},
		}).
		WithCustom("tokens", []types.Address{
			types.MustAddressFromHex("0x7f39C581F595B53c5cb19bD0b3f8dA6c935E2Ca0"),
			types.MustAddressFromHex("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
		}).
		WithFunctionData("symbols", []origin.FunctionData{
			{
				Args:   []any{},
				Return: []any{"wstETH"},
			},
			{
				Args:   []any{},
				Return: []any{"WETH"},
			},
		}).
		WithFunctionData("decimals", []origin.FunctionData{
			{
				Args:   []any{},
				Return: []any{big.NewInt(18)},
			},
			{
				Args:   []any{},
				Return: []any{big.NewInt(18)},
			},
		}).
		WithCustom("blockNumber", blockNumber).
		WithPrice(2)
	err = infestor.NewMocksBuilder().Debug().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/", s.url)

	// Calling RPC url via multicall contract, slot0()
	reqJSON := fmt.Sprintf(origin.RPCCallRequestJSON, "null", "0xeefba1e63905ef1d7acba5a8513c70307c1ce441", "0x252dba42000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000020000000000000000000000000109830a1aaad605bbf02a9dfa7b0b92ec2fb7daa000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000043850c7bd00000000000000000000000000000000000000000000000000000000", "latest")
	jsonStr := []byte(fmt.Sprint(reqJSON))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var response jsonrpcMessage
	err = parseBody(resp, &response)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	// Multicall response should be equal to the expected one
	s.Require().Equal(response.Result, "0x000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000011135c1a5a808593f506d1a14000000000000000000000000000000000000000000000000000000000000051500000000000000000000000000000000000000000000000000000000000000170000000000000000000000000000000000000000000000000000000000000096000000000000000000000000000000000000000000000000000000000000009600000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

	// token0(), token1()
	reqJSON = fmt.Sprintf(origin.RPCCallRequestJSON, "null", "0xeefba1e63905ef1d7acba5a8513c70307c1ce441", "0x252dba4200000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000109830a1aaad605bbf02a9dfa7b0b92ec2fb7daa000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000040dfe168100000000000000000000000000000000000000000000000000000000000000000000000000000000109830a1aaad605bbf02a9dfa7b0b92ec2fb7daa00000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000004d21220a700000000000000000000000000000000000000000000000000000000", "latest")
	jsonStr = []byte(fmt.Sprint(reqJSON))
	resp, err = http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	err = parseBody(resp, &response)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	s.Require().Equal(response.Result, "0x0000000000000000000000000000000000000000000000000000000000000064000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000200000000000000000000000007f39c581f595b53c5cb19bd0b3f8da6c935e2ca00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")

	// symbol(), decimals()
	reqJSON = fmt.Sprintf(origin.RPCCallRequestJSON, "null", "0xeefba1e63905ef1d7acba5a8513c70307c1ce441", "0x252dba420000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000018000000000000000000000000000000000000000000000000000000000000002000000000000000000000000007f39c581f595b53c5cb19bd0b3f8da6c935e2ca00000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000495d89b41000000000000000000000000000000000000000000000000000000000000000000000000000000007f39c581f595b53c5cb19bd0b3f8da6c935e2ca000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000004313ce56700000000000000000000000000000000000000000000000000000000000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc20000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000495d89b4100000000000000000000000000000000000000000000000000000000000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000004313ce56700000000000000000000000000000000000000000000000000000000", "latest")
	jsonStr = []byte(fmt.Sprint(reqJSON))
	resp, err = http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	err = parseBody(resp, &response)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	s.Require().Equal(response.Result, "0x00000000000000000000000000000000000000000000000000000000000000640000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000014000000000000000000000000000000000000000000000000000000000000001c0000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000006777374455448000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000004574554480000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000012")

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)
	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
	defer func() { _ = resp.Body.Close() }()
}

func (s *ExchangesE2ESuite) TestUpbit() {
	ts := time.Now()
	ex := origin.NewExchange("upbit").
		WithSymbol("ETH/BTC").
		WithPrice(1).
		WithVolume(4).
		WithTime(ts)

	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/v1/ticker?markets=BTC-ETH", s.url)
	resp, err := http.Get(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var response []struct {
		Market    string
		Price     float64 `json:"trade_price"`
		Volume    float64 `json:"trade_volume"`
		Timestamp int64
	}

	err = parseBody(resp, &response)

	s.Require().NoError(err)
	s.Require().Len(response, 1)
	s.Require().Equal(float64(1), response[0].Price)
	s.Require().Equal(float64(4), response[0].Volume)
	s.Require().Equal(ts.UnixMilli(), response[0].Timestamp)

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)

	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *ExchangesE2ESuite) TestWSTETH() {
	ex := origin.NewExchange("wsteth").WithSymbol("WSTETH/STETH").WithPrice(1)
	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().Error(err)

	const blockNumber int = 100
	rate := big.NewInt(0.94 * 1e18)
	ex = ex.
		WithCustom("WSTETH/STETH", types.MustAddressFromHex("0x7f39C581F595B53c5cb19bD0b3f8dA6c935E2Ca0")).
		WithFunctionData("stEthPerToken", []origin.FunctionData{
			{
				Args:   []any{},
				Return: []any{rate},
			},
		}).
		WithCustom("blockNumber", blockNumber).
		WithPrice(0.94)
	err = infestor.NewMocksBuilder().Debug().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/", s.url)

	// Calling RPC url via multicall contract, coins()
	reqJSON := fmt.Sprintf(origin.RPCCallRequestJSON, "null", "0xeefba1e63905ef1d7acba5a8513c70307c1ce441", "0x252dba420000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000200000000000000000000000007f39c581f595b53c5cb19bd0b3f8da6c935e2ca000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000004035faf8200000000000000000000000000000000000000000000000000000000", "latest")
	jsonStr := []byte(fmt.Sprint(reqJSON))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var response jsonrpcMessage
	err = parseBody(resp, &response)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	// Multicall response should be equal to the expected one
	s.Require().Equal(response.Result, "0x000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000d0b8d0508de0000")

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)
	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
	defer func() { _ = resp.Body.Close() }()
}
