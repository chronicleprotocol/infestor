package origin

import "fmt"

type Kraken struct{}

func (k Kraken) BuildMock(e Exchange) ([]byte, error) {
	yaml := `
- request:
    method: GET
    path: '/0/public/Ticker'
    query_params:
      pair: '%s'
  dynamic_response:
    engine: go_template
    script: >
      headers:
        Content-Type: [application/json]
      body: >
        {
          "error": [],
          "result": {
            "%s": {
              "a": [
                "%f",
                "10",
                "10.000"
              ],
              "b": [
                "%f",
                "2",
                "2.000"
              ],
              "c": [
                "%f",
                "0.22651150"
              ],
              "v": [
                "%f",
                "5803.57144830"
              ]
            }
          }
        }`
	symbol := e.Symbol.Format("%s%s")
	return []byte(fmt.Sprintf(yaml, symbol, symbol, e.Ask, e.Bid, e.Price, e.Volume)), nil
}
