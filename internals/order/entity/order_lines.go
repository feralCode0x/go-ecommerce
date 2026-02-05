package entity

import (
	productEntity "go-ecommerce/internals/product/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderLine struct {
	ID        string `json:"id" gorm:"unique;not null;index;primary_key"`
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
	Product   *productEntity.Product
	Quantity  uint            `json:"quantity"`
	Price     float64         `json:"price"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (line *OrderLine) BeforeCreate(tx *gorm.DB) error {
	line.ID = uuid.New().String()
	return nil
}