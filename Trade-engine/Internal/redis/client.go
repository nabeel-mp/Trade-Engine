package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func New() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
