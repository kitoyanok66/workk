package projects

import (
	"time"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
	"github.com/kitoyanok66/workk/internal/skills"
)

type ProjectORM struct {
	ID          uuid.UUID          `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID      uuid.UUID          `gorm:"type:uuid;not null;index"`
	Title       string             `gorm:"type:varchar(255);not null"`
	Description string             `gorm:"type:text"`
	Budget      float64            `gorm:"type:numeric(10,2)"`
	Deadline    time.Time          `gorm:"type:timestamp"`
	CreatedAt   time.Time          `gorm:"default:now()"`
	UpdatedAt   time.Time          `gorm:"default:now()"`
	Skills      []*skills.SkillORM `gorm:"many2many:project_skills;joinForeignKey:ProjectID;joinReferences:SkillID;constraint:OnDelete:CASCADE"`
}

func (ProjectORM) TableName() string {
	return "projects"
}

func (o *ProjectORM) ToDomain() *domain.Project {
	if o == nil {
		return nil
	}

	var skillDomains []domain.Skill
	for _, s := range o.Skills {
		skillDomains = append(skillDomains, *s.ToDomain())
	}

	return &domain.Project{
		ID:          o.ID,
		UserID:      o.UserID,
		Title:       o.Title,
		Description: o.Description,
		Budget:      o.Budget,
		Deadline:    o.Deadline,
		Skills:      skillDomains,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}
}

func FromDomain(p *domain.Project) *ProjectORM {
	if p == nil {
		return nil
	}

	var skillORMs []*skills.SkillORM
	for _, s := range p.Skills {
		skillORMs = append(skillORMs, skills.FromDomain(&s))
	}

	return &ProjectORM{
		ID:          p.ID,
		UserID:      p.UserID,
		Title:       p.Title,
		Description: p.Description,
		Budget:      p.Budget,
		Deadline:    p.Deadline,
		Skills:      skillORMs,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
