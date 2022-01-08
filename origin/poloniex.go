package origin

import (
	"fmt"
	"net/http"
	"strings"
)

type Poloniex struct{}

func (p Poloniex) BuildMock(e []ExchangeMock) ([]byte, error) {
	yaml := `
- request:
    method: GET
    path: '/public'
    query_params:
      command: 'returnTicker'
  response:
    status: %d
    headers:
      Content-Type: [application/json]
    body: |-
      {
        %s
      }`
	status := http.StatusOK
	if len(e) > 0 {
		status = e[0].StatusCode
	}
	return []byte(fmt.Sprintf(yaml, status, p.build(e))), nil
}

func (p Poloniex) build(mocks []ExchangeMock) string {
	var result []string
	yaml := `
        "%s": {
          "last": "%.8f",
          "lowestAsk": "%.8f",
          "highestBid": "%.8f",
          "baseVolume": "%.8f"
        }`
	for _, e := range p.filter(mocks) {
		result = append(
			result,
			fmt.Sprintf(
				yaml,
				fmt.Sprintf("%s_%s", e.Symbol.Quote, e.Symbol.Base),
				e.Price,
				e.Ask,
				e.Bid,
				e.Volume,
			),
		)
	}
	return strings.Join(result, ",")
}

// removes duplicated pairs from list, otherwise it will cause issues
func (p Poloniex) filter(mocks []ExchangeMock) map[string]ExchangeMock {
	filtered := map[string]ExchangeMock{}
	for _, e := range mocks {
		filtered[e.Symbol.String()] = e
	}
	return filtered
}
