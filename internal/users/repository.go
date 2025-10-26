package users

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	var ormUsers []UserORM
	if err := r.db.WithContext(ctx).Find(&ormUsers).Error; err != nil {
		return nil, err
	}
	result := make([]*domain.User, len(ormUsers))
	for i, u := range ormUsers {
		result[i] = u.ToDomain()
	}
	return result, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var ormUser UserORM
	if err := r.db.WithContext(ctx).First(&ormUser, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return ormUser.ToDomain(), nil
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	ormUser := FromDomain(user)
	return r.db.WithContext(ctx).Create(&ormUser).Error
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	ormUser := FromDomain(user)
	return r.db.WithContext(ctx).Model(&ormUser).Where("id = ?", user.ID).Updates(&ormUser).Error
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&UserORM{}, "id = ?", id).Error
}
