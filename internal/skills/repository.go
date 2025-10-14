package skills

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
	"gorm.io/gorm"
)

type SkillRepository interface {
	GetAll(ctx context.Context) ([]*domain.Skill, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Skill, error)
	GetByName(ctx context.Context, name string) (*domain.Skill, error)
	Create(ctx context.Context, skill *domain.Skill) error
	Update(ctx context.Context, skill *domain.Skill) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type skillRepository struct {
	db *gorm.DB
}

func NewSkillRepository(db *gorm.DB) SkillRepository {
	return &skillRepository{db: db}
}

func (r *skillRepository) GetAll(ctx context.Context) ([]*domain.Skill, error) {
	var ormSkills []SkillORM
	if err := r.db.WithContext(ctx).Find(&ormSkills).Error; err != nil {
		return nil, err
	}
	result := make([]*domain.Skill, len(ormSkills))
	for i, s := range ormSkills {
		result[i] = s.ToDomain()
	}
	return result, nil
}

func (r *skillRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Skill, error) {
	var ormSkill SkillORM
	if err := r.db.WithContext(ctx).First(&ormSkill, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return ormSkill.ToDomain(), nil
}

func (r *skillRepository) GetByName(ctx context.Context, name string) (*domain.Skill, error) {
	var ormSkill SkillORM
	if err := r.db.WithContext(ctx).First(&ormSkill, "name = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return ormSkill.ToDomain(), nil
}

func (r *skillRepository) Create(ctx context.Context, skill *domain.Skill) error {
	if skill == nil {
		return errors.New("skill is nil")
	}
	ormSkill := FromDomain(skill)
	return r.db.WithContext(ctx).Create(&ormSkill).Error
}

func (r *skillRepository) Update(ctx context.Context, skill *domain.Skill) error {
	if skill == nil {
		return errors.New("skill is nil")
	}
	ormSkill := FromDomain(skill)
	return r.db.WithContext(ctx).Model(&ormSkill).Where("id = ?", skill.ID).Updates(&ormSkill).Error
}

func (r *skillRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&SkillORM{}, "id = ?", id).Error
}
