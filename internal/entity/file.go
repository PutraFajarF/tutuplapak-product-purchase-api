package entity

import "time"

type File struct {
	ID           int64  `gorm:"primaryKey;autoIncrement"`
	URI          string `gorm:"not null"`
	ThumbnailURI *string
	MimeType     *string
	SizeBytes    *int64
	CreatedAt    time.Time
}
