package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	ServerPort string `env:"SERVER_PORT" default:"8080"`
	Env        string `env:"ENV" default:"development"`
	Database   DatabaseConfig
	Redis      RedisConfig
	OpenAI     OpenAIConfig
	JWT        JWTConfig
	Google     GoogleConfig
}

type DatabaseConfig struct {
	Host     string `env:"POSTGRES_HOST" default:"localhost"`
	Port     string `env:"POSTGRES_PORT" default:"5432"`
	User     string `env:"POSTGRES_USER" default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" default:"postgres"`
	DBName   string `env:"POSTGRES_DB" default:"itdev_translator"`
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" default:"localhost"`
	Port     string `env:"REDIS_PORT" default:"6379"`
	Password string `env:"REDIS_PASSWORD" default:""`
	DB       int    `env:"REDIS_DB" default:"0"`
}

type OpenAIConfig struct {
	APIKey string `env:"OPENAI_API_KEY" default:""`
}

type JWTConfig struct {
	SecretKey string `env:"JWT_SECRET_KEY" default:"your-secret-key"`
	ExpiresIn int    `env:"JWT_EXPIRES_IN" default:"24"` // hours
}

type GoogleConfig struct {
	ProjectID         string `env:"GOOGLE_PROJECT_ID"`
	CredentialsFile   string `env:"GOOGLE_APPLICATION_CREDENTIALS"`
	GeminiAPIKey      string `env:"GOOGLE_GEMINI_API_KEY"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		// Don't return error if .env file doesn't exist
		// This allows using environment variables without .env file
	}

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		redisDB = 0 // Use default if not set or invalid
	}

	jwtExpiresIn, err := strconv.Atoi(getEnvWithDefault("JWT_EXPIRES_IN", "24"))
	if err != nil {
		jwtExpiresIn = 24 // default to 24 hours if invalid
	}

	return &Config{
		ServerPort: getEnvWithDefault("SERVER_PORT", "8080"),
		Env:        getEnvWithDefault("ENV", "development"),
		Database: DatabaseConfig{
			Host:     getEnvWithDefault("POSTGRES_HOST", "localhost"),
			Port:     getEnvWithDefault("POSTGRES_PORT", "5432"),
			User:     getEnvWithDefault("POSTGRES_USER", "postgres"),
			Password: getEnvWithDefault("POSTGRES_PASSWORD", "postgres"),
			DBName:   getEnvWithDefault("POSTGRES_DB", "itdev_translator"),
		},
		Redis: RedisConfig{
			Host:     getEnvWithDefault("REDIS_HOST", "localhost"),
			Port:     getEnvWithDefault("REDIS_PORT", "6379"),
			Password: getEnvWithDefault("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		OpenAI: OpenAIConfig{
			APIKey: getEnvWithDefault("OPENAI_API_KEY", ""),
		},
		JWT: JWTConfig{
			SecretKey: getEnvWithDefault("JWT_SECRET_KEY", "your-secret-key"),
			ExpiresIn: jwtExpiresIn,
		},
		Google: GoogleConfig{
			ProjectID:         getEnvWithDefault("GOOGLE_PROJECT_ID", ""),
			CredentialsFile:   getEnvWithDefault("GOOGLE_APPLICATION_CREDENTIALS", ""),
			GeminiAPIKey:       getEnvWithDefault("GOOGLE_GEMINI_API_KEY", ""),
		},
	}, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

