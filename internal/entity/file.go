package entity

import "time"

type File struct {
	ID           string    `gorm:"primaryKey;column:id" json:"id"`
	Uri          string    `gorm:"column:uri" json:"uri"`
	ThumbnailUri string    `gorm:"column:thumbnail_uri" json:"thumbnailUri"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updatedAt"`
}
