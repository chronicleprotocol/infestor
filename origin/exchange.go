package origin

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chronicleprotocol/infestor/smocker"
)

// Mockable interface for exchange implementation.
type Mockable interface {
	// BuildMocks builds yaml specification for exchange specific mock
	BuildMocks(e []ExchangeMock) ([]*smocker.Mock, error)
}

type MockableFunc func(e ExchangeMock) (*smocker.Mock, error)

// CombineMocks is helper function that helps most of exchanges to build mocks.
func CombineMocks(e []ExchangeMock, f MockableFunc) ([]*smocker.Mock, error) {
	var mocks []*smocker.Mock
	for _, ex := range e {
		m, err := f(ex)
		if err != nil {
			return nil, fmt.Errorf("failed to build mock: %w", err)
		}
		mocks = append(mocks, m)
	}
	return mocks, nil
}

var exchanges = map[string]Mockable{
	"rocketpool":    RocketPool{},
	"balancer":      Balancer{},
	"binance":       Binance{},
	"bitfinex":      Bitfinex{},
	"bitthumb":      Bithumb{},
	"bithumb":       Bithumb{},
	"bitstamp":      BitStamp{},
	"bittrex":       BitTrex{},
	"coinbase":      Coinbase{},
	"cryptocompare": CryptoCompare{},
	"ftx":           Ftx{},
	"gateio":        GateIO{},
	"gemini":        Gemini{},
	"hitbtc":        HitBTC{},
	"huobi":         Huobi{},
	"kraken":        Kraken{},
	"kucoin":        KuCoin{},
	"kyber":         Kyber{},
	"okex":          Okex{},
	"poloniex":      Poloniex{},
	"upbit":         Upbit{},
	"wsteth":        WSTETH{},
	"balancerV2":    BalancerV2{},
	"curve":         Curve{},
	"ethrpc":        EthRPC{},
	"uniswap_v3":    UniswapV3{},
}

// Symbol represents an asset pair.
type Symbol struct {
	Base  string
	Quote string
}

// NewSymbol returns a new Pair for given string. The string must be formatted
// as "BASE/QUOTE".
func NewSymbol(s string) Symbol {
	ss := strings.Split(s, "/")
	if len(ss) != 2 {
		panic("Invalid symbol !")
	}
	return Symbol{Base: strings.ToUpper(ss[0]), Quote: strings.ToUpper(ss[1])}
}

func (p Symbol) String() string {
	return fmt.Sprintf("%s/%s", p.Base, p.Quote)
}

func (p Symbol) Format(format string) string {
	return fmt.Sprintf(format, p.Base, p.Quote)
}

type ExchangeMock struct {
	Name       string
	StatusCode int
	Symbol     Symbol
	Price      float64
	Volume     float64
	Ask        float64
	Bid        float64
	Timestamp  time.Time
	Custom     map[string]string
}

func NewExchange(name string) *ExchangeMock {
	return &ExchangeMock{
		StatusCode: http.StatusOK,
		Name:       name,
		Timestamp:  time.Now(),
		Custom:     make(map[string]string),
	}
}

func (e *ExchangeMock) WithStatusCode(statusCode int) *ExchangeMock {
	e.StatusCode = statusCode
	return e
}

func (e *ExchangeMock) WithSymbol(symbol string) *ExchangeMock {
	e.Symbol = NewSymbol(symbol)
	return e
}

func (e *ExchangeMock) WithPrice(price float64) *ExchangeMock {
	e.Price = price
	return e
}

func (e *ExchangeMock) WithVolume(volume float64) *ExchangeMock {
	e.Volume = volume
	return e
}

func (e *ExchangeMock) WithAsk(ask float64) *ExchangeMock {
	e.Ask = ask
	return e
}

func (e *ExchangeMock) WithBid(bid float64) *ExchangeMock {
	e.Bid = bid
	return e
}

func (e *ExchangeMock) WithTime(timestamp time.Time) *ExchangeMock {
	e.Timestamp = timestamp
	return e
}

func (e *ExchangeMock) WithCustom(key, value string) *ExchangeMock {
	e.Custom[key] = value
	return e
}

func BuildMocksForExchanges(exchangeName string, e []ExchangeMock) ([]*smocker.Mock, error) {
	ex, ok := exchanges[exchangeName]
	if !ok {
		return nil, fmt.Errorf("failed to find exchange name %s", exchangeName)
	}
	mocks, err := ex.BuildMocks(e)
	if err != nil {
		return nil, err
	}
	for _, m := range mocks {
		mErr := m.Validate()
		if mErr != nil {
			return nil, mErr
		}
	}
	return mocks, nil
}
