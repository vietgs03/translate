package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/vietgs03/translate/backend/internal/errors"
	"github.com/vietgs03/translate/backend/internal/repository"
	"github.com/vietgs03/translate/backend/internal/service"
	"github.com/vietgs03/translate/backend/internal/types"
)

type TranslationHandler struct {
	translationService service.TranslationService
}

func NewTranslationHandler(ts service.TranslationService) *TranslationHandler {
	return &TranslationHandler{
		translationService: ts,
	}
}

// @Summary Create translation
// @Description Create a new translation
// @Tags translations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body service.CreateTranslationInput true "Translation details"
// @Success 201 {object} model.Translation
// @Failure 400 {object} types.APIError
// @Failure 401 {object} types.APIError
// @Failure 403 {object} types.APIError
// @Router /translations [post]
func (h *TranslationHandler) Create(c *fiber.Ctx) error {
	var input service.CreateTranslationInput
	if err := c.BodyParser(&input); err != nil {
		return errors.NewValidationError("invalid request body: %v", err)
	}

	// Get user from JWT claims
	user, ok := c.Locals("user").(*types.JWTClaims)
	if !ok {
		return errors.NewUnauthorizedError("user not authenticated")
	}
	
	input.CreatedBy = user.Username

	translation, err := h.translationService.CreateTranslation(c.Context(), input)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(translation)
}

func (h *TranslationHandler) Get(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return errors.NewValidationError("Invalid ID format")
	}

	translation, err := h.translationService.GetTranslation(c.Context(), uint(id))
	if err != nil {
		return err
	}

	return c.JSON(translation)
}

// @Summary List translations
// @Description Get a list of translations with pagination
// @Tags translations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} types.PaginatedResponse{data=[]model.Translation}
// @Failure 400 {object} types.APIError
// @Failure 401 {object} types.APIError
// @Router /translations [get]
func (h *TranslationHandler) List(c *fiber.Ctx) error {
	filter := repository.TranslationFilter{
		SourceLanguage: c.Query("source_lang"),
		TargetLanguage: c.Query("target_lang"),
		Category:       c.Query("category"),
	}

	// Parse pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	filter.Page = page
	filter.PageSize = pageSize

	translations, err := h.translationService.ListTranslations(c.Context(), filter)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"data": translations,
		"pagination": fiber.Map{
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func (h *TranslationHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return errors.NewValidationError("Invalid ID format")
	}

	var input service.UpdateTranslationInput
	if err := c.BodyParser(&input); err != nil {
		return errors.NewValidationError("Invalid request body")
	}

	translation, err := h.translationService.UpdateTranslation(c.Context(), uint(id), input)
	if err != nil {
		return err
	}

	return c.JSON(translation)
}

func (h *TranslationHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return errors.NewValidationError("Invalid ID format")
	}

	if err := h.translationService.DeleteTranslation(c.Context(), uint(id)); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
} 