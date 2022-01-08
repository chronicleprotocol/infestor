package example

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/chronicleprotocol/infestor/origin"

	"github.com/chronicleprotocol/infestor"
	"github.com/chronicleprotocol/infestor/smocker"

	"github.com/stretchr/testify/require"
)

func callSetzer(params ...string) (string, error) {
	out, err := exec.Command("setzer", params...).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func TestETHBTC(t *testing.T) {
	api := smocker.API{
		Host: "http://smocker",
		Port: 8081,
	}

	err := infestor.NewMocksBuilder().
		Reset().
		Add(origin.NewExchange("binance").WithSymbol("ETH/BTC").WithPrice(1)).
		Add(origin.NewExchange("bitfinex").WithSymbol("ETH/BTC").WithPrice(1)).
		Add(origin.NewExchange("coinbase").WithSymbol("ETH/BTC").WithPrice(1)).
		Add(origin.NewExchange("huobi").WithSymbol("ETH/BTC").WithPrice(1)).
		Add(origin.NewExchange("poloniex").WithSymbol("ETH/BTC").WithPrice(1)).
		Add(origin.NewExchange("kraken").WithSymbol("XETH/XXBT").WithPrice(1)).
		Deploy(api)

	// Build your test further
	require.NoError(t, err)

	out, err := callSetzer("price", "ethbtc")

	require.NoError(t, err)
	require.Equal(t, "1.0000000000", out)
}
