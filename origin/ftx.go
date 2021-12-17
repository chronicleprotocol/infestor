package origin

import "fmt"

type Ftx struct{}

func (b Ftx) BuildMock(e Exchange) ([]byte, error) {
	yaml := `
- request:
    method: GET
    path: '/api/markets/%s/%s'
  response:
    status: %d
    headers:
      Content-Type: application/json
    body: |-
      {
				"success": true,
				"result": {
					"name": "%s",
					"enabled": true,
					"postOnly": false,
					"priceIncrement": 2.5e-6,
					"sizeIncrement": 0.001,
					"minProvideSize": 0.001,
					"last": %f,
					"bid": %f,
					"ask": %f,
					"price": %f,
					"type": "spot",
					"baseCurrency": "%s",
					"quoteCurrency": "%s",
					"underlying": null,
					"restricted": false,
					"highLeverageFeeExempt": true,
					"change1h": -0.0015453653897589824,
					"change24h": 0.05481774512574174,
					"changeBod": 0.01809090909090909,
					"quoteVolume24h": 1351.9492182925,
					"volumeUsd24h": %f
				}
			}`
	return []byte(fmt.Sprintf(
		yaml,
		e.Symbol.Base,
		e.Symbol.Quote,
		e.StatusCode,
		e.Symbol.String(),
		e.Price,
		e.Bid,
		e.Ask,
		e.Price,
		e.Symbol.Base,
		e.Symbol.Quote,
		e.Volume,
	)), nil
}
