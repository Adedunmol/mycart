package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

const (
	TypeEmailDelivery = "mail:deliver"
)

type EmailDeliveryPayload struct {
	UserID int
	Data   interface{}
}

func NewEmailDeliveryTask(userID int, data interface{}) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailDeliveryPayload{UserID: userID, Data: data})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeEmailDelivery, payload), nil
}

func HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {
	var p EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Sending Email to User: user_id=%ds", p.UserID)

	return nil
}
