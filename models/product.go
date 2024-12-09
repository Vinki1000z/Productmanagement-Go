package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	UserID                  uint           `json:"user_id"`
	ProductName             string         `json:"product_name"`
	ProductDescription      string         `json:"product_description"`
	ProductImages           pq.StringArray `gorm:"type:text[]" json:"product_images"`
	ProductPrice            float64        `json:"product_price"`
	CompressedProductImages pq.StringArray `gorm:"type:text[]" json:"compressed_product_images"`
}
