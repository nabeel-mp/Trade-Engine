package candle

import "time"

type Candle struct {
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
	Start  time.Time
}

func Update(c *Candle, price, qty float64) {
	if c.Open == 0 {
		c.Open = price
		c.High = price
		c.Low = price
	}
	c.Close = price
	c.Volume += qty

	if price > c.High {
		c.High = price
	}
	if price < c.Low {
		c.Low = price
	}
}
