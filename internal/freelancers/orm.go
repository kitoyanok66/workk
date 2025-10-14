package freelancers

import (
	"time"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
)

type FreelancerORM struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;index"`
	Title           string    `gorm:"type:varchar(255);not null"`
	Description     string    `gorm:"type:text"`
	HourlyRate      float64   `gorm:"type:numeric(10,2)"`
	PortfolioURL    string    `gorm:"type:text"`
	ExperienceYears int       `gorm:"type:int"`
	Rating          float64   `gorm:"type:float;default:0"`
	CreatedAt       time.Time `gorm:"default:now()"`
	UpdatedAt       time.Time `gorm:"default:now()"`
}

// ORM -> domain
func (o *FreelancerORM) ToDomain() *domain.Freelancer {
	if o == nil {
		return nil
	}
	return &domain.Freelancer{
		ID:               o.ID,
		TelegramUserID:   o.TelegramUserID,
		TelegramUsername: o.TelegramUsername,
		Fullname:         o.Fullname,
		CreatedAt:        o.CreatedAt,
		UpdatedAt:        o.UpdatedAt,
	}
}

// domain -> ORM
func FromDomain(u *domain.User) *UserORM {
	if u == nil {
		return nil
	}
	return &UserORM{
		ID:               u.ID,
		TelegramUserID:   u.TelegramUserID,
		TelegramUsername: u.TelegramUsername,
		Fullname:         u.Fullname,
		CreatedAt:        u.CreatedAt,
		UpdatedAt:        u.UpdatedAt,
	}
}
