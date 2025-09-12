package entity

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID           int64          `gorm:"column:id;primaryKey;autoIncrement"`
	AuthID       string         `gorm:"column:auth_id;not null"`
	TypeCategory string         `gorm:"column:type;not null"`
	Name         string         `gorm:"column:name;not null"`
	Qty          int            `gorm:"column:qty;not null"`
	Price        int64          `gorm:"column:price;not null"`
	Sku          string         `gorm:"column:sku;not null"`
	FileId       string         `gorm:"column:file_id;not null"` // references files.fileId (string)
	CreatedAt    time.Time      `gorm:"column:created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`

	File File `gorm:"foreignKey:FileId;references:FileID"`
}

func (Product) TableName() string {
	return "products"
}
