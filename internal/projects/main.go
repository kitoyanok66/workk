package projects

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
	"github.com/kitoyanok66/workk/internal/deps"
	"github.com/kitoyanok66/workk/internal/skills"
	"github.com/kitoyanok66/workk/internal/users"
)

type ProjectService interface {
	GetAll(ctx context.Context) ([]*domain.Project, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Project, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Project, error)
	Create(ctx context.Context, userID uuid.UUID, title, description string, budget float64, deadline time.Time, skillIDs []uuid.UUID) (*domain.Project, error)
	Update(ctx context.Context, id uuid.UUID, title, description string, budget float64, deadline time.Time, skillIDs []uuid.UUID) (*domain.Project, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetFeedForFreelancer(ctx context.Context, freelancerID uuid.UUID) (*domain.Project, error)

	SetFreelancerDep(dep deps.FreelancerFetcher)
}

type projectService struct {
	repo          ProjectRepository
	skillSvc      skills.SkillService
	freelancerDep deps.FreelancerFetcher
	userSvc       users.UserService
}

func NewProjectService(repo ProjectRepository, skillSvc skills.SkillService, userSvc users.UserService) ProjectService {
	return &projectService{
		repo:     repo,
		skillSvc: skillSvc,
		userSvc:  userSvc,
	}
}

func (s *projectService) SetFreelancerDep(dep deps.FreelancerFetcher) {
	s.freelancerDep = dep
}

func (s *projectService) GetAll(ctx context.Context) ([]*domain.Project, error) {
	return s.repo.GetAll(ctx)
}

func (s *projectService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *projectService) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Project, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *projectService) Create(ctx context.Context, userID uuid.UUID, title, description string, budget float64, deadline time.Time, skillIDs []uuid.UUID) (*domain.Project, error) {
	existing, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("project for this user already exists")
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

	freelancer, err := s.freelancerDep.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if freelancer != nil {
		if err := s.freelancerDep.Delete(ctx, freelancer.ID); err != nil {
			return nil, err
		}
	}

	project, err := domain.NewProject(userID, title, description, budget, deadline, skillsDomain)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, project); err != nil {
		return nil, err
	}

	if err := s.userSvc.UpdateRole(ctx, userID, "project"); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *projectService) Update(ctx context.Context, id uuid.UUID, title, description string, budget float64, deadline time.Time, skillIDs []uuid.UUID) (*domain.Project, error) {
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, fmt.Errorf("project not found")
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

	if err := project.UpdateProject(title, description, budget, deadline, skillsDomain); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, project); err != nil {
		return nil, err
	}

	return project, err
}

func (s *projectService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *projectService) GetFeedForFreelancer(ctx context.Context, freelancerID uuid.UUID) (*domain.Project, error) {
	freelancer, err := s.freelancerDep.GetByID(ctx, freelancerID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetBySkillIDs(freelancer.SkillIDs, 1)
}
