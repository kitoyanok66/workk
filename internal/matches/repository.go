package matches

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
	"gorm.io/gorm"
)

type MatchRepository interface {
	GetAll(ctx context.Context) ([]*domain.Match, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Match, error)
	Create(ctx context.Context, match *domain.Match) error
}

type matchRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) MatchRepository {
	return &matchRepository{db: db}
}

func (r *matchRepository) GetAll(ctx context.Context) ([]*domain.Match, error) {
	var ormMatches []MatchORM
	if err := r.db.WithContext(ctx).Find(&ormMatches).Error; err != nil {
		return nil, err
	}
	result := make([]*domain.Match, len(ormMatches))
	for i, m := range ormMatches {
		result[i] = m.ToDomain()
	}
	return result, nil
}

func (r *matchRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Match, error) {
	var ormMatch MatchORM
	if err := r.db.WithContext(ctx).First(&ormMatch, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return ormMatch.ToDomain(), nil
}

func (r *matchRepository) Create(ctx context.Context, match *domain.Match) error {
	if match == nil {
		return errors.New("match is nil")
	}
	ormMatch := FromDomain(match)
	return r.db.WithContext(ctx).Create(&ormMatch).Error
}
