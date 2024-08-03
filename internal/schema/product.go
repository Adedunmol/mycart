package schema

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type CreateProductDto struct {
	Name     string `json:"name" validate:"required"`
	Details  string `json:"details" validate:"required"`
	Price    int    `json:"price" validate:"required,min=0"`
	Quantity int    `json:"quantity" validate:"required,min=1"`
	Category string `json:"category" validate:"required"`
	// Date     time.Time `json:"date"`
}

func (u *CreateProductDto) Valid(ctx context.Context) (problems map[string]string) {
	problems = make(map[string]string)

	v := validator.New(validator.WithRequiredStructEnabled())

	if err := v.Struct(u); err != nil {
		for _, err := range err.(validator.ValidationErrors) {

			switch err.Tag() {
			case "required":
				field := err.Field()
				message := fmt.Sprintf("Field '%s' cannot be blank", err.Field())
				problems[field] = message
			default:
				field := err.Field()
				message := fmt.Sprintf("Field '%s': '%v' must satisfy '%s' '%v' criteria", err.Field(), err.Value(), err.Tag(), err.Param())
				problems[field] = message
			}
		}
	}

	return problems
}

type UpdateProductDto struct {
	Name     string `json:"name"`
	Details  string `json:"details"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
	Category string `json:"category"`
}

func (u *UpdateProductDto) Valid(ctx context.Context) (problems map[string]string) {
	problems = make(map[string]string)

	v := validator.New(validator.WithRequiredStructEnabled())

	if err := v.Struct(u); err != nil {
		for _, err := range err.(validator.ValidationErrors) {

			switch err.Tag() {
			case "required":
				field := err.Field()
				message := fmt.Sprintf("Field '%s' cannot be blank", err.Field())
				problems[field] = message
			default:
				field := err.Field()
				message := fmt.Sprintf("Field '%s': '%v' must satisfy '%s' '%v' criteria", err.Field(), err.Value(), err.Tag(), err.Param())
				problems[field] = message
			}
		}
	}

	return problems
}
