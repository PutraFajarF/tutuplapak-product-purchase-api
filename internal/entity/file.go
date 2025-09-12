package entity

import "time"

type File struct {
	ID               string    `gorm:"column:id;primaryKey"`
	FileID           string    `gorm:"column:fileId;uniqueIndex;not null"`
	FileUri          string    `gorm:"column:fileUri;not null"`
	FileThumbnailUri string    `gorm:"column:fileThumbnailUri;not null"`
	CreatedAt        time.Time `gorm:"column:created_at"`
}

func (File) TableName() string {
	return "files"
}
