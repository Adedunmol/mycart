package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

const (
	TypeInvoiceGeneration = "invoice:generate"
)

type InvoiceGenerationPayload struct {
	InvoiceID int
	CartID    int
	Data      interface{}
}

func NewInvoiceGenerationTask(invoiceID int, cartID int, data interface{}) (*asynq.Task, error) {
	payload, err := json.Marshal(InvoiceGenerationPayload{InvoiceID: invoiceID, CartID: cartID, Data: data})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeInvoiceGeneration, payload), nil
}

func HandleInvoiceGenerationTask(ctx context.Context, t *asynq.Task) error {
	var p InvoiceGenerationPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Generating invoice for: invoice_id=%ds", p.InvoiceID)

	return nil
}
