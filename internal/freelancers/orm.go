package freelancers

import (
	"time"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
	"github.com/kitoyanok66/workk/internal/skills"
)

type FreelancerORM struct {
	ID              uuid.UUID          `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID          uuid.UUID          `gorm:"type:uuid;not null;index"`
	Title           string             `gorm:"type:varchar(255);not null"`
	Description     string             `gorm:"type:text"`
	HourlyRate      float64            `gorm:"type:numeric(10,2)"`
	PortfolioURL    string             `gorm:"type:text"`
	ExperienceYears int                `gorm:"type:int"`
	Rating          float64            `gorm:"type:float;default:0"`
	CreatedAt       time.Time          `gorm:"default:now()"`
	UpdatedAt       time.Time          `gorm:"default:now()"`
	Skills          []*skills.SkillORM `gorm:"many2many:freelancer_skills;constraint:OnDelete:CASCADE"`
}

func (o *FreelancerORM) ToDomain() *domain.Freelancer {
	if o == nil {
		return nil
	}

	var skillDomains []domain.Skill
	for _, s := range o.Skills {
		skillDomains = append(skillDomains, *s.ToDomain())
	}

	return &domain.Freelancer{
		ID:              o.ID,
		UserID:          o.UserID,
		Title:           o.Title,
		Description:     o.Description,
		HourlyRate:      o.HourlyRate,
		PortfolioURL:    o.PortfolioURL,
		ExperienceYears: o.ExperienceYears,
		Rating:          o.Rating,
		Skills:          skillDomains,
		CreatedAt:       o.CreatedAt,
		UpdatedAt:       o.UpdatedAt,
	}
}

func FromDomain(f *domain.Freelancer) *FreelancerORM {
	if f == nil {
		return nil
	}

	var skillORMs []*skills.SkillORM
	for _, s := range f.Skills {
		skillORMs = append(skillORMs, skills.FromDomain(&s))
	}

	return &FreelancerORM{
		ID:              f.ID,
		UserID:          f.UserID,
		Title:           f.Title,
		Description:     f.Description,
		HourlyRate:      f.HourlyRate,
		PortfolioURL:    f.PortfolioURL,
		ExperienceYears: f.ExperienceYears,
		Rating:          f.Rating,
		Skills:          skillORMs,
		CreatedAt:       f.CreatedAt,
		UpdatedAt:       f.UpdatedAt,
	}
}
