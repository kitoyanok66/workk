package auth

import (
	"time"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
)

type AuthORM struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
	Provider   string    `gorm:"type:text;not null"`
	ExternalID string    `gorm:"type:text;not null"`
	Username   string    `gorm:"type:text"`
	CreatedAt  time.Time `gorm:"default:now()"`
	UpdatedAt  time.Time `gorm:"default:now()"`
}

func (AuthORM) TableName() string {
	return "auth"
}

func (o *AuthORM) ToDomain() *domain.Auth {
	if o == nil {
		return nil
	}
	return &domain.Auth{
		ID:         o.ID,
		UserID:     o.UserID,
		Provider:   o.Provider,
		ExternalID: o.ExternalID,
		Username:   o.Username,
		CreatedAt:  o.CreatedAt,
		UpdatedAt:  o.UpdatedAt,
	}
}

func FromDomain(a *domain.Auth) *AuthORM {
	if a == nil {
		return nil
	}
	return &AuthORM{
		ID:         a.ID,
		UserID:     a.UserID,
		Provider:   a.Provider,
		ExternalID: a.ExternalID,
		Username:   a.Username,
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
	}
}
