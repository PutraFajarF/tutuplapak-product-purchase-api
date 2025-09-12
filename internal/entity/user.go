package entity

import (
	"database/sql"
	"time"
)

type User struct {
	ID        string         `gorm:"column:id;primaryKey"`
	Email     sql.NullString `gorm:"column:email;uniqueIndex"`
	Phone     sql.NullString `gorm:"column:phone;uniqueIndex"`
	Password  string         `gorm:"column:password_hash;not null"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt sql.NullTime   `gorm:"column:updated_at"`
	DeletedAt sql.NullTime   `gorm:"column:deleted_at"`
}

func (User) TableName() string {
	return "users"
}
