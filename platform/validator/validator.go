package validator

import (
	"fmt"
	"reflect"
	"strings"

	"codeberg.org/mrrizkin/nihil"
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
	v := validator.New()
	v.RegisterCustomTypeFunc(func(value reflect.Value) any {
		nullable := value.Interface().(nihil.NilString)
		if !nullable.Valid {
			return nil
		}
		return nullable.String
	}, nihil.NilString{})

	return &Validator{validator: v}
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

	err := v.validator.Struct(data)
	if err == nil {
		return errorResponse
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return append(errorResponse, ErrorResponse{
			Error:       true,
			FailedField: "payload",
			Tag:         "invalid",
			Value:       data,
		})
	}

	for _, validationError := range validationErrors {
		errorResponse = append(errorResponse, ErrorResponse{
			FailedField: validationError.Field(),
			Tag:         validationError.Tag(),
			Value:       validationError.Value(),
			Error:       true,
		})
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
