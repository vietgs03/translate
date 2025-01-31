package service

import (
	"context"

	"github.com/vietgs03/translate/backend/internal/model"
	"github.com/vietgs03/translate/backend/internal/repository"
)

type TranslationService interface {
	CreateTranslation(ctx context.Context, input CreateTranslationInput) (*model.Translation, error)
	GetTranslation(ctx context.Context, id uint) (*model.Translation, error)
	UpdateTranslation(ctx context.Context, id uint, input UpdateTranslationInput) (*model.Translation, error)
	DeleteTranslation(ctx context.Context, id uint) error
	ListTranslations(ctx context.Context, filter repository.TranslationFilter) ([]model.Translation, error)
}

type CreateTranslationInput struct {
	SourceText     string `json:"source_text" validate:"required,min=1,max=1000"`
	SourceLanguage string `json:"source_language" validate:"required,len=2"`
	TargetLanguage string `json:"target_language" validate:"required,len=2"`
	Context        string `json:"context" validate:"omitempty,max=500"`
	Category       string `json:"category" validate:"omitempty,max=50"`
	CreatedBy      string `json:"created_by" validate:"required"`
}

type UpdateTranslationInput struct {
	TranslatedText string `json:"translated_text" validate:"required,min=1,max=1000"`
	Context        string `json:"context" validate:"omitempty,max=500"`
	Category       string `json:"category" validate:"omitempty,max=50"`
} 