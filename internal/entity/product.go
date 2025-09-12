package entity

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID           int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	AuthId       string `json:"auth_id" gorm:"column:auth_id;not null"`
	TypeCategory string `json:"type" gorm:"column:type;not null"`
	Name         string `json:"name" gorm:"column:name;not null"`
	Qty          int    `json:"qty" gorm:"column:qty;not null"`
	Price        int    `json:"price" gorm:"column:price;not null"`
	Sku          string `json:"sku" gorm:"column:sku;not null"`
	FileId       string `json:"file_id" gorm:"column:file_id;not null"`

	File File `json:"file" gorm:"foreignKey:ID;references:FileId"`

	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
}
