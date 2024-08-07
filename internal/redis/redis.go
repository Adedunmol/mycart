package redis

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/Adedunmol/mycart/internal/logger"
	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	once        sync.Once
)

func Init(redisAddress string) error {
	var err error
	once.Do(func() {
		logger.Logger.Info("setting up connection to redis")
		redisClient = redis.NewClient(&redis.Options{
			Addr:     redisAddress,
			Password: "",
			DB:       0,
		})

		_, err = redisClient.Do(context.Background(), "CONFIG", "SET", "notify-keyspace-events", "KEA").Result()

		if err != nil {
			logger.Logger.Error("error setting up publish event")
			logger.Logger.Error(err.Error())
		}

		// subscribe to events published in the keyevent channel, specifically for expired events
		pubsub := redisClient.PSubscribe(context.Background(), "__keyevent@0__:expired")

		for {
			message, err := pubsub.ReceiveMessage(context.Background())

			if err != nil {
				logger.Logger.Error("error receiving messages from pub/sub channel")
			}

			fmt.Printf("Keyspace event recieved %v  \n", message.String())

			key := message.Payload

			if strings.Contains(key, "shadowKey") {
				logger.Logger.Info("getting expired key")
				key = strings.ReplaceAll(key, "shadowKey", "")
				value := strings.Split(key, ":")
				userId, _ := strconv.Atoi(value[2])

				logger.Logger.Info("getting expired cart for user: ")
				_ = GetCart(userId)

				// write cart data to postgres
				logger.Logger.Info("writing expired cart to db")
				WriteCartToDB(userId)

				logger.Logger.Info("deleting expired cart")
				DeleteCart(userId)
			}
		}
	})

	return err
}

func GetClient() *redis.Client {
	return redisClient
}

func Close() {
	logger.Logger.Info("closing connection to redis")
	if redisClient != nil {
		redisClient.Close()
	}
}
