package users

import (
	"time"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
)

type UserORM struct {
	ID               uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TelegramUserID   int64     `gorm:"unique;not null"`
	TelegramUsername string    `gorm:"type:text"`
	Fullname         string    `gorm:"column:full_name;type:text;not null"`
	CreatedAt        time.Time `gorm:"default:now()"`
	UpdatedAt        time.Time `gorm:"default:now()"`
}

func (UserORM) TableName() string {
	return "users"
}

func (o *UserORM) ToDomain() *domain.User {
	if o == nil {
		return nil
	}
	return &domain.User{
		ID:               o.ID,
		TelegramUserID:   o.TelegramUserID,
		TelegramUsername: o.TelegramUsername,
		Fullname:         o.Fullname,
		CreatedAt:        o.CreatedAt,
		UpdatedAt:        o.UpdatedAt,
	}
}

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
