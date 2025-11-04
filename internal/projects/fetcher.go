package projects

import (
	"context"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
)

type ProjectFetcherAdapter struct {
	svc ProjectService
}

func NewProjectFetcherAdapter(svc ProjectService) *ProjectFetcherAdapter {
	return &ProjectFetcherAdapter{svc: svc}
}

func (a *ProjectFetcherAdapter) GetByID(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	return a.svc.GetByID(ctx, id)
}

func (a *ProjectFetcherAdapter) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Project, error) {
	return a.svc.GetByUserID(ctx, userID)
}

func (a *ProjectFetcherAdapter) Delete(ctx context.Context, id uuid.UUID) error {
	return a.svc.Delete(ctx, id)
}

func (a *ProjectFetcherAdapter) GetFeedForFreelancer(ctx context.Context, freelancerID, currentUserID uuid.UUID) (*domain.Project, error) {
	return a.svc.GetFeedForFreelancer(ctx, freelancerID, currentUserID)
}
