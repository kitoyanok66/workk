package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
)

type AuthDTO struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Provider   string    `json:"provider"`
	ExternalID string    `json:"external_id"`
	Username   string    `json:"username"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type AuthRequest struct {
	UserID     uuid.UUID `json:"user_id,omitempty"`
	Provider   string    `json:"provider"`
	ExternalID string    `json:"external_id"`
	Username   string    `json:"username"`
}

type AuthResponse struct {
	Token string   `json:"token"`
	User  *UserDTO `json:"user"`
}

func NewAuthDTO(a *domain.Auth) *AuthDTO {
	if a == nil {
		return nil
	}
	return &AuthDTO{
		ID:         a.ID,
		UserID:     a.UserID,
		Provider:   a.Provider,
		ExternalID: a.ExternalID,
		Username:   a.Username,
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
	}
}

func (dto *AuthDTO) ToDomain() *domain.Auth {
	if dto == nil {
		return nil
	}
	return &domain.Auth{
		ID:         dto.ID,
		UserID:     dto.UserID,
		Provider:   dto.Provider,
		ExternalID: dto.ExternalID,
		Username:   dto.Username,
		CreatedAt:  dto.CreatedAt,
		UpdatedAt:  dto.UpdatedAt,
	}
}
