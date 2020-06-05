package config

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

// Redis struct
type Redis struct {
	client      *redis.Client
	cacheExpire time.Duration
}

// Output struct
type Output struct {
	Result interface{}
	Error  error
}

// Cache interface
type Cache interface {
	Get(key string) Output
	Set(key string, data interface{}, age time.Duration) Output
	Del(key string) Output
	Flush() Output
}

// NewCacheRedis Redis's Constructor
func NewCacheRedis(client *redis.Client, cacheExpire time.Duration) *Redis {
	return &Redis{client: client, cacheExpire: cacheExpire}
}

// Get method
func (r *Redis) Get(key string) Output {
	ctx := context.Background()

	res := r.client.Get(ctx, key)
	if res.Err() != nil {
		return Output{Result: nil, Error: res.Err()}
	}

	var jsonBlob []byte
	err := res.Scan(&jsonBlob)
	if err != nil {
		return Output{Result: nil, Error: err}
	}

	return Output{Result: jsonBlob, Error: nil}
}

// Set method, if age = 0 using default expiration, if age < 0 no expiration in cache
func (r *Redis) Set(key string, data interface{}, age time.Duration) Output {
	ctx := context.Background()

	payload, err := json.Marshal(data)
	if err != nil {
		return Output{Result: nil, Error: err}
	}

	if age == 0 {
		age = r.cacheExpire
	}

	res := r.client.Set(ctx, key, payload, age)
	if res.Err() != nil {
		return Output{Result: nil, Error: res.Err()}
	}
	return Output{Result: nil, Error: nil}
}

// Del method
func (r *Redis) Del(key string) Output {
	ctx := context.Background()
	res := r.client.Keys(ctx, key)
	if res.Err() != nil {
		return Output{Result: nil, Error: res.Err()}
	}

	vals, err := res.Result()
	if err != nil {
		return Output{
			Result: nil,
			Error:  res.Err(),
		}
	}

	for _, v := range vals {
		err = r.client.Del(ctx, v).Err()
		if err != nil {
			return Output{Result: nil, Error: err}
		}
	}
	return Output{Result: nil, Error: nil}
}

// Flush method
func (r *Redis) Flush() Output {
	ctx := context.Background()
	res, err := r.client.FlushDB(ctx).Result()
	return Output{Result: res, Error: err}
}