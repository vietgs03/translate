package main

import (
	"fmt"
	"log"
	
	"github.com/vietgs03/translate/backend/internal/config"
	"github.com/vietgs03/translate/backend/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"github.com/vietgs03/translate/backend/internal/middleware"
	"github.com/vietgs03/translate/backend/internal/handler"
	"github.com/vietgs03/translate/backend/internal/openai"
	"github.com/vietgs03/translate/backend/internal/repository"
	"github.com/vietgs03/translate/backend/internal/service"
	"github.com/vietgs03/translate/backend/internal/cache"
	"go.uber.org/zap"
	"github.com/vietgs03/translate/backend/internal/service/google"
	"github.com/vietgs03/translate/backend/internal/service/translator"
)

type App struct {
	config            *config.Config
	db               *gorm.DB
	redis            *redis.Client
	fiber            *fiber.App
	openai           *openai.Client
	cache            *cache.TranslationCache
	logger           *zap.Logger
	authService      service.AuthService
	translationService service.TranslationService
	authHandler     *handler.AuthHandler
	translationHandler *handler.TranslationHandler
}

func main() {
	app, err := initApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	// Start server
	serverAddr := fmt.Sprintf(":%s", app.config.ServerPort)
	log.Printf("Server starting on port %s", app.config.ServerPort)
	if err := app.fiber.Listen(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initApp() (*App, error) {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	// Initialize Redis
	redisClient, err := database.NewRedisClient(&cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis: %v", err)
	}

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Initialize cache
	translationCache := cache.NewTranslationCache(redisClient)

	// Initialize OpenAI client with rate limiter
	openaiClient := openai.NewClient(&cfg.OpenAI, redisClient)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	translationRepo := repository.NewTranslationRepository(db)

	// Initialize translation service
	var translatorService translator.Translator
	if cfg.Google.GeminiAPIKey != "" {
		// Use Gemini if API key is provided
		translatorService, err = google.NewTranslateService(cfg.Google.GeminiAPIKey)
	} else {
		// Fallback to OpenAI
		translatorService = openaiClient
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create translator service: %v", err)
	}

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.JWT)
	translationService := service.NewTranslationService(
		translationRepo,
		translationCache,
		translatorService,
	)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	translationHandler := handler.NewTranslationHandler(translationService)

	// Create Fiber app with custom error handler
	fiberApp := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Add middleware
	fiberApp.Use(middleware.Logger(logger))
	fiberApp.Use(recover.New())

	app := &App{
		config:            cfg,
		db:               db,
		redis:            redisClient,
		fiber:            fiberApp,
		openai:           openaiClient,
		cache:            translationCache,
		logger:           logger,
		authService:      authService,
		translationService: translationService,
		authHandler:     authHandler,
		translationHandler: translationHandler,
	}

	// Setup routes
	setupRoutes(app)

	return app, nil
}

func setupRoutes(app *App) {
	api := app.fiber.Group("/api/v1")
	
	// Public routes
	public := api.Group("/public")
	public.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"services": fiber.Map{
				"api":    "up",
				"db":     "up",
				"redis":  "up",
			},
		})
	})

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/register", middleware.ValidateRequest(&service.RegisterInput{}), app.authHandler.Register)
	auth.Post("/login", middleware.ValidateRequest(&service.LoginInput{}), app.authHandler.Login)

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.JWTAuth(app.config.JWT.SecretKey))
	protected.Use(middleware.RateLimiter())

	// Admin routes
	admin := protected.Group("/admin")
	admin.Use(middleware.RequireRole("admin"))
	admin.Put("/users/:id/role", app.authHandler.UpdateRole)

	// Translation routes with role-based access
	translations := protected.Group("/translations")
	
	// Create translation - requires validation
	translations.Post("/", 
		middleware.ValidateRequest(&service.CreateTranslationInput{}),
		middleware.RequireRole("user", "translator", "admin"),
		app.translationHandler.Create,
	)

	// Read operations - any authenticated user
	translations.Get("/:id", app.translationHandler.Get)
	translations.Get("/", app.translationHandler.List)

	// Update operations - requires translator role
	translations.Put("/:id",
		middleware.ValidateRequest(&service.UpdateTranslationInput{}),
		middleware.RequireRole("translator", "admin"),
		app.translationHandler.Update,
	)

	// Delete operations - requires admin role
	translations.Delete("/:id",
		middleware.RequireRole("admin"),
		app.translationHandler.Delete,
	)
}