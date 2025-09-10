package entity

import "time"

type Product struct {
	ID         int64  `gorm:"primaryKey;autoIncrement"`
	SellerID   string `gorm:"type:uuid;index;not null"`
	CategoryID int64  `gorm:"index;not null"`
	Name       string `gorm:"size:64;not null"`
	Qty        int    `gorm:"not null"`
	Price      int64  `gorm:"not null"` // store in smallest unit
	SKU        string `gorm:"size:32;not null"`
	FileID     int64  `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Category ProductCategory `gorm:"foreignKey:CategoryID"`
	File     File            `gorm:"foreignKey:FileID"`
}

type ProductCategory struct {
	ID   int64  `gorm:"primaryKey;autoIncrement"`
	Code string `gorm:"uniqueIndex;not null"`
	Name string `gorm:"not null"`
}
