package skills

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
)

type SkillService interface {
	GetAll(ctx context.Context) ([]*domain.Skill, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Skill, error)
	Create(ctx context.Context, name, category, description string) (*domain.Skill, error)
	Update(ctx context.Context, id uuid.UUID, name, category, description string) (*domain.Skill, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type skillService struct {
	repo SkillRepository
}

func NewSkillService(repo SkillRepository) SkillService {
	return &skillService{repo: repo}
}

func (s *skillService) GetAll(ctx context.Context) ([]*domain.Skill, error) {
	return s.repo.GetAll(ctx)
}

func (s *skillService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Skill, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *skillService) Create(ctx context.Context, name, category, description string) (*domain.Skill, error) {
	existing, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("skill with name '%s' already exists", name)
	}

	skill, err := domain.NewSkill(name, category, description)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, skill); err != nil {
		return nil, err
	}

	return skill, nil
}

func (s *skillService) Update(ctx context.Context, id uuid.UUID, name, category, description string) (*domain.Skill, error) {
	skill, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if skill == nil {
		return nil, fmt.Errorf("user not found")
	}

	if err := skill.UpdateSkill(name, category, description); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, skill); err != nil {
		return nil, err
	}

	return skill, err
}

func (s *skillService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
