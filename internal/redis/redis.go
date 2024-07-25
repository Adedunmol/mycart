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

func Init(redisAddress string) {
	once.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     redisAddress,
			Password: "",
			DB:       0,
		})

		_, err := redisClient.Do(context.Background(), "CONFIG", "SET", "notify-keyspace-events", "KEA").Result()

		if err != nil {
			logger.Error.Println("error setting up publish event", err)
		}

		// subscribe to events published in the keyevent channel, specifically for expired events
		pubsub := redisClient.PSubscribe(context.Background(), "__keyevent@0__:expired")

		for {
			message, err := pubsub.ReceiveMessage(context.Background())

			if err != nil {
				logger.Error.Println("error receiving messages from pub/sub channel")
			}

			fmt.Printf("Keyspace event recieved %v  \n", message.String())

			key := message.Payload

			if strings.Contains(key, "shadowKey") {
				key = strings.ReplaceAll(key, "shadowKey", "")
				value := strings.Split(key, ":")
				userId, _ := strconv.Atoi(value[1])

				_ = GetCart(userId)

				// write cart data to postgres
				WriteCartToDB(userId)

				DeleteCart(userId)
			}
		}
	})
}

func GetClient() *redis.Client {
	return redisClient
}
