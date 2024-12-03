package validator

import (
	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	valiadtor *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	return &CustomValidator{
		valiadtor: validator.New(),
	}
}

func (c *CustomValidator) Validate(i interface{}) error {
	if err := c.valiadtor.Struct(i); err != nil {
		return err
	}
	return nil
}
