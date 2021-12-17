package origin

import (
	"fmt"
)

type Upbit struct{}

func (o Upbit) BuildMock(e Exchange) ([]byte, error) {
	yaml := `
- request:
    method: GET
    path: '/v1/ticker'
		query_params:
			markets: %s
  response:
    status: %d
    headers:
      Content-Type: application/json
    body: |-
      [
				{
					"market": "%s",
					"trade_date": "%s",
					"trade_time": "%s",
					"trade_date_kst": "20211217",
					"trade_time_kst": "181713",
					"trade_timestamp": %d,
					"opening_price": 0.08281743,
					"high_price": 0.08345911,
					"low_price": 0.08214412,
					"trade_price": %f,
					"prev_closing_price": 0.08345011,
					"change": "FALL",
					"change_price": 0.00104761,
					"change_rate": 0.0125537282,
					"signed_change_price": -0.00104761,
					"signed_change_rate": -0.0125537282,
					"trade_volume": %f,
					"acc_trade_price": 1.8122924226326650,
					"acc_trade_price_24h": 12.68662859,
					"acc_trade_volume": 21.88766354,
					"acc_trade_volume_24h": 152.69138015,
					"highest_52_week_price": 0.09,
					"highest_52_week_date": "2021-12-10",
					"lowest_52_week_price": 0.02250000,
					"lowest_52_week_date": "2020-12-27",
					"timestamp": %d
				}
			]`

	symbol := fmt.Sprintf("%s-%s", e.Symbol.Quote, e.Symbol.Base)
	return []byte(fmt.Sprintf(
		yaml,
		symbol,
		e.StatusCode,
		symbol,
		e.Timestamp.Format("20060102"),
		e.Timestamp.Format("150405"),
		e.Timestamp.UnixMilli(),
		e.Price,
		e.Volume,
		e.Timestamp.UnixMilli(),
	)), nil
}
