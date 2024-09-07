package redisdb

import (
	"context"
	"sync"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	redis "github.com/redis/go-redis/v9"
)

var (
	rdb  *redis.Client
	once sync.Once
	err  error
)

func InitRedis() error {
	once.Do(func() {
		opts, errParseURL := redis.ParseURL(config.Env.REDIS_URL)
		if errParseURL != nil {
			err = errParseURL
			return
		}

		rdb = redis.NewClient(opts)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if _, errPing := rdb.Ping(ctx).Result(); errPing != nil {
			err = errPing
			return
		}
		err = nil
	})
	return err
}

