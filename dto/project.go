package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
)

type ProjectDTO struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Budget      float64    `json:"budget"`
	Deadline    time.Time  `json:"deadline"`
	Skills      []SkillDTO `json:"skills,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type ProjectRequest struct {
	UserID      uuid.UUID   `json:"user_id"`
	Title       string      `json:"title"`
	Description string      `json:"description,omitempty"`
	Budget      float64     `json:"budget"`
	Deadline    time.Time   `json:"deadline"`
	SkillIDs    []uuid.UUID `json:"skill_ids,omitempty"`
}

func NewProjectDTO(p *domain.Project) *ProjectDTO {
	if p == nil {
		return nil
	}

	var skills []SkillDTO
	for _, s := range p.Skills {
		skills = append(skills, *NewSkillDTO(&s))
	}

	return &ProjectDTO{
		ID:          p.ID,
		UserID:      p.UserID,
		Title:       p.Title,
		Description: p.Description,
		Budget:      p.Budget,
		Deadline:    p.Deadline,
		Skills:      skills,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func (dto *ProjectDTO) ToDomain() *domain.Project {
	if dto == nil {
		return nil
	}

	var skills []domain.Skill
	for _, s := range dto.Skills {
		skills = append(skills, *s.ToDomain())
	}

	return &domain.Project{
		ID:          dto.ID,
		UserID:      dto.UserID,
		Title:       dto.Title,
		Description: dto.Description,
		Budget:      dto.Budget,
		Deadline:    dto.Deadline,
		Skills:      skills,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
	}
}
