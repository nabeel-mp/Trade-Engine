package matching

import (
	"container/list"
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

	// Use pre-sorted slices
	var prices *[]float64
	var book map[float64]*list.List

	if order.Side == Buy {
		prices = &ob.AskPrices
		book = ob.Asks
	} else {
		prices = &ob.BidPrices
		book = ob.Bids
	}

	for i := 0; i < len(*prices); {
		price := (*prices)[i]

		// Stop if limit price is hit (Market orders ignore this)
		if order.Type == Limit {
			if order.Side == Buy && price > order.Price {
				break
			}
			if order.Side == Sell && price < order.Price {
				break
			}
		}

		processLevel(book[price], order, &trades, price, order.Side)

		if book[price].Len() == 0 {
			delete(book, price)
			*prices = append((*prices)[:i], (*prices)[i+1:]...) // Remove price
		} else {
			i++
		}

		if order.Quantity <= 0 {
			break
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
