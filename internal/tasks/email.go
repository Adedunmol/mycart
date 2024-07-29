package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/models"
	"github.com/Adedunmol/mycart/internal/util"
	"github.com/hibiken/asynq"
)

const (
	TypeEmailDelivery = "mail:deliver"
)

type EmailDeliveryPayload struct {
	Template string
	Subject  string
	UserID   int
	Data     interface{}
}

func NewEmailDeliveryTask(template string, subject string, userID int, data interface{}) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailDeliveryPayload{Template: template, Subject: subject, UserID: userID, Data: data})
	if err != nil {
		return nil, err
	}
	log.Printf("Creating mail task for User: user_id=%d", userID)

	return asynq.NewTask(TypeEmailDelivery, payload), nil
}

func HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {
	var p EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Sending Email to User: user_id=%d", p.UserID)

	var user models.User

	result := database.DB.First(&user, p.UserID)

	if result.Error != nil {
		message := fmt.Sprintf("no user found with this id: %d", p.UserID)
		log.Println(message)
		return errors.New(message)
	}

	util.SendMailWithTemplate(p.Template, user.Email, p.Subject, p.Data, "")

	return nil
}
