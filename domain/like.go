package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Like struct {
	ID         uuid.UUID
	FromUserID uuid.UUID
	ToUserID   uuid.UUID
	CreatedAt  time.Time
}

func NewLike(from, to uuid.UUID) (*Like, error) {
	if from == uuid.Nil {
		return nil, errors.New("user id cannot be empty")
	}
	if to == uuid.Nil {
		return nil, errors.New("user id cannot be empty")
	}

	return &Like{
		ID:         uuid.New(),
		FromUserID: from,
		ToUserID:   to,
		CreatedAt:  time.Now(),
	}, nil
}
