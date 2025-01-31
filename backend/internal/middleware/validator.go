package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/vietgs03/translate/backend/internal/errors"
)

var validate = validator.New()

func ValidateRequest(payload interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.BodyParser(payload); err != nil {
			return errors.NewValidationError("invalid request body: %v", err)
		}

		if err := validate.Struct(payload); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				errorMessages := make([]string, len(validationErrors))
				for i, err := range validationErrors {
					errorMessages[i] = formatValidationError(err)
				}
				return errors.NewValidationError("validation failed: %v", errorMessages)
			}
			return errors.NewValidationError("validation failed")
		}

		c.Locals("validated", payload)
		return c.Next()
	}
}

func formatValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " is required"
	case "min":
		return err.Field() + " must be at least " + err.Param()
	case "max":
		return err.Field() + " must be at most " + err.Param()
	case "len":
		return err.Field() + " must be exactly " + err.Param() + " characters long"
	default:
		return err.Field() + " is invalid"
	}
} 