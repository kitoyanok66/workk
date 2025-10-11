package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserDTO struct {
	ID               uuid.UUID `json:"uuid"`
	TelegramUserID   int64     `json:"telegram_user_id"`
	TelegramUsername *string   `json:"telegram_username,omitempty"`
	Fullname         string    `json:"full_name"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func NewUserDTO(dm domain.User) UserDTO {
	return UserDTO{
		ID:               dm.ID,
		TelegramUserID:   dm.TelegramUserID,
		TelegramUsername: dm.TelegramUsername,
		Fullname:         dm.Fullname,
		CreatedAt:        dm.CreatedAt,
		UpdatedAt:        dm.UpdatedAt,
	}
}
