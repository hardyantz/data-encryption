package conn

import (
	"os"

	"github.com/go-redis/redis/v8"
)

func RedisConnect() *redis.Client{
	host := os.Getenv("REDIS_HOST")
	pass := os.Getenv("REDIS_PASS")
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pass,
		DB:       0,
	})

	return client
}
