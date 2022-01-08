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

	"github.com/chronicleprotocol/infestor/smocker"
	"github.com/stretchr/testify/suite"
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

func (s *ExchangesE2ESuite) SetupSuite() {
	smockerHost, exist := os.LookupEnv("SMOCKER_HOST")
	s.Require().True(exist, "SMOCKER_HOST env variable have to be set")

	s.api = smocker.API{
		Host: smockerHost,
		Port: 8081,
	}

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
	err = infestor.NewMocksBuilder().Reset().Add(ex).Deploy(s.api)
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
