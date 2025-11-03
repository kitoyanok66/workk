package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
)

type LikeDTO struct {
	ID         uuid.UUID `json:"id"`
	FromUserID uuid.UUID `json:"from_user_id"`
	ToUserID   uuid.UUID `json:"to_user_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type LikeRequest struct {
	FromUserID uuid.UUID `json:"from_user_id"`
	ToUserID   uuid.UUID `json:"to_user_id"`
}

type LikeResponse struct {
	Like  *LikeDTO  `json:"like"`
	Match *MatchDTO `json:"match,omitempty"`
}

type NextFeedResponse struct {
	Next interface{} `json:"next"`
}

func NewLikeDTO(l *domain.Like) *LikeDTO {
	if l == nil {
		return nil
	}
	return &LikeDTO{
		ID:         l.ID,
		FromUserID: l.FromUserID,
		ToUserID:   l.ToUserID,
		CreatedAt:  l.CreatedAt,
	}
}

func (dto *LikeDTO) ToDomain() *domain.Like {
	if dto == nil {
		return nil
	}
	return &domain.Like{
		ID:         dto.ID,
		FromUserID: dto.FromUserID,
		ToUserID:   dto.ToUserID,
		CreatedAt:  dto.CreatedAt,
	}
}
