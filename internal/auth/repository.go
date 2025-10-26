package auth

import (
	"context"
	"errors"

	"github.com/kitoyanok66/workk/domain"
	"gorm.io/gorm"
)

type AuthRepository interface {
	GetByProviderAndExternalID(ctx context.Context, provider, externalID string) (*domain.Auth, error)
	Create(ctx context.Context, auth *domain.Auth) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) GetByProviderAndExternalID(ctx context.Context, provider, externalID string) (*domain.Auth, error) {
	var ormAuth AuthORM
	if err := r.db.WithContext(ctx).Where("provider = ? AND external_id = ?", provider, externalID).First(&ormAuth).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return ormAuth.ToDomain(), nil
}

func (r *authRepository) Create(ctx context.Context, auth *domain.Auth) error {
	if auth == nil {
		return errors.New("auth is nil")
	}
	ormAuth := FromDomain(auth)
	return r.db.WithContext(ctx).Create(&ormAuth).Error
}
