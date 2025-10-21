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
}

type ProjectFetcher interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Project, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Project, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
