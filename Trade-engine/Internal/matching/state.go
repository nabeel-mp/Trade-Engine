package matching

import (
	"Trade-engine/Internal/redis"
	"encoding/json"
)

var rdb = redis.New()

func SaveOrderBook(ob *OrderBook) {
	data, _ := json.Marshal(ob)
	rdb.Set(redis.Ctx, "orderbook", data, 0)
}

func GetStoredOrderBook() (map[string]interface{}, error) {
	val, err := rdb.Get(redis.Ctx, "orderbook").Result()
	if err != nil {
		return nil, err
	}
	var res map[string]interface{}
	json.Unmarshal([]byte(val), &res)
	return res, nil
}

func SaveTrade(trade Trade) {
	data, _ := json.Marshal(trade)
	rdb.LPush(redis.Ctx, "trades", data)
	rdb.LTrim(redis.Ctx, "trades", 0, 99) // Keep last 100 trades
}
