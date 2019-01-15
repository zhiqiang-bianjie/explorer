package cache

import (
	"github.com/go-redis/redis"
	"github.com/irisnet/explorer/backend/logger"
	"time"
)

var client Cache

type RedisClient struct {
	*redis.Client
}

func (r RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	err := r.Client.Set(key, value, expiration).Err()
	if err != nil {
		logger.Error("redis set error", logger.String("err", err.Error()))
	}
	return err
}

func (r RedisClient) Get(key string) ([]byte, error) {
	val, err := r.Client.Get(key).Bytes()
	if err != nil {
		logger.Warn("get value from redis,return nil", logger.String("key", key))
	}
	return val, err
}

func init() {
	cli := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := cli.Ping().Result()
	logger.Info("redis ping", logger.Any("pong", pong), logger.Any("err", err.Error()))

	client = RedisClient{
		cli,
	}
}

type Cache interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) ([]byte, error)
}

func Instance() Cache {
	return client
}
