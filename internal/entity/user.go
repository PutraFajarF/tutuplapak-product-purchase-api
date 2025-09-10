package entity

import "time"

type User struct {
	ID        string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email     *string `gorm:"uniqueIndex"`
	Phone     *string `gorm:"uniqueIndex"`
	Password  *string
	CreatedAt time.Time
	UpdatedAt time.Time
	Profile   *UserProfile `gorm:"constraint:OnDelete:CASCADE"`
}
