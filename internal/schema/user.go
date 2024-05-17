package schema

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type CreateUser struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required,min=6"`
	Email     string `json:"email" validate:"required,email"`
}

func (u *CreateUser) Valid(ctx context.Context) (problems map[string]string) {
	problems = make(map[string]string)

	v := validator.New(validator.WithRequiredStructEnabled())

	if err := v.Struct(u); err != nil {
		for _, err := range err.(validator.ValidationErrors) {

			switch err.Tag() {
			case "required":
				field := err.Field()
				message := fmt.Sprintf("Field '%s' cannot be blank", err.Field())
				problems[field] = message
			case "email":
				field := err.Field()
				message := fmt.Sprintf("Field '%s' must be a valid email address", err.Field())
				problems[field] = message
			case "len":
				field := err.Field()
				message := fmt.Sprintf("Field '%s' must be exactly %v characters long", err.Field(), err.Param())
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

type LoginUser struct {
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required"`
}

func (u *LoginUser) Valid(ctx context.Context) (problems map[string]string) {
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
