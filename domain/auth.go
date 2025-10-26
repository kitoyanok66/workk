package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Auth struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Provider   string
	ExternalID string
	Username   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewAuth(userID uuid.UUID, provider, externalID, username string) (*Auth, error) {
	if userID == uuid.Nil {
		return nil, errors.New("userID cannot be nil")
	}
	provider = strings.TrimSpace(provider)
	if provider == "" {
		return nil, errors.New("provider cannot be empty")
	}
	externalID = strings.TrimSpace(externalID)
	if externalID == "" {
		return nil, errors.New("externalID cannot be empty")
	}

	return &Auth{
		ID:         uuid.New(),
		UserID:     userID,
		Provider:   provider,
		ExternalID: externalID,
		Username:   strings.TrimSpace(username),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}

func (a *Auth) UpdateUsername(newUsername string) {
	a.Username = strings.TrimSpace(newUsername)
	a.UpdatedAt = time.Now()
}
