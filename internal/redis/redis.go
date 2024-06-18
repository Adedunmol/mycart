package redis

var (
	client *asynq.Client
	once   sync.Once
)

func Init(redisAddress) {
	
}