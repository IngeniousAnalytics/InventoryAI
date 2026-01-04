package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/inventory_ai/backend/internal/handlers"
	"github.com/inventory_ai/backend/internal/middleware"
	"github.com/inventory_ai/backend/internal/models"
	"github.com/inventory_ai/backend/internal/queue"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func openPostgresWithRetry(dsn string, maxAttempts int, delay time.Duration) (*gorm.DB, error) {
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			return db, nil
		}
		lastErr = err
		log.Printf("Postgres not ready (attempt %d/%d): %v", attempt, maxAttempts, err)
		time.Sleep(delay)
	}
	return nil, lastErr
}

func main() {
	// Initialize Fiber
	app := fiber.New(fiber.Config{
		AppName: "InventoryAI Backend",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using defaults")
	}

	// Database Connection (Postgres)
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			os.Getenv("BASE_DB_HOST"),
			os.Getenv("BASE_DB_USER"),
			os.Getenv("BASE_DB_PASSWORD"),
			os.Getenv("BASE_DB_NAME"),
			os.Getenv("BASE_DB_PORT"),
		)
	}

	pgDb, err := openPostgresWithRetry(dsn, 30, 2*time.Second)
	if err != nil {
		log.Fatalf("Could not connect to Postgres: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	// uuid-ossp is optional (we generate UUIDs in-app), but enable it when available.
	if err := pgDb.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		log.Printf("Warning: could not enable uuid-ossp extension: %v", err)
	}

	// Database Connection (MongoDB)
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v", err)
	}

	// Verify Mongo
	if err := mongoClient.Ping(context.TODO(), nil); err != nil {
		log.Printf("Warning: MongoDB ping failed: %v", err)
	} else {
		log.Println("Connected to MongoDB")
	}
	mongoDb := mongoClient.Database("inventory_ai")

	// Redis Connection
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		// Heuristic: if we're configured to talk to docker-compose Postgres service,
		// assume Redis is also the docker-compose service name.
		dbHost := os.Getenv("BASE_DB_HOST")
		databaseURL := os.Getenv("DATABASE_URL")
		if strings.EqualFold(dbHost, "postgres") || strings.Contains(databaseURL, "host=postgres") {
			redisURL = "redis://redis:6379"
		} else {
			redisURL = "redis://localhost:6379"
		}
	}
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Printf("Warning: invalid REDIS_URL %q (%v); falling back to localhost", redisURL, err)
		opt = &redis.Options{Addr: "localhost:6379"}
	}
	redisClient := redis.NewClient(opt)

	// Rate Limiter (100 req / minute)
	limiter := middleware.NewRateLimiter(redisClient, 100, time.Minute)
	app.Use(limiter.LimitMiddleware())

	// RabbitMQ Connection
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}
	rabbitPub, err := queue.NewPublisher(rabbitURL)
	if err != nil {
		log.Printf("Warning: Could not connect to RabbitMQ: %v", err)
	} else {
		defer rabbitPub.Close()
		log.Println("Connected to RabbitMQ")
	}

	// AutoMigrate
	err = pgDb.AutoMigrate(&models.User{}, &models.Tenant{}, &models.Warehouse{}, &models.Category{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Handlers
	authHandler := handlers.NewAuthHandler(pgDb)
	inventoryHandler := handlers.NewInventoryHandler(pgDb, mongoDb)
	aiHandler := handlers.NewAIHandler(rabbitPub)

	// Routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Auth Routes
	auth := v1.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Health Check
	v1.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "inventory-ai-backend",
		})
	})

	// Protected Routes
	protected := v1.Group("/", middleware.Protected())

	// Warehouses
	protected.Post("/warehouses", inventoryHandler.CreateWarehouse)
	protected.Get("/warehouses", inventoryHandler.GetWarehouses)
	protected.Put("/warehouses/:id", inventoryHandler.UpdateWarehouse)
	protected.Delete("/warehouses/:id", inventoryHandler.DeleteWarehouse)

	// Categories
	protected.Post("/categories", inventoryHandler.CreateCategory)
	protected.Get("/categories", inventoryHandler.GetCategories)
	protected.Put("/categories/:id", inventoryHandler.UpdateCategory)
	protected.Delete("/categories/:id", inventoryHandler.DeleteCategory)

	// Items
	protected.Post("/items", inventoryHandler.CreateItem)
	protected.Get("/items", inventoryHandler.GetItems)
	protected.Put("/items/:id", inventoryHandler.UpdateItem)
	protected.Delete("/items/:id", inventoryHandler.DeleteItem)

	// AI
	protected.Post("/ai/queue", aiHandler.QueueImageAnalysis)

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
