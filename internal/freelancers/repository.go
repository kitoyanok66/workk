package freelancers

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
	"gorm.io/gorm"
)

type FreelancerRepository interface {
	GetAll(ctx context.Context) ([]*domain.Freelancer, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Freelancer, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Freelancer, error)
	Create(ctx context.Context, freelancer *domain.Freelancer) error
	Update(ctx context.Context, freelancer *domain.Freelancer) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type freelancerRepository struct {
	db *gorm.DB
}

func NewSkillRepository(db *gorm.DB) FreelancerRepository {
	return &freelancerRepository{db: db}
}

func (r *freelancerRepository) GetAll(ctx context.Context) ([]*domain.Freelancer, error) {
	var ormFreelancers []FreelancerORM
	if err := r.db.WithContext(ctx).Preload("Skills").Find(&ormFreelancers).Error; err != nil {
		return nil, err
	}
	result := make([]*domain.Freelancer, len(ormFreelancers))
	for i, f := range ormFreelancers {
		result[i] = f.ToDomain()
	}
	return result, nil
}

func (r *freelancerRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Freelancer, error) {
	var ormFreelancer FreelancerORM
	if err := r.db.WithContext(ctx).Preload("Skills").First(&ormFreelancer, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return ormFreelancer.ToDomain(), nil
}

func (r *freelancerRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Freelancer, error) {
	var ormFreelancer FreelancerORM
	if err := r.db.WithContext(ctx).Preload("Skills").First(&ormFreelancer, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return ormFreelancer.ToDomain(), nil
}

func (r *freelancerRepository) Create(ctx context.Context, freelancer *domain.Freelancer) error {
	if freelancer == nil {
		return errors.New("freelancer is nil")
	}
	ormFreelancer := FromDomain(freelancer)
	return r.db.WithContext(ctx).Create(&ormFreelancer).Error
}

func (r *freelancerRepository) Update(ctx context.Context, freelancer *domain.Freelancer) error {
	if freelancer == nil {
		return errors.New("freelancer is nil")
	}
	ormFreelancer := FromDomain(freelancer)
	tx := r.db.WithContext(ctx).Model(&FreelancerORM{}).Where("id = ?", freelancer.ID).Updates(&ormFreelancer)
	if tx.Error != nil {
		return tx.Error
	}
	if err := r.db.WithContext(ctx).Model(&ormFreelancer).Association("Skills").Replace(ormFreelancer.Skills); err != nil {
		return err
	}
	return nil
}

func (r *freelancerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&FreelancerORM{}, "id = ?", id).Error
}
