package repository

import (
	"context"

	"github.com/vietgs03/translate/backend/internal/model"
)

type TranslationRepository interface {
	Create(ctx context.Context, translation *model.Translation) error
	GetByID(ctx context.Context, id uint) (*model.Translation, error)
	Update(ctx context.Context, translation *model.Translation) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, filter TranslationFilter) ([]model.Translation, error)
}

type TranslationFilter struct {
	SourceLanguage string
	TargetLanguage string
	Category       string
	Page          int
	PageSize      int
} 