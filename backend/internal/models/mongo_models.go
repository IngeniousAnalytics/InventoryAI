package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	ID          primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	TenantID    string                 `bson:"tenant_id" json:"tenant_id"`
	WarehouseID string                 `bson:"warehouse_id" json:"warehouse_id"`
	CategoryID  string                 `bson:"category_id" json:"category_id"`
	Name        string                 `bson:"name" json:"name"`
	Description string                 `bson:"description" json:"description"`
	SKU         string                 `bson:"sku" json:"sku"`
	Quantity    int                    `bson:"quantity" json:"quantity"`
	Price       float64                `bson:"price" json:"price"`
	Images      []string               `bson:"images" json:"images"`
	Attributes  map[string]interface{} `bson:"attributes" json:"attributes"` // Flexible schema
	CreatedAt   time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time              `bson:"updated_at" json:"updated_at"`
}
