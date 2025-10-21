package freelancers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"

	"github.com/kitoyanok66/workk/internal/deps"
	"github.com/kitoyanok66/workk/internal/skills"
	"github.com/kitoyanok66/workk/internal/users"
)

type FreelancerService interface {
	GetAll(ctx context.Context) ([]*domain.Freelancer, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Freelancer, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Freelancer, error)
	Create(ctx context.Context, userID uuid.UUID, title, description string, hourlyRate float64, portfolio string, experience int, skillIDs []uuid.UUID) (*domain.Freelancer, error)
	Update(ctx context.Context, id uuid.UUID, title, description string, hourlyRate float64, portfolio string, experience int, skillIDs []uuid.UUID) (*domain.Freelancer, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetFeedForProject(ctx context.Context, projectID uuid.UUID) ([]*domain.Freelancer, error)

	SetProjectDep(dep deps.ProjectFetcher)
}

type freelancerService struct {
	repo       FreelancerRepository
	skillSvc   skills.SkillService
	projectDep deps.ProjectFetcher
	userSvc    users.UserService
}

func NewFreelancerService(repo FreelancerRepository, skillSvc skills.SkillService, userSvc users.UserService) FreelancerService {
	return &freelancerService{
		repo:     repo,
		skillSvc: skillSvc,
		userSvc:  userSvc,
	}
}

func (s *freelancerService) SetProjectDep(dep deps.ProjectFetcher) {
	s.projectDep = dep
}

func (s *freelancerService) GetAll(ctx context.Context) ([]*domain.Freelancer, error) {
	return s.repo.GetAll(ctx)
}

func (s *freelancerService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Freelancer, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *freelancerService) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Freelancer, error) {
	return s.repo.GetByUserID(ctx, userID)
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
		skill, err := s.skillSvc.GetByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to load skill %s: %w", id, err)
		}
		if skill == nil {
			return nil, fmt.Errorf("skill %s not found", id)
		}
		skillsDomain = append(skillsDomain, *skill)
	}

	project, err := s.projectDep.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if project != nil {
		if err := s.projectDep.Delete(ctx, project.ID); err != nil {
			return nil, err
		}
	}

	freelancer, err := domain.NewFreelancer(userID, title, description, hourlyRate, portfolio, experience, skillsDomain)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, freelancer); err != nil {
		return nil, err
	}

	if err := s.userSvc.UpdateRole(ctx, userID, "freelancer"); err != nil {
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
		skill, err := s.skillSvc.GetByID(ctx, id)
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

func (s *freelancerService) GetFeedForProject(ctx context.Context, projectID uuid.UUID) ([]*domain.Freelancer, error) {
	project, err := s.projectDep.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	var skillIDs []uuid.UUID
	for _, sk := range project.Skills {
		skillIDs = append(skillIDs, sk.ID)
	}

	return s.repo.GetBySkillIDs(ctx, skillIDs, projectID, 1)
}
