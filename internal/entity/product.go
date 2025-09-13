package entity

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID           int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	AuthID       string `gorm:"column:auth_id;not null;uniqueIndex:unique_auth_sku,priority:1;constraint:OnDelete:CASCADE;"`
	TypeCategory string `gorm:"column:type;not null;size:32;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Name         string `json:"name" gorm:"column:name;not null"`
	Qty          int    `json:"qty" gorm:"column:qty;not null"`
	Price        int    `json:"price" gorm:"column:price;not null"`
	Sku          string `gorm:"column:sku;not null;size:32;uniqueIndex:unique_auth_sku,priority:2"`
	FileId       string `gorm:"column:file_id;not null;size:255;constraint:OnDelete:CASCADE;"`

	File File `json:"file" gorm:"foreignKey:id;references:FileId"`

	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
}
