package mongo_models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InventoryItem represents a flexible inventory product
type InventoryItem struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TenantID    string             `bson:"tenant_id" json:"tenant_id"` // Indexed
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description,omitempty" json:"description"`
	
	// Core Quantities
	Quantity    float64            `bson:"quantity" json:"quantity"`
	Unit        string             `bson:"unit" json:"unit"` // kg, pcs, liters
	
	// Valuation
	CostPrice   float64            `bson:"cost_price" json:"cost_price"`
	SellingPrice float64           `bson:"selling_price" json:"selling_price"`
	Currency    string             `bson:"currency" json:"currency"`
	
	// Tracking
	BatchNumber string             `bson:"batch_number,omitempty" json:"batch_number"`
	ExpiryDate  *time.Time         `bson:"expiry_date,omitempty" json:"expiry_date"`
	Location    string             `bson:"location,omitempty" json:"location"` // Warehouse A, Shelf B
	
	// Metadata & AI
	Images      []string           `bson:"images" json:"images"` // URLs
	Tags        []string           `bson:"tags" json:"tags"`
	AILog       *AILog             `bson:"ai_log,omitempty" json:"ai_log"`
	
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type AILog struct {
	ScannedAt   time.Time `bson:"scanned_at" json:"scanned_at"`
	Confidence  float64   `bson:"confidence" json:"confidence"`
	DetectedBy  string    `bson:"detected_by" json:"detected_by"` // generic-yolo-v8
	RawOCRText  string    `bson:"raw_ocr_text" json:"raw_ocr_text"`
}

// AuditLog for compliance
type AuditLog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	TenantID  string             `bson:"tenant_id"`
	UserID    string             `bson:"user_id"`
	Action    string             `bson:"action"` // CREATE, UPDATE, DELETE, IMPORT
	Entity    string             `bson:"entity"` // InventoryItem, Settings
	EntityID  string             `bson:"entity_id"`
	Changes   map[string]interface{} `bson:"changes,omitempty"` // Previous vs New
	Timestamp time.Time          `bson:"timestamp"`
}
