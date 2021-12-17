package origin

import "fmt"

type Poloniex struct{}

func (b Poloniex) BuildMock(e Exchange) ([]byte, error) {
	yaml := `
- request:
    method: GET
    path: '/public'
    query_params:
      command: 'returnTicker'
  dynamic_response:
    engine: go_template
    script: >
      headers:
        Content-Type: [application/json]
      body: >
        {
          "%s": {
            "id": 148,
            "last": "%f",
            "lowestAsk": "%f",
            "highestBid": "%f",
            "percentChange": "0.01131570",
            "baseVolume": "%f",
            "quoteVolume": "1087.92098487",
            "isFrozen": "0",
            "postOnly": "0",
            "marginTradingEnabled": "1",
            "high24hr": "0.08840000",
            "low24hr": "0.08626057"
          }
        }`
	return []byte(fmt.Sprintf(yaml, e.Symbol.Format("%s%s"), e.Price, e.Ask, e.Bid, e.Volume)), nil
}
