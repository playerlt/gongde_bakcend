package common

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gongde/config"
	"os"
	"time"
)

var RDB *redis.Client

func InitRedis(addr string) {
	RDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.RedisPassword, // no password set
		DB:       0,                    // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		os.Exit(1)
	}
}
