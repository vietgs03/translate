package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vietgs03/translate/backend/internal/errors"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	
	if e, ok := err.(errors.AppError); ok {
		switch e.Type {
		case errors.NotFound:
			code = fiber.StatusNotFound
		case errors.ValidationErr:
			code = fiber.StatusBadRequest
		case errors.Unauthorized:
			code = fiber.StatusUnauthorized
		case errors.DatabaseErr:
			code = fiber.StatusServiceUnavailable
		}
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
} 