package tasks

import (
	"sync"

	"github.com/Adedunmol/mycart/internal/logger"
	"github.com/hibiken/asynq"
)

var (
	client *asynq.Client
	once   sync.Once
)

func Init(redisAddress string) {
	once.Do(func() {
		logger.Logger.Info("setting up connection for asynq queue")

		client = asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddress, Password: "", DB: 0})

	})
}

func Close() {
	logger.Logger.Info("closing connection for asynq queue")

	if client != nil {
		client.Close()
	}
}

func GetClient() *asynq.Client {
	return client
}
