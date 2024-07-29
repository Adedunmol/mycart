package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Adedunmol/mycart/internal/redis"
	"github.com/hibiken/asynq"
)

const (
	TypeCartDelete = "cart:delete"
)

type CartDeletePayload struct {
	UserID int
}

func NewCartDeleteTask(userID int) (*asynq.Task, error) {
	payload, err := json.Marshal(CartUpdatePayload{UserID: userID})
	if err != nil {
		return nil, err
	}
	log.Printf("creating cart delete task for User: user_id=%d", userID)

	return asynq.NewTask(TypeCartDelete, payload), nil
}

func HandleCartDeleteTask(ctx context.Context, t *asynq.Task) error {
	var p CartDeletePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("deleting cart for User: user_id=%d", p.UserID)

	redis.ClearCartAndDB(int(p.UserID))

	return nil
}
