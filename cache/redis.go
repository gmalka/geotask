package cache

import (
	"fmt"

	"github.com/go-redis/redis"
)

func NewRedisClient(host, port string) *redis.Client {
	// реализуйте создание клиента для Redis

	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
	})
}
