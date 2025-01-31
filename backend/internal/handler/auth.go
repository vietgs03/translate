package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vietgs03/translate/backend/internal/service"
	"github.com/vietgs03/translate/backend/internal/errors"
	"github.com/vietgs03/translate/backend/internal/types"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// @Summary Register new user
// @Description Register a new user with username, email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param input body service.RegisterInput true "Registration details"
// @Success 201 {object} model.User
// @Failure 400 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var input service.RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return types.NewValidationError("invalid request body")
	}

	user, err := h.authService.Register(c.Context(), input)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// @Summary Login user
// @Description Login with username and password to get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body service.LoginInput true "Login credentials"
// @Success 200 {object} types.LoginResponse
// @Failure 400 {object} types.APIError
// @Failure 401 {object} types.APIError
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var input service.LoginInput
	if err := c.BodyParser(&input); err != nil {
		return err
	}

	token, err := h.authService.Login(c.Context(), input)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}

func (h *AuthHandler) UpdateRole(c *fiber.Ctx) error {
	var input struct {
		UserID uint   `json:"user_id"`
		Role   string `json:"role" validate:"required,oneof=user translator admin"`
	}

	if err := c.BodyParser(&input); err != nil {
		return errors.NewValidationError("invalid request body: %v", err)
	}

	// Only admin can update roles
	currentUser := c.Locals("user").(*types.JWTClaims)
	if currentUser.Role != "admin" {
		return errors.NewUnauthorizedError("only admin can update roles")
	}

	user, err := h.authService.UpdateRole(c.Context(), input.UserID, input.Role)
	if err != nil {
		return err
	}

	return c.JSON(user)
} 