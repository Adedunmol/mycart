package schema

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type CreateReviewDto struct {
	Comment string `json:"comment" validate:"required"`
	Rating  uint   `json:"rating" validate:"required"`
}

func (u *CreateReviewDto) Valid(ctx context.Context) (problems map[string]string) {
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
