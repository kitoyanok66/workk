package skills

import (
	"time"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
)

type SkillORM struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string    `gorm:"type:varchar(255);unique;not null"`
	Category    string    `gorm:"type:varchar(100)"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"default:now()"`
	UpdatedAt   time.Time `gorm:"default:now()"`
}

func (SkillORM) TableName() string {
	return "skills"
}

func (o *SkillORM) ToDomain() *domain.Skill {
	if o == nil {
		return nil
	}
	return &domain.Skill{
		ID:          o.ID,
		Name:        o.Name,
		Category:    o.Category,
		Description: o.Description,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}
}

func FromDomain(s *domain.Skill) *SkillORM {
	if s == nil {
		return nil
	}
	return &SkillORM{
		ID:          s.ID,
		Name:        s.Name,
		Category:    s.Category,
		Description: s.Description,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}
