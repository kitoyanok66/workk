package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
)

type MatchDTO struct {
	ID           uuid.UUID `json:"id"`
	FreelancerID uuid.UUID `json:"freelancer_id"`
	ProjectID    uuid.UUID `json:"project_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type MatchRequest struct {
	FreelancerID uuid.UUID `json:"freelancer_id"`
	ProjectID    uuid.UUID `json:"project_id"`
}

func NewMatchDTO(m *domain.Match) *MatchDTO {
	if m == nil {
		return nil
	}
	return &MatchDTO{
		ID:           m.ID,
		FreelancerID: m.FreelancerID,
		ProjectID:    m.ProjectID,
		CreatedAt:    m.CreatedAt,
	}
}

func (dto *MatchDTO) ToDomain() *domain.Match {
	if dto == nil {
		return nil
	}
	return &domain.Match{
		ID:           dto.ID,
		FreelancerID: dto.FreelancerID,
		ProjectID:    dto.ProjectID,
		CreatedAt:    dto.CreatedAt,
	}
}
