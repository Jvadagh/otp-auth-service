package model

import "time"

type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PhoneNumber string    `gorm:"uniqueIndex;not null" json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
}
