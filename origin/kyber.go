package origin

import (
	"fmt"
	"strings"
)

type Kyber struct{}

func (k Kyber) BuildMock(e []ExchangeMock) ([]byte, error) {
	yaml := `
- request:
    method: GET
    path: '/change24h'
  response:
    status: %d
    headers:
      Content-Type: [application/json]
    body: |-
      {
        %s
      }`
	status := 200
	if len(e) > 0 {
		status = e[0].StatusCode
	}
	return []byte(fmt.Sprintf(yaml, status, k.build(e))), nil
}

func (k Kyber) build(mocks []ExchangeMock) string {
	var result []string
	yaml := `
        "%s": {
          "timestamp":%d,
          "token_symbol":"%s",
          "token_decimal":18,
          "rate_eth_now":%.18f
        }`
	for _, e := range k.filter(mocks) {
		result = append(result,
			fmt.Sprintf(yaml,
				fmt.Sprintf("%s_%s", e.Symbol.Quote, e.Symbol.Base),
				e.Timestamp.UnixMilli(),
				e.Symbol.Base,
				e.Price,
			))
	}
	return strings.Join(result, ",")
}

// removes duplicated pairs from list, otherwise it will cause issues
func (k Kyber) filter(mocks []ExchangeMock) map[string]ExchangeMock {
	filtered := map[string]ExchangeMock{}
	for _, e := range mocks {
		filtered[e.Symbol.String()] = e
	}
	return filtered
}
