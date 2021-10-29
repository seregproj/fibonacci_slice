package redisstorage

import (
	"context"
	"fmt"
	"strconv"
	"time"
)
import "github.com/go-redis/redis/v8"

type RedisCache struct {
	client  *redis.Client
	expires time.Duration
}

func NewRedisCache(addr string, db int, expires time.Duration) *RedisCache {
	client := redis.NewClient(&redis.Options{Addr: addr, DB: db, Password: ""})

	return &RedisCache{client: client, expires: expires}
}

func (rc *RedisCache) Get(ctx context.Context, k uint64) (uint64, error) {
	ks := strconv.FormatUint(k, 10)

	v, err := rc.client.Get(ctx, ks).Result()
	if err != nil {
		return 0, fmt.Errorf("cant get result: %w", err)
	}

	vi, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("cant convert %v to uint64: %w", v, err)
	}

	return vi, nil
}

func (rc *RedisCache) Set(ctx context.Context, k uint64, v uint64) error {
	ks := strconv.FormatUint(k, 10)

	status := rc.client.Set(ctx, ks, v, rc.expires)
	if status.Err() != nil {
		return fmt.Errorf("can write to redis: %w", status.Err())
	}

	return nil
}
