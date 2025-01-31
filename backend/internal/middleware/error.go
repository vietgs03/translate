package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vietgs03/translate/backend/internal/types"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	// Convert error to APIError
	apiError := &types.APIError{
		Error: err.Error(),
	}

	code := fiber.StatusInternalServerError

	switch err.(type) {
	case *fiber.Error:
		e := err.(*fiber.Error)
		code = e.Code
	case *types.ValidationError:
		code = fiber.StatusBadRequest
	case *types.UnauthorizedError:
		code = fiber.StatusUnauthorized
	case *types.NotFoundError:
		code = fiber.StatusNotFound
	}

	return c.Status(code).JSON(apiError)
}