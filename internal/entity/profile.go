package entity

import "time"

type UserProfile struct {
	UserID            string `gorm:"type:uuid;primaryKey"`
	FileID            *int64 `gorm:"type:bigint"`
	BankAccountName   *string
	BankAccountHolder *string
	BankAccountNumber *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
