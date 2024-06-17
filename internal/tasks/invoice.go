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
	TypeInvoiceGeneration = "invoice:generate"
)

type InvoiceGenerationPayload struct {
	InvoiceID int
	CartID    int
	UserID    int
	Data      interface{}
}

func NewInvoiceGenerationTask(invoiceID int, cartID int, userID int, data interface{}) (*asynq.Task, error) {
	payload, err := json.Marshal(InvoiceGenerationPayload{InvoiceID: invoiceID, CartID: cartID, UserID: userID, Data: data})
	if err != nil {
		return nil, err
	}
	log.Printf("Creating invoice task for Invoice: invoice_id=%ds", invoiceID)

	return asynq.NewTask(TypeInvoiceGeneration, payload), nil
}

func HandleInvoiceGenerationTask(ctx context.Context, t *asynq.Task) error {
	var p InvoiceGenerationPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Generating invoice for: invoice_id=%ds", p.InvoiceID)

	var user models.User

	result := database.DB.First(&user, p.UserID)

	if result.Error != nil {
		message := fmt.Sprintf("no user found with this id: %d", p.UserID)
		log.Println(message)
		return errors.New(message)
	}

	var cart models.Cart

	result = database.DB.First(&cart, p.CartID)

	if result.Error != nil {
		message := fmt.Sprintf("no cart found with this id: %d", p.CartID)
		log.Println(message)
		return errors.New(message)
	}

	// generate receipt
	filePath, err := util.GeneratePdf(cart, user)
	if err != nil {
		message := fmt.Sprintf("error generating receipt for this id: %d", p.CartID)
		log.Println(message)
		return errors.New(message)
	}

	// send receipt to user
	util.SendMailWithTemplate("purchase", user.Email, "Successful purchase", struct{}{}, filePath)

	return nil
}
