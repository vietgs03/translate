package repository

import (
	"context"
	"fmt"

	"github.com/vietgs03/translate/backend/internal/model"
	"gorm.io/gorm"
)

type translationRepo struct {
	db *gorm.DB
}

func NewTranslationRepository(db *gorm.DB) TranslationRepository {
	return &translationRepo{db: db}
}

func (r *translationRepo) Create(ctx context.Context, translation *model.Translation) error {
	return r.db.WithContext(ctx).Create(translation).Error
}

func (r *translationRepo) GetByID(ctx context.Context, id uint) (*model.Translation, error) {
	var translation model.Translation
	if err := r.db.WithContext(ctx).First(&translation, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("translation not found")
		}
		return nil, err
	}
	return &translation, nil
}

func (r *translationRepo) Update(ctx context.Context, translation *model.Translation) error {
	return r.db.WithContext(ctx).Save(translation).Error
}

func (r *translationRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Translation{}, id).Error
}

func (r *translationRepo) List(ctx context.Context, filter TranslationFilter) ([]model.Translation, error) {
	var translations []model.Translation
	query := r.db.WithContext(ctx)

	if filter.SourceLanguage != "" {
		query = query.Where("source_language = ?", filter.SourceLanguage)
	}
	if filter.TargetLanguage != "" {
		query = query.Where("target_language = ?", filter.TargetLanguage)
	}
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}

	// Pagination
	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	if err := query.Find(&translations).Error; err != nil {
		return nil, err
	}
	return translations, nil
} 