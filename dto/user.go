package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
)

type UserDTO struct {
	ID               uuid.UUID `json:"id"`
	TelegramUserID   int64     `json:"telegram_user_id"`
	TelegramUsername string    `json:"telegram_username,omitempty"`
	Fullname         string    `json:"full_name"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type UserRequest struct {
	TelegramUserID   int64  `json:"telegram_user_id"`
	TelegramUsername string `json:"telegram_username,omitempty"`
	Fullname         string `json:"full_name"`
}

func NewUserDTO(dm *domain.User) *UserDTO {
	if dm == nil {
		return nil
	}
	return &UserDTO{
		ID:               dm.ID,
		TelegramUserID:   dm.TelegramUserID,
		TelegramUsername: dm.TelegramUsername,
		Fullname:         dm.Fullname,
		CreatedAt:        dm.CreatedAt,
		UpdatedAt:        dm.UpdatedAt,
	}
}

func (dto *UserDTO) ToDomain() *domain.User {
	if dto == nil {
		return nil
	}
	return &domain.User{
		ID:               dto.ID,
		TelegramUserID:   dto.TelegramUserID,
		TelegramUsername: dto.TelegramUsername,
		Fullname:         dto.Fullname,
		CreatedAt:        dto.CreatedAt,
		UpdatedAt:        dto.UpdatedAt,
	}
}
