package likes

import (
	"context"

	"github.com/google/uuid"
)

type LikeFetcherAdapter struct {
	svc LikeService
}

func NewLikeFetcherAdapter(svc LikeService) *LikeFetcherAdapter {
	return &LikeFetcherAdapter{svc: svc}
}

func (a *LikeFetcherAdapter) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return a.svc.DeleteByUserID(ctx, userID)
}
