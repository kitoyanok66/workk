package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Match struct {
	ID           uuid.UUID
	FreelancerID uuid.UUID
	ProjectID    uuid.UUID
	CreatedAt    time.Time
}

func NewMatch(freelancerID, projectID uuid.UUID) (*Match, error) {
	if freelancerID == uuid.Nil {
		return nil, errors.New("freelancer id cannot be empty")
	}
	if projectID == uuid.Nil {
		return nil, errors.New("project id cannot be empty")
	}

	return &Match{
		ID:           uuid.New(),
		FreelancerID: freelancerID,
		ProjectID:    projectID,
		CreatedAt:    time.Now(),
	}, nil
}
