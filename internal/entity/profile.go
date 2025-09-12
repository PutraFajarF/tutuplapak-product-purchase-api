package entity

import (
	"database/sql"
	"time"
)

type Profile struct {
	ID                string         `gorm:"column:id;primaryKey"`
	UserID            string         `gorm:"column:user_id;uniqueIndex;not null"`
	FileID            sql.NullString `gorm:"column:file_id"`
	BankAccountName   sql.NullString `gorm:"column:bank_account_name"`
	BankAccountHolder sql.NullString `gorm:"column:bank_account_holder"`
	BankAccountNumber sql.NullString `gorm:"column:bank_account_number"`
	CreatedAt         time.Time      `gorm:"column:created_at"`
	UpdatedAt         sql.NullTime   `gorm:"column:updated_at"`
}

func (Profile) TableName() string {
	return "profiles"
}
