package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	ErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       any
	}

	Validator struct {
		validator *validator.Validate
	}
)

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

func (v *Validator) MustValidate(data any) error {
	errs := v.Validate(data)
	if len(errs) == 0 {
		return nil
	}

	return fiber.NewError(fiber.StatusUnprocessableEntity, strings.Join(v.Format(errs), " and "))
}

func (v *Validator) Validate(data any) []ErrorResponse {
	errorResponse := make([]ErrorResponse, 0)

	errs := v.validator.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse

			elem.FailedField = err.Field()
			elem.Tag = err.Tag()
			elem.Value = err.Value()
			elem.Error = true

			errorResponse = append(errorResponse, elem)
		}
	}

	return errorResponse
}

func (v *Validator) Format(errs []ErrorResponse) []string {
	errorMessages := make([]string, 0)
	for _, err := range errs {
		if !err.Error {
			continue
		}
		errorMessages = append(errorMessages, fmt.Sprintf(
			"[%s]: '%v' | Needs to implement '%s'",
			err.FailedField,
			err.Value,
			err.Tag,
		))
	}
	return errorMessages
}

func (v *Validator) ParseBodyAndValidate(c *fiber.Ctx, out any) error {
	err := c.BodyParser(out)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "payload not valid")
	}

	err = v.MustValidate(out)
	if err != nil {
		return err
	}

	return nil
}
