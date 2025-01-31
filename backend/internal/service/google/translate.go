package google

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"github.com/vietgs03/translate/backend/internal/service/translator"
)

var _ translator.Translator = (*TranslateService)(nil) // Verify interface implementation

type TranslateService struct {
	client *genai.Client
}

func NewTranslateService(apiKey string) (*TranslateService, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %v", err)
	}

	return &TranslateService{
		client: client,
	}, nil
}

func (s *TranslateService) Translate(ctx context.Context, text, sourceLang, targetLang string) (string, error) {
	model := s.client.GenerativeModel("gemini-pro")

	prompt := fmt.Sprintf(
		"Translate the following text from %s to %s. Only return the translation, no explanations:\n\n%s",
		sourceLang,
		targetLang,
		text,
	)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate translation: %v", err)
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no translation generated")
	}

	// Get the response text
	var translation string
	for _, part := range resp.Candidates[0].Content.Parts {
		if textValue, ok := part.(genai.Text); ok {
			translation += string(textValue)
		}
	}

	translation = strings.TrimSpace(translation)
	if translation == "" {
		return "", fmt.Errorf("empty translation received")
	}

	return translation, nil
}

func (s *TranslateService) Close() error {
	return s.client.Close()
} 