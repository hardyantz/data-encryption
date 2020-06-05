package conn

import (
	"github.com/go-redis/redis/v8"
)

func RedisConnect() *redis.Client{
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "redis",
		DB:       0,
	})

	return client
}
