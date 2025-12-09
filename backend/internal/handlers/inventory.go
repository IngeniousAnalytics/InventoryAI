package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/inventory_ai/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type InventoryHandler struct {
	PG    *gorm.DB
	Mongo *mongo.Database
}

func NewInventoryHandler(pg *gorm.DB, mongo *mongo.Database) *InventoryHandler {
	return &InventoryHandler{PG: pg, Mongo: mongo}
}

// --- Warehouses ---

func (h *InventoryHandler) CreateWarehouse(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var wh models.Warehouse
	if err := c.BodyParser(&wh); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	wh.TenantID = uuid.MustParse(tenantID)

	if err := h.PG.Create(&wh).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create warehouse"})
	}
	return c.JSON(wh)
}

func (h *InventoryHandler) GetWarehouses(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var warehouses []models.Warehouse
	h.PG.Where("tenant_id = ?", tenantID).Find(&warehouses)
	return c.JSON(warehouses)
}

// --- Categories ---

func (h *InventoryHandler) CreateCategory(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var cat models.Category
	if err := c.BodyParser(&cat); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	cat.TenantID = uuid.MustParse(tenantID)

	if err := h.PG.Create(&cat).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create category"})
	}
	return c.JSON(cat)
}

func (h *InventoryHandler) GetCategories(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var categories []models.Category
	h.PG.Where("tenant_id = ?", tenantID).Find(&categories)
	return c.JSON(categories)
}

// --- Items (MongoDB) ---

func (h *InventoryHandler) CreateItem(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var item models.Item
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	item.TenantID = tenantID
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	item.ID = primitive.NewObjectID()

	collection := h.Mongo.Collection("items")
	_, err := collection.InsertOne(context.TODO(), item)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create item"})
	}
	return c.JSON(item)
}

func (h *InventoryHandler) GetItems(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	collection := h.Mongo.Collection("items")

	cursor, err := collection.Find(context.TODO(), bson.M{"tenant_id": tenantID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not fetch items"})
	}

	var items []models.Item
	if err = cursor.All(context.TODO(), &items); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not parse items"})
	}
	return c.JSON(items)
}
