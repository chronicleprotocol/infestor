package origin

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Mockable interface for exchange implementation.
type Mockable interface {
	// BuildMock builds yaml specification for exchange specific mock
	BuildMock(e Exchange) ([]byte, error)
}

var exchanges = map[string]Mockable{
	"balancer":      Balancer{},
	"binance":       Binance{},
	"bitfinex":      Bitfinex{},
	"butthumb":      Bithumb{},
	"buthumb":       Bithumb{},
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

type Exchange struct {
	Name       string
	StatusCode int
	Symbol     Symbol
	Price      float64
	Volume     float64
	Ask        float64
	Bid        float64
	Timestamp  time.Time
}

func NewExchange(name string) *Exchange {
	return &Exchange{StatusCode: http.StatusOK, Name: name}
}

func (e *Exchange) WithStatusCode(statusCode int) *Exchange {
	e.StatusCode = statusCode
	return e
}

func (e *Exchange) WithSymbol(symbol string) *Exchange {
	e.Symbol = NewSymbol(symbol)
	return e
}

func (e *Exchange) WithPrice(price float64) *Exchange {
	e.Price = price
	return e
}

func (e *Exchange) WithVolume(volume float64) *Exchange {
	e.Volume = volume
	return e
}

func (e *Exchange) WithAsk(ask float64) *Exchange {
	e.Ask = ask
	return e
}

func (e *Exchange) WithBid(bid float64) *Exchange {
	e.Bid = bid
	return e
}

func (e *Exchange) WithTime(timestamp time.Time) *Exchange {
	e.Timestamp = timestamp
	return e
}

func BuildMock(e Exchange) ([]byte, error) {
	ex, ok := exchanges[e.Name]
	if !ok {
		return nil, fmt.Errorf("failed to find exchange name %s", e.Name)
	}
	return ex.BuildMock(e)
}
