package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
)

type SkillDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SkillRequest struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description,omitempty"`
}

func NewSkillDTO(s *domain.Skill) *SkillDTO {
	if s == nil {
		return nil
	}
	return &SkillDTO{
		ID:          s.ID,
		Name:        s.Name,
		Category:    s.Category,
		Description: s.Description,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

func (dto *SkillDTO) ToDomain() *domain.Skill {
	if dto == nil {
		return nil
	}
	return &domain.Skill{
		ID:          dto.ID,
		Name:        dto.Name,
		Category:    dto.Category,
		Description: dto.Description,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
	}
}
