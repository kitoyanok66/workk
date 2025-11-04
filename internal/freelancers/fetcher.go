package freelancers

import (
	"context"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
)

type FreelancerFetcherAdapter struct {
	svc FreelancerService
}

func NewFreelancerFetcherAdapter(svc FreelancerService) *FreelancerFetcherAdapter {
	return &FreelancerFetcherAdapter{svc: svc}
}

func (a *FreelancerFetcherAdapter) GetByID(ctx context.Context, id uuid.UUID) (*domain.Freelancer, error) {
	return a.svc.GetByID(ctx, id)
}

func (a *FreelancerFetcherAdapter) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Freelancer, error) {
	return a.svc.GetByUserID(ctx, userID)
}

func (a *FreelancerFetcherAdapter) Delete(ctx context.Context, id uuid.UUID) error {
	return a.svc.Delete(ctx, id)
}

func (a *FreelancerFetcherAdapter) GetFeedForProject(ctx context.Context, projectID, currentUserID uuid.UUID) (*domain.Freelancer, error) {
	return a.svc.GetFeedForProject(ctx, projectID, currentUserID)
}
