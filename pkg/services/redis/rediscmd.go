package redisdb

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)


func Set(ctx context.Context, key string, value interface{}, expiry time.Duration) error {
	statusCmd := rdb.Set(ctx, key, value, expiry)
	if statusCmd.Err() != nil {
		return statusCmd.Err()
	}
	return nil
}

func Get(ctx context.Context, key string) (string, error) {
	value, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil // Key does not exist
		}
		return "", err
	}
	return value, err
}

func Delete(ctx context.Context, key string) (bool, error) {
	state, err := rdb.Del(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if state == 0 {
		return false, nil
	}
	return true, nil
}