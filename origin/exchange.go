package origin

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Mockable interface for exchange implementation.
type Mockable interface {
	// BuildMock builds yaml specification for exchange specific mock
	BuildMock(e []ExchangeMock) ([]byte, error)
}

type MockableFunc func(e ExchangeMock) ([]byte, error)

// CombineMocks is helper function that helps most of exchanges to build mocks.
func CombineMocks(e []ExchangeMock, f MockableFunc) ([]byte, error) {
	var result bytes.Buffer
	for _, mock := range e {
		m, err := f(mock)
		if err != nil {
			return nil, err
		}
		result.Write(m)
	}
	return result.Bytes(), nil
}

var exchanges = map[string]Mockable{
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

func BuildMock(exchangeName string, e []ExchangeMock) ([]byte, error) {
	ex, ok := exchanges[exchangeName]
	if !ok {
		return nil, fmt.Errorf("failed to find exchange name %s", exchangeName)
	}
	return ex.BuildMock(e)
}
