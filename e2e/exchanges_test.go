package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/defiweb/go-eth/types"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/chronicleprotocol/infestor"
	"github.com/chronicleprotocol/infestor/origin"
	"github.com/chronicleprotocol/infestor/smocker"
)

func parseBody(resp *http.Response, r interface{}) error {
	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
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
	reqJSON := fmt.Sprintf(origin.RpcCallRequestJSON, "null", "0xeefba1e63905ef1d7acba5a8513c70307c1ce441", "0x252dba42000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001800000000000000000000000000000000000000000000000000000000000000200000000000000000000000000ae7ab96520de3a18e5e111b5eaab095312d7fe840000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000495d89b4100000000000000000000000000000000000000000000000000000000000000000000000000000000ae7ab96520de3a18e5e111b5eaab095312d7fe8400000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000004313ce56700000000000000000000000000000000000000000000000000000000000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc20000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000495d89b4100000000000000000000000000000000000000000000000000000000000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000004313ce56700000000000000000000000000000000000000000000000000000000", "latest")
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
}

func (s *ExchangesE2ESuite) TestBalancerV2() {
	ex := origin.NewExchange("balancerV2").WithSymbol("STETH/ETH").WithPrice(1)
	// No `pool` field should fail
	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().Error(err)

	const blockNumber int = 100
	//price := types.Bytes(big.NewInt(0.94 * 1e18).Bytes()).PadLeft(32)
	price := big.NewInt(0.94 * 1e18)
	ex = ex.
		WithCustom("pool", types.MustAddressFromHex("0x32296969ef14eb0c6d29669c550d4a0449130230")).
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
	reqJSON := fmt.Sprintf(origin.RpcCallRequestJSON, "null", "0xeefba1e63905ef1d7acba5a8513c70307c1ce441", "0x252dba4200000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000032296969ef14eb0c6d29669c550d4a044913023000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000024b10be739000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", "latest")
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
	s.Require().Equal(ts.Format(format), response.Time)

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)

	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

// todo, Curve
func (s *ExchangesE2ESuite) TestCurve() {

}

// todo, DSR

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

// todo, rocketpool
// todo, sdai
// todo, sushiswap
// todo, uniswapV2

func (s *ExchangesE2ESuite) TestUniswapV3() {
	ex := origin.NewExchange("uniswap_v3").WithPrice(1)

	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/subgraphs/name/uniswap/uniswap-v3", s.url)
	jsonStr := []byte(`{"match":"0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640"}`)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	defer func() { _ = resp.Body.Close() }()
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
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

// todo, wsteth
