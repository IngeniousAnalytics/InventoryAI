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

func (h *InventoryHandler) UpdateWarehouse(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	warehouseIDStr := c.Params("id")
	warehouseID, err := uuid.Parse(warehouseIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid warehouse id"})
	}

	var req struct {
		Name     string `json:"name"`
		Location string `json:"location"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	// Allow location to be set to empty string intentionally.
	updates["location"] = req.Location
	updates["updated_at"] = time.Now()

	res := h.PG.Model(&models.Warehouse{}).
		Where("id = ? AND tenant_id = ?", warehouseID, tenantID).
		Updates(updates)
	if res.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not update warehouse"})
	}
	if res.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Warehouse not found"})
	}

	var wh models.Warehouse
	if err := h.PG.Where("id = ? AND tenant_id = ?", warehouseID, tenantID).First(&wh).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not fetch updated warehouse"})
	}
	return c.JSON(wh)
}

func (h *InventoryHandler) DeleteWarehouse(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	warehouseIDStr := c.Params("id")
	warehouseID, err := uuid.Parse(warehouseIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid warehouse id"})
	}

	res := h.PG.Where("id = ? AND tenant_id = ?", warehouseID, tenantID).Delete(&models.Warehouse{})
	if res.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not delete warehouse"})
	}
	if res.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Warehouse not found"})
	}
	return c.JSON(fiber.Map{"message": "Warehouse deleted"})
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

func (h *InventoryHandler) UpdateCategory(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	categoryIDStr := c.Params("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid category id"})
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	if req.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name is required"})
	}

	res := h.PG.Model(&models.Category{}).
		Where("id = ? AND tenant_id = ?", categoryID, tenantID).
		Updates(map[string]interface{}{"name": req.Name})
	if res.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not update category"})
	}
	if res.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
	}

	var cat models.Category
	if err := h.PG.Where("id = ? AND tenant_id = ?", categoryID, tenantID).First(&cat).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not fetch updated category"})
	}
	return c.JSON(cat)
}

func (h *InventoryHandler) DeleteCategory(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	categoryIDStr := c.Params("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid category id"})
	}

	res := h.PG.Where("id = ? AND tenant_id = ?", categoryID, tenantID).Delete(&models.Category{})
	if res.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not delete category"})
	}
	if res.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
	}
	return c.JSON(fiber.Map{"message": "Category deleted"})
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

func (h *InventoryHandler) UpdateItem(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	itemIDStr := c.Params("id")
	itemID, err := primitive.ObjectIDFromHex(itemIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid item id"})
	}

	var req struct {
		WarehouseID string                 `json:"warehouse_id"`
		CategoryID  string                 `json:"category_id"`
		Name        string                 `json:"name"`
		Description string                 `json:"description"`
		SKU         string                 `json:"sku"`
		Quantity    *int                   `json:"quantity"`
		Price       *float64               `json:"price"`
		Images      []string               `json:"images"`
		Attributes  map[string]interface{} `json:"attributes"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	set := bson.M{"updated_at": time.Now()}
	if req.WarehouseID != "" {
		set["warehouse_id"] = req.WarehouseID
	}
	if req.CategoryID != "" {
		set["category_id"] = req.CategoryID
	}
	if req.Name != "" {
		set["name"] = req.Name
	}
	// Allow description to be cleared explicitly.
	set["description"] = req.Description
	if req.SKU != "" {
		set["sku"] = req.SKU
	}
	if req.Quantity != nil {
		set["quantity"] = *req.Quantity
	}
	if req.Price != nil {
		set["price"] = *req.Price
	}
	if req.Images != nil {
		set["images"] = req.Images
	}
	if req.Attributes != nil {
		set["attributes"] = req.Attributes
	}

	collection := h.Mongo.Collection("items")
	filter := bson.M{"_id": itemID, "tenant_id": tenantID}
	update := bson.M{"$set": set}

	res, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not update item"})
	}
	if res.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Item not found"})
	}

	var updated models.Item
	err = collection.FindOne(context.TODO(), filter).Decode(&updated)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not fetch updated item"})
	}
	return c.JSON(updated)
}

func (h *InventoryHandler) DeleteItem(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	itemIDStr := c.Params("id")
	itemID, err := primitive.ObjectIDFromHex(itemIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid item id"})
	}

	collection := h.Mongo.Collection("items")
	res, err := collection.DeleteOne(context.TODO(), bson.M{"_id": itemID, "tenant_id": tenantID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not delete item"})
	}
	if res.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Item not found"})
	}
	return c.JSON(fiber.Map{"message": "Item deleted"})
}

