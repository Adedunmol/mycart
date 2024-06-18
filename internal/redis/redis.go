package redis

import (
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	once        sync.Once
)

func Init(redisAddress string) {
	once.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     redisAddress,
			Password: "",
			DB:       0,
		})
	})
}

func GetClient() *redis.Client {
	return redisClient
}
