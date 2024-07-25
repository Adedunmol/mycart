package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/models"
	"github.com/Adedunmol/mycart/internal/redis"
	"github.com/hibiken/asynq"
)

const (
	TypeCartUpdate = "cart:update"
)

type CartUpdatePayload struct {
	UserID int
}

func NewCartUpdateTask(userID int) (*asynq.Task, error) {
	payload, err := json.Marshal(CartUpdatePayload{UserID: userID})
	if err != nil {
		return nil, err
	}
	log.Printf("Creating cart update task for User: user_id=%ds", userID)

	return asynq.NewTask(TypeCartUpdate, payload), nil
}

func HandleCartUpdateTask(ctx context.Context, t *asynq.Task) error {
	var p CartUpdatePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Updating cart for User: user_id=%ds", p.UserID)

	var user models.User

	result := database.DB.First(&user, p.UserID)

	if result.Error != nil {
		message := fmt.Sprintf("no user found with this id: %d", p.UserID)
		log.Println(message)
		return errors.New(message)
	}

	redis.UpdateCartFromDB(int(user.ID))

	return nil
}
