package likes

import (
	"time"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
)

type LikeORM struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FromUserID uuid.UUID `gorm:"type:uuid;not null"`
	ToUserID   uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt  time.Time `gorm:"default:now()"`
}

func (LikeORM) TableName() string {
	return "likes"
}

func (l *LikeORM) ToDomain() *domain.Like {
	if l == nil {
		return nil
	}
	return &domain.Like{
		ID:         l.ID,
		FromUserID: l.FromUserID,
		ToUserID:   l.ToUserID,
		CreatedAt:  l.CreatedAt,
	}
}

func FromDomain(l *domain.Like) *LikeORM {
	if l == nil {
		return nil
	}
	return &LikeORM{
		ID:         l.ID,
		FromUserID: l.FromUserID,
		ToUserID:   l.ToUserID,
		CreatedAt:  l.CreatedAt,
	}
}
