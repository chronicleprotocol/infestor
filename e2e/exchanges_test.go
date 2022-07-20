package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/chronicleprotocol/infestor/origin"

	"github.com/chronicleprotocol/infestor"

	"github.com/stretchr/testify/suite"

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

func (s *ExchangesE2ESuite) TestBalancer() {
	ex := origin.NewExchange("balancer").WithSymbol("BAL/USD").WithPrice(1)
	// No `contract` field should fail
	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().Error(err)

	// With contract
	contract := "0xba100000625a3754423978a60c9317c58a424e3d"
	ex = ex.WithCustom("contract", contract)
	err = infestor.NewMocksBuilder().Debug().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/subgraphs/name/balancer-labs/balancer", s.url)
	reqJSON := `{"query":"query($id:String) { tokenPrices(where: { id: $id }) { symbol price poolLiquidity } }", "variables": { "id": "%s" } }`
	jsonStr := []byte(fmt.Sprintf(reqJSON, contract))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	response := struct {
		Data struct {
			TokenPrices []struct {
				Price  string `json:"price"`
				Symbol string `json:"symbol"`
			} `json:"tokenPrices"`
		}
		Price string
	}{}
	err = parseBody(resp, &response)

	s.Require().NoError(err)
	s.Require().NotNil(response)
	s.Require().Len(response.Data.TokenPrices, 1)
	s.Require().Equal("1.00000000", response.Data.TokenPrices[0].Price)
	s.Require().Equal("BAL", response.Data.TokenPrices[0].Symbol)

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)
	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
	defer func() { _ = resp.Body.Close() }()
}

func (s *ExchangesE2ESuite) TestBalancerV2_getLatest() {
	ex := origin.NewExchange("balancerV2").WithSymbol("STETH/ETH").WithPrice(1)
	// No `pool` field should fail
	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().Error(err)

	// With set rate and price
	price := "0x0000000000000000000000000000000000000000000000000dd22d6848e229b8"
	ex = ex.WithCustom("rate", price).WithCustom("price", price)
	err = infestor.NewMocksBuilder().Debug().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/", s.url)
	reqJSON := `{"method":"eth_call","params":[{"from":null,"to":"0x0000000000000000000000000000000000000000","data":"0xb10be7390000000000000000000000000000000000000000000000000000000000000000"}, "latest"],"id":1,"jsonrpc":"2.0"}`
	jsonStr := []byte(fmt.Sprintf(reqJSON))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var response jsonrpcMessage
	err = parseBody(resp, &response)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	s.Require().Equal(price, response.Result)

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)
	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
	defer func() { _ = resp.Body.Close() }()
}

func (s *ExchangesE2ESuite) TestBalancerV2_getPriceRateCache() {
	ex := origin.NewExchange("balancerV2").WithSymbol("STETH/ETH").WithPrice(1)
	// No `pool` field should fail
	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().Error(err)

	// With set rate and price
	rate := "0x0000000000000000000000000000000000000000000000000dd22d6848e229b8"
	ex = ex.WithCustom("rate", rate).WithCustom("price", rate)
	err = infestor.NewMocksBuilder().Debug().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/", s.url)
	reqJSON := `{"method":"eth_call","params":[{"from":null,"to":"0x0000000000000000000000000000000000000000","data":"0xb867ee5a"}, "latest"],"id":1,"jsonrpc":"2.0"}`
	jsonStr := []byte(fmt.Sprintf(reqJSON))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var response jsonrpcMessage
	err = parseBody(resp, &response)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	fmt.Println(response.Result)
	s.Require().Equal(rate, response.Result)

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

func (s *ExchangesE2ESuite) TestBithumb() {
	ts := time.Now()
	ex := origin.NewExchange("bitthumb").
		WithSymbol("ETH/BTC").
		WithPrice(1).
		WithVolume(2).
		WithTime(ts)

	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/openapi/v1/spot/ticker?symbol=ETH-BTC", s.url)
	resp, err := http.Get(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	response := struct {
		Data []struct {
			Vol string
			C   string
			S   string
		}
		Timestamp int64
	}{}

	err = parseBody(resp, &response)

	s.Require().NoError(err)
	s.Require().Len(response.Data, 1)
	s.Require().Equal("ETH-BTC", response.Data[0].S)
	s.Require().Equal("1.000000", response.Data[0].C)
	s.Require().Equal("2.000000", response.Data[0].Vol)
	s.Require().Equal(ts.UnixMilli(), response.Timestamp)

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

func (s *ExchangesE2ESuite) TestBitTrex() {
	ex := origin.NewExchange("bittrex").
		WithSymbol("ETH/BTC").
		WithPrice(1).
		WithBid(2).
		WithAsk(3)

	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/api/v1.1/public/getticker?market=BTC-ETH", s.url)
	resp, err := http.Get(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	response := struct {
		Result struct {
			Last float64
			Bid  float64
			Ask  float64
		}
	}{}

	err = parseBody(resp, &response)

	s.Require().NoError(err)
	s.Require().Equal(float64(1), response.Result.Last)
	s.Require().Equal(float64(2), response.Result.Bid)
	s.Require().Equal(float64(3), response.Result.Ask)

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

func (s *ExchangesE2ESuite) TestCryptoCompare() {
	ex := origin.NewExchange("cryptocompare").
		WithSymbol("ETH/BTC").
		WithPrice(1)

	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/data/price?fsym=ETH&tsyms=BTC", s.url)
	resp, err := http.Get(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	response := map[string]float64{}

	err = parseBody(resp, &response)

	s.Require().NoError(err)
	s.Require().Len(response, 1)
	s.Require().Equal(float64(1), response["BTC"])

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)

	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *ExchangesE2ESuite) TestFtx() {
	ex := origin.NewExchange("ftx").
		WithSymbol("ETH/BTC").
		WithPrice(1).
		WithBid(2).
		WithAsk(3).
		WithVolume(4)

	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/api/markets/ETH/BTC", s.url)
	resp, err := http.Get(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	response := struct {
		Result struct {
			Name          string
			Last          float64
			Bid           float64
			Ask           float64
			Price         float64
			BaseCurrency  string
			QuoteCurrency string
			Volume        float64 `json:"volumeUsd24h"`
		}
	}{}

	err = parseBody(resp, &response)

	s.Require().NoError(err)
	s.Require().Equal(float64(1), response.Result.Last)
	s.Require().Equal(float64(1), response.Result.Price)
	s.Require().Equal(float64(2), response.Result.Bid)
	s.Require().Equal(float64(3), response.Result.Ask)
	s.Require().Equal(float64(4), response.Result.Volume)
	s.Require().Equal("ETH", response.Result.BaseCurrency)
	s.Require().Equal("BTC", response.Result.QuoteCurrency)

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)

	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *ExchangesE2ESuite) TestGateIO() {
	ex := origin.NewExchange("gateio").
		WithSymbol("ETH/BTC").
		WithPrice(1).
		WithBid(2).
		WithAsk(3).
		WithVolume(4)

	err := infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	url := fmt.Sprintf("%s/api/v4/spot/tickers?currency_pair=ETH_BTC", s.url)
	resp, err := http.Get(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var response []struct {
		Pair   string `json:"currency_pair"`
		Last   string
		Bid    string `json:"highest_bid"`
		Ask    string `json:"lowest_ask"`
		Volume string `json:"base_volume"`
	}

	err = parseBody(resp, &response)

	s.Require().NoError(err)
	s.Require().Len(response, 1)
	s.Require().Equal("1.000000", response[0].Last)
	s.Require().Equal("2.000000", response[0].Bid)
	s.Require().Equal("3.000000", response[0].Ask)
	s.Require().Equal("4.000000", response[0].Volume)
	s.Require().Equal("ETH_BTC", response[0].Pair)

	// Test status code
	ex = ex.WithStatusCode(http.StatusNotFound)

	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
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

func (s *ExchangesE2ESuite) TestKyber() {
	ts := time.Now()

	err := infestor.NewMocksBuilder().
		Debug().
		Reset().
		Add(origin.NewExchange("kyber").WithSymbol("ETH/BTC").WithPrice(1).WithTime(ts)).
		Add(origin.NewExchange("kyber").WithSymbol("MKR/BTC").WithPrice(2).WithTime(ts)).
		Deploy(s.api)

	s.Require().NoError(err)

	url := fmt.Sprintf("%s/change24h", s.url)
	resp, err := http.Get(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	response := map[string]struct {
		Timestamp   int64   `json:"timestamp"`
		TokenSymbol string  `json:"token_symbol"`
		Price       float64 `json:"rate_eth_now"`
	}{}
	err = parseBody(resp, &response)

	s.Require().NoError(err)
	s.Require().NotNil(response["BTC_ETH"])
	s.Require().Equal(ts.UnixMilli(), response["BTC_ETH"].Timestamp)
	s.Require().Equal("ETH", response["BTC_ETH"].TokenSymbol)
	s.Require().Equal(float64(1), response["BTC_ETH"].Price)
	s.Require().NotNil(response["BTC_MKR"])
	s.Require().Equal(ts.UnixMilli(), response["BTC_MKR"].Timestamp)
	s.Require().Equal("MKR", response["BTC_MKR"].TokenSymbol)
	s.Require().Equal(float64(2), response["BTC_MKR"].Price)

	// Test status code
	err = infestor.NewMocksBuilder().
		Debug().
		Reset().
		Add(origin.NewExchange("kyber").WithSymbol("ETH/BTC").WithStatusCode(http.StatusNotFound)).
		Deploy(s.api)

	s.Require().NoError(err)

	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *ExchangesE2ESuite) TestPoloniex() {
	err := infestor.NewMocksBuilder().
		Reset().
		Add(origin.NewExchange("poloniex").WithSymbol("ETH/BTC").WithPrice(1).WithAsk(2).WithBid(3).WithVolume(4)).
		Add(origin.NewExchange("poloniex").WithSymbol("MKR/BTC").WithPrice(5).WithAsk(6).WithBid(7).WithVolume(8)).
		Deploy(s.api)

	s.Require().NoError(err)

	url := fmt.Sprintf("%s/public?command=returnTicker", s.url)
	resp, err := http.Get(url)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	response := map[string]struct {
		Last       string
		LowestAsk  string `json:"lowestAsk"`
		HighestBid string `json:"highestBid"`
		BaseVolume string `json:"baseVolume"`
	}{}
	err = parseBody(resp, &response)

	s.Require().NoError(err)
	s.Require().NotNil(response["BTC_ETH"])
	s.Require().Equal("1.00000000", response["BTC_ETH"].Last)
	s.Require().Equal("2.00000000", response["BTC_ETH"].LowestAsk)
	s.Require().Equal("3.00000000", response["BTC_ETH"].HighestBid)
	s.Require().Equal("4.00000000", response["BTC_ETH"].BaseVolume)
	s.Require().NotNil(response["BTC_MKR"])
	s.Require().Equal("5.00000000", response["BTC_MKR"].Last)
	s.Require().Equal("6.00000000", response["BTC_MKR"].LowestAsk)
	s.Require().Equal("7.00000000", response["BTC_MKR"].HighestBid)
	s.Require().Equal("8.00000000", response["BTC_MKR"].BaseVolume)

	// Test status code
	err = infestor.NewMocksBuilder().
		Reset().
		Add(origin.NewExchange("poloniex").WithSymbol("ETH/BTC").WithStatusCode(http.StatusConflict)).
		Deploy(s.api)
	s.Require().NoError(err)

	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Require().Equal(http.StatusConflict, resp.StatusCode)
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
