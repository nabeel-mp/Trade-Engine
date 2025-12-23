package matching

type Trade struct {
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
}

func (ob *OrderBook) Match(order *Order) []Trade {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	trades := []Trade{}

	if order.Side == Buy {
		for price, asks := range ob.Asks {
			if order.Type == Limit && price > order.Price {
				continue
			}

			for e := asks.Front(); e != nil && order.Quantity > 0; {
				ask := e.Value.(*Order)
				qty := min(order.Quantity, ask.Quantity)

				trades = append(trades, Trade{Price: price, Quantity: qty})

				order.Quantity -= qty
				ask.Quantity -= qty

				next := e.Next()
				if ask.Quantity == 0 {
					asks.Remove(e)
				}
				e = next
			}
		}
	} else { // SELL
		for price, bids := range ob.Bids {
			if order.Type == Limit && price < order.Price {
				continue
			}

			for e := bids.Front(); e != nil && order.Quantity > 0; {
				bid := e.Value.(*Order)
				qty := min(order.Quantity, bid.Quantity)

				trades = append(trades, Trade{Price: price, Quantity: qty})

				order.Quantity -= qty
				bid.Quantity -= qty

				next := e.Next()
				if bid.Quantity == 0 {
					bids.Remove(e)
				}
				e = next
			}
		}
	}

	if order.Quantity > 0 && order.Type == Limit {
		ob.Add(order)
	}

	return trades
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
