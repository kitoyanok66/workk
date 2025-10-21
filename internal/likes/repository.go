package likes

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
	"gorm.io/gorm"
)

type LikeRepository interface {
	Create(ctx context.Context, like *domain.Like) error
	ExistsReverse(ctx context.Context, userID, targetID uuid.UUID) (bool, error)
}

type likeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) LikeRepository {
	return &likeRepository{db: db}
}

func (r *likeRepository) Create(ctx context.Context, like *domain.Like) error {
	if like == nil {
		return errors.New("like is nil")
	}
	ormLike := FromDomain(like)
	return r.db.WithContext(ctx).Create(&ormLike).Error
}

func (r *likeRepository) ExistsReverse(ctx context.Context, userID, targetID uuid.UUID) (bool, error) {
	var ormLike LikeORM
	err := r.db.WithContext(ctx).Where("from_user_id = ? AND to_user_id = ?", targetID, userID).First(&ormLike).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
