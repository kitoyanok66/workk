package matches

import (
	"time"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
)

type MatchORM struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FreelancerID uuid.UUID `gorm:"type:uuid;not null;index"`
	ProjectID    uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt    time.Time `gorm:"default:now()"`
}

func (MatchORM) TableName() string {
	return "matches"
}

func (m *MatchORM) ToDomain() *domain.Match {
	if m == nil {
		return nil
	}
	return &domain.Match{
		ID:           m.ID,
		FreelancerID: m.FreelancerID,
		ProjectID:    m.ProjectID,
		CreatedAt:    m.CreatedAt,
	}
}

func FromDomain(m *domain.Match) *MatchORM {
	if m == nil {
		return nil
	}
	return &MatchORM{
		ID:           m.ID,
		FreelancerID: m.FreelancerID,
		ProjectID:    m.ProjectID,
		CreatedAt:    m.CreatedAt,
	}
}
