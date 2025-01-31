package main

import (
	"fmt"
	"log"
	
	"github.com/vietgs03/translate/backend/internal/config"
	"github.com/vietgs03/translate/backend/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"github.com/vietgs03/translate/backend/internal/middleware"
	"github.com/vietgs03/translate/backend/internal/handler"
	"github.com/vietgs03/translate/backend/internal/openai"
	"github.com/vietgs03/translate/backend/internal/repository"
	"github.com/vietgs03/translate/backend/internal/service"
	"github.com/vietgs03/translate/backend/internal/cache"
)

type App struct {
	config            *config.Config
	db               *gorm.DB
	redis            *redis.Client
	fiber            *fiber.App
	openai           *openai.Client
	cache            *cache.TranslationCache
	translationService service.TranslationService
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

	// Create Fiber app
	fiberApp := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Add middleware
	fiberApp.Use(logger.New())
	fiberApp.Use(recover.New())

	// Initialize cache
	translationCache := cache.NewTranslationCache(redisClient)

	// Initialize OpenAI client with rate limiter
	openaiClient := openai.NewClient(&cfg.OpenAI, redisClient)

	// Initialize repositories
	translationRepo := repository.NewTranslationRepository(db)

	// Initialize services
	translationService := service.NewTranslationService(translationRepo, translationCache, openaiClient)

	app := &App{
		config:            cfg,
		db:               db,
		redis:            redisClient,
		fiber:            fiberApp,
		openai:           openaiClient,
		cache:            translationCache,
		translationService: translationService,
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

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.JWTAuth(app.config.JWT.SecretKey))
	protected.Use(middleware.RateLimiter())

	// Translation routes with role-based access
	translations := protected.Group("/translations")
	
	// Create translation - requires validation
	translations.Post("/", 
		middleware.ValidateRequest(&service.CreateTranslationInput{}),
		middleware.RequireRole("translator", "admin"),
		translationHandler.Create,
	)

	// Read operations - any authenticated user
	translations.Get("/:id", translationHandler.Get)
	translations.Get("/", translationHandler.List)

	// Update operations - requires translator role
	translations.Put("/:id",
		middleware.ValidateRequest(&service.UpdateTranslationInput{}),
		middleware.RequireRole("translator", "admin"),
		translationHandler.Update,
	)

	// Delete operations - requires admin role
	translations.Delete("/:id",
		middleware.RequireRole("admin"),
		translationHandler.Delete,
	)
}