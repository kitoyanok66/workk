package deps

import (
	"context"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
)

type FreelancerFetcher interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Freelancer, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Freelancer, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetFeedForProject(ctx context.Context, projectID, currentUserID uuid.UUID) (*domain.Freelancer, error)
}

type ProjectFetcher interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Project, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Project, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetFeedForFreelancer(ctx context.Context, freelancerID, currentUserID uuid.UUID) (*domain.Project, error)
}

type LikeFetcher interface {
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}
