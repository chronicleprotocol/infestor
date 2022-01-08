package origin

import (
	"fmt"
)

type Bithumb struct{}

func (b Bithumb) BuildMock(e []ExchangeMock) ([]byte, error) {
	return CombineMocks(e, b.build)
}

func (b Bithumb) build(e ExchangeMock) ([]byte, error) {
	yaml := `
- request:
    method: GET
    path: '/openapi/v1/spot/ticker'
    query_params:
      symbol: %s
  response:
    status: %d
    headers:
      Content-Type: [application/json]
    body: |-
      {
        "data": [
          {
            "p": "0.068600",
            "ver": "5218655",
            "vol": "%f",
            "c": "%f",
            "s": "%s",
            "t": "359.860172000000",
            "v": "4439.616702",
            "h": "0.710000",
            "l": "0.068600"
          }
        ],
        "code": "0",
        "msg": "success",
        "timestamp": %d,
        "startTime": null
      }`

	symbol := e.Symbol.Format("%s-%s")
	return []byte(fmt.Sprintf(
		yaml, symbol, e.StatusCode, e.Volume, e.Price, symbol, e.Timestamp.UnixMilli(),
	)), nil
}
