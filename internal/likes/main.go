package likes

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
	"github.com/kitoyanok66/workk/internal/freelancers"
	"github.com/kitoyanok66/workk/internal/matches"
	"github.com/kitoyanok66/workk/internal/projects"
	"github.com/kitoyanok66/workk/internal/users"
)

type LikeService interface {
	GetFeed(ctx context.Context, userID uuid.UUID) (interface{}, error)
	Like(ctx context.Context, fromUserID, toUserID uuid.UUID) (*domain.Like, error)
	Dislike(ctx context.Context, fromUserID, toUserID uuid.UUID) error
}

type likeService struct {
	repo          LikeRepository
	userSvc       users.UserService
	freelancerSvc freelancers.FreelancerService
	projectSvc    projects.ProjectService
	matchSvc      matches.MatchService
}

func NewLikeService(repo LikeRepository, userSvc users.UserService, freelancerSvc freelancers.FreelancerService, projectSvc projects.ProjectService, matchSvc matches.MatchService) LikeService {
	return &likeService{
		repo:          repo,
		userSvc:       userSvc,
		freelancerSvc: freelancerSvc,
		projectSvc:    projectSvc,
		matchSvc:      matchSvc,
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

func (s *likeService) Like(ctx context.Context, fromUserID, toUserID uuid.UUID) (*domain.Like, error) {
	if fromUserID == uuid.Nil || toUserID == uuid.Nil {
		return nil, errors.New("user IDs must not be empty")
	}

	fromUser, err := s.userSvc.GetByID(ctx, fromUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get fromUser: %w", err)
	}
	if fromUser == nil {
		return nil, fmt.Errorf("user not found: %s", fromUserID)
	}

	toUser, err := s.userSvc.GetByID(ctx, toUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get toUser: %w", err)
	}
	if toUser == nil {
		return nil, fmt.Errorf("user not found: %s", toUserID)
	}

	like, err := domain.NewLike(fromUserID, toUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to create like: %w", err)
	}
	if err := s.repo.Create(ctx, like); err != nil {
		return nil, fmt.Errorf("failed to save like: %w", err)
	}

	exists, err := s.repo.ExistsReverse(ctx, fromUserID, toUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to check reverse like: %w", err)
	}
	if !exists {
		return like, nil
	}

	var freelancerID, projectID uuid.UUID

	switch fromUser.Role {
	case "freelancer":
		freelancer, err := s.freelancerSvc.GetByUserID(ctx, fromUserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get freelancer by user id: %w", err)
		}
		project, err := s.projectSvc.GetByUserID(ctx, toUserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get project by user id: %w", err)
		}
		freelancerID = freelancer.ID
		projectID = project.ID

	case "project":
		project, err := s.projectSvc.GetByUserID(ctx, fromUserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get project by user id: %w", err)
		}
		freelancer, err := s.freelancerSvc.GetByUserID(ctx, toUserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get freelancer by user id: %w", err)
		}
		freelancerID = freelancer.ID
		projectID = project.ID

	default:
		return nil, fmt.Errorf("unsupported role for matching: %s", fromUser.Role)
	}

	if _, err := s.matchSvc.Create(ctx, freelancerID, projectID); err != nil {
		return nil, fmt.Errorf("failed to create match: %w", err)
	}

	return like, nil
}

func (s *likeService) Dislike(ctx context.Context, fromUserID, toUserID uuid.UUID) error {
	return nil
}
