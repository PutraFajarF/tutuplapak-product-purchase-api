package entity

import (
	"database/sql"
	"time"
)

type Purchase struct {
	ID                  int64          `gorm:"column:id;primaryKey;autoIncrement"`
	SenderName          string         `gorm:"column:sender_name;not null"`
	SenderContactType   string         `gorm:"column:sender_contact_type;not null"`
	SenderContactDetail string         `gorm:"column:sender_contact_detail;not null"`
	TotalPrice          int64          `gorm:"column:total_price;not null;default:0"`
	Status              string         `gorm:"column:status;not null;default:'PENDING'"`
	CreatedAt           time.Time      `gorm:"column:created_at"`
	UpdatedAt           time.Time      `gorm:"column:updated_at"`
	Items               []PurchaseItem `gorm:"foreignKey:PurchaseID"`
}

func (Purchase) TableName() string {
	return "purchases"
}

type PurchaseItem struct {
	ID               int64          `gorm:"column:id;primaryKey;autoIncrement"`
	PurchaseID       int64          `gorm:"column:purchase_id;index;not null"`
	ProductID        int64          `gorm:"column:product_id;not null"`
	SellerID         string         `gorm:"column:seller_id;index;not null"`
	Name             string         `gorm:"column:name;not null"`
	CategoryCode     string         `gorm:"column:category_code;not null"`
	Price            int64          `gorm:"column:price;not null"`
	SKU              string         `gorm:"column:sku"`
	FileID           sql.NullString `gorm:"column:file_id"` // files.fileId
	FileURI          sql.NullString `gorm:"column:file_uri"`
	FileThumbnailURI sql.NullString `gorm:"column:file_thumbnail_uri"`
	ProductCreatedAt time.Time      `gorm:"column:product_created_at"`
	ProductUpdatedAt time.Time      `gorm:"column:product_updated_at"`
	QtyBefore        int            `gorm:"column:qty_before;not null"`
	BuyQty           int            `gorm:"column:buy_qty;not null"`
	CreatedAt        time.Time      `gorm:"column:created_at"`
}

func (PurchaseItem) TableName() string {
	return "purchase_items"
}

type PurchasePaymentProof struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement"`
	PurchaseID int64     `gorm:"column:purchase_id;index;not null"`
	SellerID   string    `gorm:"column:seller_id;index;not null"`
	FileID     string    `gorm:"column:file_id;not null"` // files.fileId
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func (PurchasePaymentProof) TableName() string {
	return "purchase_payment_proofs"
}
