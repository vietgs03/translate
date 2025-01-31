package openai

import (
	"context"
	"fmt"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/vietgs03/translate/backend/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/vietgs03/translate/backend/internal/service/translator"
)

var _ translator.Translator = (*Client)(nil) // Verify interface implementation

type Client struct {
	client      *openai.Client
	rateLimiter *RateLimiter
}

func NewClient(cfg *config.OpenAIConfig, redis *redis.Client) *Client {
	return &Client{
		client: openai.NewClient(cfg.APIKey),
		rateLimiter: NewRateLimiter(redis, 
			60,                // 60 requests
			time.Minute,      // per minute
		),
	}
}

func (c *Client) Translate(ctx context.Context, text, sourceLang, targetLang string) (string, error) {
	if err := c.rateLimiter.Allow(ctx); err != nil {
		return "", fmt.Errorf("rate limit check failed: %v", err)
	}

	prompt := fmt.Sprintf(
		"Translate the following text from %s to %s. Only return the translated text without any explanation:\n\n%s",
		sourceLang,
		targetLang,
		text,
	)

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7,
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to get translation from OpenAI: %v", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no translation received from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}

// Add Close method to satisfy Translator interface
func (c *Client) Close() error {
	// OpenAI client doesn't need cleanup
	return nil
} 