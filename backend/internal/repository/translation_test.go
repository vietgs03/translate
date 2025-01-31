package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vietgs03/translate/backend/internal/model"
	"github.com/vietgs03/translate/backend/internal/testutil"
)

func TestTranslationRepository(t *testing.T) {
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := NewTranslationRepository(db)

	t.Run("Create", func(t *testing.T) {
		translation := &model.Translation{
			SourceText:     "Hello",
			TranslatedText: "Xin chào",
			SourceLanguage: "en",
			TargetLanguage: "vi",
			Category:       "greeting",
		}

		err := repo.Create(context.Background(), translation)
		assert.NoError(t, err)
		assert.NotZero(t, translation.ID)
	})

	t.Run("GetByID", func(t *testing.T) {
		translation := &model.Translation{
			SourceText:     "Test",
			TranslatedText: "Kiểm tra",
			SourceLanguage: "en",
			TargetLanguage: "vi",
		}
		err := repo.Create(context.Background(), translation)
		assert.NoError(t, err)

		found, err := repo.GetByID(context.Background(), translation.ID)
		assert.NoError(t, err)
		assert.Equal(t, translation.SourceText, found.SourceText)
	})

	t.Run("List", func(t *testing.T) {
		filter := TranslationFilter{
			SourceLanguage: "en",
			TargetLanguage: "vi",
			Page:          1,
			PageSize:      10,
		}

		translations, err := repo.List(context.Background(), filter)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(translations), 2)
	})
} 