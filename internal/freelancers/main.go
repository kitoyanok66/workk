package freelancers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
	"github.com/kitoyanok66/workk/internal/skills"
)

type FreelancerService interface {
	GetAll(ctx context.Context) ([]*domain.Freelancer, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Freelancer, error)
	Create(ctx context.Context, userID uuid.UUID, title, description string, hourlyRate float64, portfolio string, experience int, skillIDs []uuid.UUID) (*domain.Freelancer, error)
	Update(ctx context.Context, id uuid.UUID, title, description string, hourlyRate float64, portfolio string, experience int, skillIDs []uuid.UUID) (*domain.Freelancer, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type freelancerService struct {
	repo       FreelancerRepository
	skillsRepo skills.SkillRepository
}

func NewFreelancerService(repo FreelancerRepository, skillsRepo skills.SkillRepository) FreelancerService {
	return &freelancerService{
		repo:       repo,
		skillsRepo: skillsRepo,
	}
}

func (s *freelancerService) GetAll(ctx context.Context) ([]*domain.Freelancer, error) {
	return s.repo.GetAll(ctx)
}

func (s *freelancerService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Freelancer, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *freelancerService) Create(ctx context.Context, userID uuid.UUID, title, description string, hourlyRate float64, portfolio string, experience int, skillIDs []uuid.UUID) (*domain.Freelancer, error) {
	existing, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("freelancer profile for this user already exists")
	}

	var skillsDomain []domain.Skill
	for _, id := range skillIDs {
		skill, err := s.skillsRepo.GetByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to load skill %s: %w", id, err)
		}
		if skill == nil {
			return nil, fmt.Errorf("skill %s not found", id)
		}
		skillsDomain = append(skillsDomain, *skill)
	}

	freelancer, err := domain.NewFreelancer(userID, title, description, hourlyRate, portfolio, experience, skillsDomain)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, freelancer); err != nil {
		return nil, err
	}

	return freelancer, nil
}

func (s *freelancerService) Update(ctx context.Context, id uuid.UUID, title, description string, hourlyRate float64, portfolio string, experience int, skillIDs []uuid.UUID) (*domain.Freelancer, error) {
	freelancer, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if freelancer == nil {
		return nil, fmt.Errorf("freelancer not found")
	}

	var skillsDomain []domain.Skill
	for _, id := range skillIDs {
		skill, err := s.skillsRepo.GetByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to load skill %s: %w", id, err)
		}
		if skill == nil {
			return nil, fmt.Errorf("skill %s not found", id)
		}
		skillsDomain = append(skillsDomain, *skill)
	}

	if err := freelancer.UpdateFreelancer(title, description, hourlyRate, portfolio, experience, skillsDomain); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, freelancer); err != nil {
		return nil, err
	}

	return freelancer, err
}

func (s *freelancerService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
