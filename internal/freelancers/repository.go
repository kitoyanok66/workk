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
	GetBySkillIDs(ctx context.Context, skillIDs []uuid.UUID, currentUserID uuid.UUID) (*domain.Freelancer, error)
}

type freelancerRepository struct {
	db *gorm.DB
}

func NewFreelancerRepository(db *gorm.DB) FreelancerRepository {
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

func (r *freelancerRepository) GetBySkillIDs(ctx context.Context, skillIDs []uuid.UUID, currentUserID uuid.UUID) (*domain.Freelancer, error) {
	if len(skillIDs) == 0 {
		return nil, nil
	}

	query := `
        SELECT f.*
        FROM freelancers f
        JOIN freelancer_skills fs ON f.id = fs.freelancer_id
        WHERE fs.skill_id IN ?
        AND f.user_id NOT IN (
            SELECT to_user_id
            FROM likes
            WHERE from_user_id = ?
        )
        GROUP BY f.id
        ORDER BY COUNT(DISTINCT fs.skill_id) DESC
        LIMIT 1;
    `

	var ormFreelancer FreelancerORM
	if err := r.db.WithContext(ctx).Raw(query, skillIDs, currentUserID).Scan(&ormFreelancer).Error; err != nil {
		return nil, err
	}

	if ormFreelancer.ID != uuid.Nil {
		if err := r.db.WithContext(ctx).Preload("Skills").First(&ormFreelancer, "id = ?", ormFreelancer.ID).Error; err != nil {
			return nil, err
		}
	}

	return ormFreelancer.ToDomain(), nil
}
