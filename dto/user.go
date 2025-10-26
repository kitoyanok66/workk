package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
)

type UserDTO struct {
	ID        uuid.UUID `json:"id"`
	Fullname  string    `json:"full_name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRequest struct {
	Fullname string `json:"full_name"`
	Role     string `json:"role"`
}

func NewUserDTO(u *domain.User) *UserDTO {
	if u == nil {
		return nil
	}
	return &UserDTO{
		ID:        u.ID,
		Fullname:  u.Fullname,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (dto *UserDTO) ToDomain() *domain.User {
	if dto == nil {
		return nil
	}
	return &domain.User{
		ID:        dto.ID,
		Fullname:  dto.Fullname,
		Role:      dto.Role,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}
