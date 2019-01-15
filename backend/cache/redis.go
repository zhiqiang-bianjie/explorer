package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/irisnet/explorer/backend/logger"
	"time"
)

var client *redis.Client

func Set(key string, value interface{}, expiration time.Duration) {
	err := client.Set(key, value, expiration).Err()
	if err != nil {
		logger.Error("redis set error", logger.String("err", err.Error()))
	}
}

func Get(key string) []byte {
	val, err := client.Get(key).Bytes()
	if err != nil {
		logger.Warn("redis get warn")
	}
	return val
}

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
}
