package likes

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
	"github.com/kitoyanok66/workk/internal/freelancers"
	"github.com/kitoyanok66/workk/internal/projects"
	"github.com/kitoyanok66/workk/internal/users"
)

type LikeService interface {
	GetFeed(ctx context.Context, userID uuid.UUID) (interface{}, error)
	Like(ctx context.Context, userID, targetID uuid.UUID) (*domain.Like, error)
	Dislike(ctx context.Context, userID, targetID uuid.UUID) error
}

type likeService struct {
	repo          LikeRepository
	userSvc       users.UserService
	freelancerSvc freelancers.FreelancerService
	projectSvc    projects.ProjectService
}

func NewLikeService(repo LikeRepository, userSvc users.UserService, freelancerSvc freelancers.FreelancerService, projectSvc projects.ProjectService) LikeService {
	return &likeService{
		repo:          repo,
		userSvc:       userSvc,
		freelancerSvc: freelancerSvc,
		projectSvc:    projectSvc,
	}
}

func (s *likeService) GetFeed(ctx context.Context, userID uuid.UUID) (interface{}, error) {
	user, err := s.userSvc.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	switch user.Role {
	case "project":
		project, err := s.projectSvc.GetByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
		return s.freelancerSvc.GetFeedForProject(ctx, project.ID)

	case "freelancer":
		freelancer, err := s.freelancerSvc.GetByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
		return s.projectSvc.GetFeedForFreelancer(ctx, freelancer.ID)

	default:
		return nil, fmt.Errorf("unsupported user role: %s", user.Role)
	}
}

func (s *likeService) Like(ctx context.Context, userID, targetID uuid.UUID) (*domain.Like, error) {
	if userID == uuid.Nil {
		return nil, errors.New("invalid user id")
	}
	if targetID == uuid.Nil {
		return nil, errors.New("invalid target id")
	}

	like, err := domain.NewLike(userID, targetID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, like); err != nil {
		return nil, err
	}

	exists, err := s.repo.ExistsReverse(ctx, userID, targetID)
	if err != nil {
		return nil, err
	}
	if exists {
		// TODO: интеграция с сервисом matchService.Create
	}

	return like, nil
}

func (s *likeService) Dislike(ctx context.Context, userID, targetID uuid.UUID) error {
	return nil
}
