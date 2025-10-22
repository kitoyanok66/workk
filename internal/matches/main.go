package matches

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
)

type MatchService interface {
	GetAll(ctx context.Context) ([]*domain.Match, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Match, error)
	GetLastBetween(ctx context.Context, freelancerID, projectID uuid.UUID) (*domain.Match, error)
	Create(ctx context.Context, freelancerID, projectID uuid.UUID) (*domain.Match, error)
}

type matchService struct {
	repo MatchRepository
}

func NewMatchService(repo MatchRepository) MatchService {
	return &matchService{repo: repo}
}

func (s *matchService) GetAll(ctx context.Context) ([]*domain.Match, error) {
	return s.repo.GetAll(ctx)
}

func (s *matchService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Match, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *matchService) GetLastBetween(ctx context.Context, freelancerID, projectID uuid.UUID) (*domain.Match, error) {
	return s.repo.GetLastBetween(ctx, freelancerID, projectID)
}

func (s *matchService) Create(ctx context.Context, freelancerID, projectID uuid.UUID) (*domain.Match, error) {
	if freelancerID == uuid.Nil || projectID == uuid.Nil {
		return nil, errors.New("freelancerID and projectID must not be empty")
	}

	match, err := domain.NewMatch(freelancerID, projectID)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, match); err != nil {
		return nil, err
	}

	return match, nil
}
