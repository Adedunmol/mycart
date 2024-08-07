package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/logger"
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
	log.Printf("Creating cart update task for User: user_id=%d", userID)

	return asynq.NewTask(TypeCartUpdate, payload), nil
}

func HandleCartUpdateTask(ctx context.Context, t *asynq.Task) error {
	var p CartUpdatePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Updating cart for User: user_id=%d", p.UserID)

	_, updatedAt := redis.GetCartAndUpdatedAt(int(p.UserID))

	fmt.Println("updated at cart update: ", updatedAt)

	// cart's TTL has expired and the updatedAt key can't be accessed
	if updatedAt == "" {
		logger.Logger.Info("updated at not set for the current cart")

		var cart models.Cart

		database.DB.Where(&models.Cart{BuyerID: uint(p.UserID)}).First(&cart)

		redis.UpdateCartFromDB(int(p.UserID))

		return nil
	}

	newUpdatedAt, _ := strconv.Atoi(updatedAt)

	var cart models.Cart

	database.DB.Where(&models.Cart{BuyerID: uint(p.UserID)}).First(&cart)

	// if the updated at of the redis cart is greater than that of the one in postgres
	// dont bother to update from postgres
	if int64(newUpdatedAt) > cart.UpdatedAt.Unix() {
		return nil
	}

	redis.UpdateCartFromDB(int(p.UserID))

	return nil
}
