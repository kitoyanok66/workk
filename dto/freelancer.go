package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
)

type FreelancerDTO struct {
	ID              uuid.UUID  `json:"id"`
	UserID          uuid.UUID  `json:"user_id"`
	Title           string     `json:"title"`
	Description     string     `json:"description,omitempty"`
	HourlyRate      float64    `json:"hourly_rate"`
	PortfolioURL    string     `json:"portfolio_url,omitempty"`
	ExperienceYears int        `json:"experience_years"`
	Rating          float64    `json:"rating"`
	Skills          []SkillDTO `json:"skills,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type FreelancerRequest struct {
	UserID          uuid.UUID   `json:"user_id"`
	Title           string      `json:"title"`
	Description     string      `json:"description,omitempty"`
	HourlyRate      float64     `json:"hourly_rate"`
	PortfolioURL    string      `json:"portfolio_url,omitempty"`
	ExperienceYears int         `json:"experience_years"`
	SkillIDs        []uuid.UUID `json:"skill_ids,omitempty"`
}

func NewFreelancerDTO(f *domain.Freelancer) *FreelancerDTO {
	if f == nil {
		return nil
	}

	var skills []SkillDTO
	for _, s := range f.Skills {
		skills = append(skills, *NewSkillDTO(&s))
	}

	return &FreelancerDTO{
		ID:              f.ID,
		UserID:          f.UserID,
		Title:           f.Title,
		Description:     f.Description,
		HourlyRate:      f.HourlyRate,
		PortfolioURL:    f.PortfolioURL,
		ExperienceYears: f.ExperienceYears,
		Rating:          f.Rating,
		Skills:          skills,
		CreatedAt:       f.CreatedAt,
		UpdatedAt:       f.UpdatedAt,
	}
}

func (dto *FreelancerDTO) ToDomain() *domain.Freelancer {
	if dto == nil {
		return nil
	}

	var skills []domain.Skill
	for _, s := range dto.Skills {
		skills = append(skills, *s.ToDomain())
	}

	return &domain.Freelancer{
		ID:              dto.ID,
		UserID:          dto.UserID,
		Title:           dto.Title,
		Description:     dto.Description,
		HourlyRate:      dto.HourlyRate,
		PortfolioURL:    dto.PortfolioURL,
		ExperienceYears: dto.ExperienceYears,
		Rating:          dto.Rating,
		Skills:          skills,
		CreatedAt:       dto.CreatedAt,
		UpdatedAt:       dto.UpdatedAt,
	}
}
