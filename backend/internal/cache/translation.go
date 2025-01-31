package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vietgs03/translate/backend/internal/model"
)

type TranslationCache struct {
	redis *redis.Client
	ttl   time.Duration
}

func NewTranslationCache(redis *redis.Client) *TranslationCache {
	return &TranslationCache{
		redis: redis,
		ttl:   24 * time.Hour, // Cache for 24 hours
	}
}

func (c *TranslationCache) generateKey(sourceText, sourceLang, targetLang string) string {
	return fmt.Sprintf("translation:%s:%s:%s", sourceLang, targetLang, sourceText)
}

func (c *TranslationCache) Set(ctx context.Context, translation *model.Translation) error {
	key := c.generateKey(translation.SourceText, translation.SourceLanguage, translation.TargetLanguage)
	data, err := json.Marshal(translation)
	if err != nil {
		return fmt.Errorf("failed to marshal translation: %v", err)
	}

	return c.redis.Set(ctx, key, data, c.ttl).Err()
}

func (c *TranslationCache) Get(ctx context.Context, sourceText, sourceLang, targetLang string) (*model.Translation, error) {
	key := c.generateKey(sourceText, sourceLang, targetLang)
	data, err := c.redis.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var translation model.Translation
	if err := json.Unmarshal(data, &translation); err != nil {
		return nil, fmt.Errorf("failed to unmarshal translation: %v", err)
	}

	return &translation, nil
}

func (c *TranslationCache) Delete(ctx context.Context, sourceText, sourceLang, targetLang string) error {
	key := c.generateKey(sourceText, sourceLang, targetLang)
	return c.redis.Del(ctx, key).Err()
} 