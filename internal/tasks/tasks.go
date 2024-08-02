package tasks

import (
	"log"
	"sync"

	"github.com/Adedunmol/mycart/internal/logger"
	"github.com/hibiken/asynq"
)

var (
	client *asynq.Client
	once   sync.Once
)

const redisAddr = "127.0.0.1:6379"

func Run() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeCartUpdate, HandleCartUpdateTask)
	mux.HandleFunc(TypeInvoiceGeneration, HandleInvoiceGenerationTask)
	mux.HandleFunc(TypeEmailDelivery, HandleEmailDeliveryTask)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

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
