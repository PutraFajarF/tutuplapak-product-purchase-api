package entity

import "time"

type Purchase struct {
	ID                  int64  `gorm:"primaryKey;autoIncrement"`
	SenderName          string `gorm:"not null"`
	SenderContactType   string `gorm:"not null"` // email|phone
	SenderContactDetail string `gorm:"not null"`
	TotalPrice          int64  `gorm:"not null;default:0"`
	Status              string `gorm:"not null;default:'PENDING'"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Items               []PurchaseItem
}

type PurchaseItem struct {
	ID               int64  `gorm:"primaryKey;autoIncrement"`
	PurchaseID       int64  `gorm:"index;not null"`
	ProductID        int64  `gorm:"not null"`
	SellerID         string `gorm:"type:uuid;index;not null"`
	Name             string `gorm:"not null"`
	CategoryCode     string `gorm:"not null"`
	Price            int64  `gorm:"not null"`
	SKU              string
	FileID           *int64
	FileURI          *string
	FileThumbnailURI *string
	ProductCreatedAt *time.Time
	ProductUpdatedAt *time.Time
	QtyBefore        int `gorm:"not null"`
	BuyQty           int `gorm:"not null"`
	CreatedAt        time.Time
}

type PurchasePaymentProof struct {
	ID         int64  `gorm:"primaryKey;autoIncrement"`
	PurchaseID int64  `gorm:"not null;index"`
	SellerID   string `gorm:"type:uuid;not null;index"`
	FileID     int64  `gorm:"not null"`
	CreatedAt  time.Time
}
