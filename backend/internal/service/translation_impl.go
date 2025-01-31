package service

import (
	"context"
	"fmt"
	"log"

	"github.com/vietgs03/translate/backend/internal/errors"
	"github.com/vietgs03/translate/backend/internal/model"
	"github.com/vietgs03/translate/backend/internal/repository"
	"github.com/vietgs03/translate/backend/internal/cache"
	"github.com/vietgs03/translate/backend/internal/service/translator"
)

type translationService struct {
	repo       repository.TranslationRepository
	cache      *cache.TranslationCache
	translator translator.Translator
}

func NewTranslationService(
	repo repository.TranslationRepository,
	cache *cache.TranslationCache,
	translator translator.Translator,
) TranslationService {
	return &translationService{
		repo:       repo,
		cache:      cache,
		translator: translator,
	}
}

func (s *translationService) CreateTranslation(ctx context.Context, input CreateTranslationInput) (*model.Translation, error) {
	// Check cache first
	if cached, err := s.cache.Get(ctx, input.SourceText, input.SourceLanguage, input.TargetLanguage); err == nil && cached != nil {
		return cached, nil
	}

	// Try to find existing translation in database
	existing, err := s.findExistingTranslation(ctx, input)
	if err == nil {
		// Cache the found translation
		if err := s.cache.Set(ctx, existing); err != nil {
			log.Printf("Failed to cache translation: %v", err)
		}
		return existing, nil
	}

	// Get translation from translator service
	translatedText, err := s.translator.Translate(ctx, input.SourceText, input.SourceLanguage, input.TargetLanguage)
	if err != nil {
		return nil, fmt.Errorf("failed to translate text: %v", err)
	}

	translation := &model.Translation{
		SourceText:     input.SourceText,
		TranslatedText: translatedText,
		SourceLanguage: input.SourceLanguage,
		TargetLanguage: input.TargetLanguage,
		Context:        input.Context,
		Category:       input.Category,
		CreatedBy:      input.CreatedBy,
	}

	if err := s.repo.Create(ctx, translation); err != nil {
		return nil, errors.NewDatabaseError("failed to save translation: %v", err)
	}

	// Cache the new translation
	if err := s.cache.Set(ctx, translation); err != nil {
		log.Printf("Failed to cache translation: %v", err)
	}

	return translation, nil
}

func (s *translationService) GetTranslation(ctx context.Context, id uint) (*model.Translation, error) {
	translation, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NewNotFoundError("translation not found")
	}
	return translation, nil
}

func (s *translationService) UpdateTranslation(ctx context.Context, id uint, input UpdateTranslationInput) (*model.Translation, error) {
	translation, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NewNotFoundError("translation not found")
	}

	if input.TranslatedText != "" {
		translation.TranslatedText = input.TranslatedText
	}
	if input.Context != "" {
		translation.Context = input.Context
	}
	if input.Category != "" {
		translation.Category = input.Category
	}

	if err := s.repo.Update(ctx, translation); err != nil {
		return nil, errors.NewDatabaseError("failed to update translation: %v", err)
	}

	return translation, nil
}

func (s *translationService) DeleteTranslation(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return errors.NewDatabaseError("failed to delete translation: %v", err)
	}
	return nil
}

func (s *translationService) ListTranslations(ctx context.Context, filter repository.TranslationFilter) ([]model.Translation, error) {
	translations, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, errors.NewDatabaseError("failed to list translations: %v", err)
	}
	return translations, nil
}

func (s *translationService) findExistingTranslation(ctx context.Context, input CreateTranslationInput) (*model.Translation, error) {
	translations, err := s.repo.List(ctx, repository.TranslationFilter{
		SourceText:     input.SourceText,
		SourceLanguage: input.SourceLanguage,
		TargetLanguage: input.TargetLanguage,
		Page:          1,
		PageSize:      1,
	})
	if err != nil || len(translations) == 0 {
		return nil, fmt.Errorf("no existing translation found")
	}
	return &translations[0], nil
} 