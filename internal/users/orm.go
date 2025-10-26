package users

import (
	"time"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
)

type UserORM struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Fullname  string    `gorm:"column:full_name;type:text;not null"`
	Role      string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"default:now()"`
	UpdatedAt time.Time `gorm:"default:now()"`
}

func (UserORM) TableName() string {
	return "users"
}

func (o *UserORM) ToDomain() *domain.User {
	if o == nil {
		return nil
	}
	return &domain.User{
		ID:        o.ID,
		Fullname:  o.Fullname,
		Role:      o.Role,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

func FromDomain(u *domain.User) *UserORM {
	if u == nil {
		return nil
	}
	return &UserORM{
		ID:        u.ID,
		Fullname:  u.Fullname,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
