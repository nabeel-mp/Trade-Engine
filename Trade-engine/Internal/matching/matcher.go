package matching

import (
	"container/list"
	"sort"
	"time"
)

type Trade struct {
	Price     float64   `json:"price"`
	Quantity  float64   `json:"quantity"`
	Side      Side      `json:"side"`
	Timestamp time.Time `json:"timestamp"`
}

func (ob *OrderBook) Match(order *Order) []Trade {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	trades := []Trade{}

	if order.Side == Buy {
		// Sort prices ascending to match lowest Asks first
		var prices []float64
		for p := range ob.Asks {
			prices = append(prices, p)
		}
		sort.Float64s(prices)

		for _, price := range prices {
			if order.Quantity <= 0 || (order.Type == Limit && price > order.Price) {
				break
			}
			processLevel(ob.Asks[price], order, &trades, price, Buy)
			if ob.Asks[price].Len() == 0 {
				delete(ob.Asks, price)
			}
		}
	} else {
		// Sort prices descending to match highest Bids first
		var prices []float64
		for p := range ob.Bids {
			prices = append(prices, p)
		}
		sort.Slice(prices, func(i, j int) bool { return prices[i] > prices[j] })

		for _, price := range prices {
			if order.Quantity <= 0 || (order.Type == Limit && price < order.Price) {
				break
			}
			processLevel(ob.Bids[price], order, &trades, price, Sell)
			if ob.Bids[price].Len() == 0 {
				delete(ob.Bids, price)
			}
		}
	}

	if order.Quantity > 0 && order.Type == Limit {
		ob.addWithoutLock(order)
	}
	return trades
}

// Helper to reduce code duplication
func processLevel(lvl *list.List, order *Order, trades *[]Trade, price float64, side Side) {
	for e := lvl.Front(); e != nil && order.Quantity > 0; {
		maker := e.Value.(*Order)
		qty := min(order.Quantity, maker.Quantity)
		*trades = append(*trades, Trade{Price: price, Quantity: qty, Side: side, Timestamp: time.Now()})
		order.Quantity -= qty
		maker.Quantity -= qty
		next := e.Next()
		if maker.Quantity == 0 {
			lvl.Remove(e)
		}
		e = next
	}
}
