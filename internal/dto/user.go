package dto

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Withdraw struct {
	Order string
	Sum   float64
}

type CreateUser struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (user *CreateUser) Validate() error {
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		var errArr []string
		for _, e := range errs {
			errArr = append(errArr, e.Error())
		}
		middlewareNamesStr := strings.Join(errArr, ", ")
		return errors.New(middlewareNamesStr)
	}
	return nil
}
