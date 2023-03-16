package model

import "time"

type User struct {
	ID         uint64    `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Password   string    `json:"-"`
	TelegramID int64     `json:"telegram_id"`
	Token      string    `gorm:"-" json:"token,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
