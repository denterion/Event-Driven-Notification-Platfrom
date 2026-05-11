package redis

import (
	"context"
	"os"
	"strconv"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *goredis.Client
}

func NewFromEnv() *Client {
	addr := getenv("REDIS_ADDR", "localhost:6379")
	password := getenv("REDIS_PASSWORD", "")
	db := getenvInt("REDIS_DB", 0)

	return &Client{
		rdb: goredis.NewClient(&goredis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
}

func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

func (c *Client) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	return c.rdb.SetNX(ctx, key, value, ttl).Result()
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getenvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
