package main

import (
	"context"
	"fmt"
	"log"
	"os"
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
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("BASE_DB_HOST"),
		os.Getenv("BASE_DB_USER"),
		os.Getenv("BASE_DB_PASSWORD"),
		os.Getenv("BASE_DB_NAME"),
		os.Getenv("BASE_DB_PORT"),
	)

	pgDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to Postgres: %v", err)
	}
	log.Println("Connected to PostgreSQL")

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
		redisURL = "redis://localhost:6379"
	}
	opt, _ := redis.ParseURL(redisURL)
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
		log.Printf("Migration failed: %v", err)
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

	// Categories
	protected.Post("/categories", inventoryHandler.CreateCategory)
	protected.Get("/categories", inventoryHandler.GetCategories)

	// Items
	protected.Post("/items", inventoryHandler.CreateItem)
	protected.Get("/items", inventoryHandler.GetItems)

	// AI
	protected.Post("/ai/queue", aiHandler.QueueImageAnalysis)

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
