package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID               uuid.UUID
	TelegramUserID   int64
	TelegramUsername string
	Fullname         string
	Role             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func NewUser(telegramUserID int64, telegramUsername, fullname string) (*User, error) {
	if telegramUserID == 0 {
		return nil, errors.New("telegram user id cannot be empty")
	}
	if strings.TrimSpace(fullname) == "" {
		return nil, errors.New("full name cannot be empty")
	}

	id := uuid.New()

	if strings.TrimSpace(telegramUsername) == "" {
		telegramUsername = "user_" + id.String()[:8]
	}

	now := time.Now()
	return &User{
		ID:               id,
		TelegramUserID:   telegramUserID,
		TelegramUsername: telegramUsername,
		Fullname:         fullname,
		Role:             "undefind",
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

func (u *User) UpdateFullName(newName string) error {
	newName = strings.TrimSpace(newName)
	if newName == "" {
		return errors.New("new name cannot be empty")
	}
	u.Fullname = newName
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) UpdateTelegramUsername(username string) {
	u.TelegramUsername = strings.TrimSpace(username)
	u.UpdatedAt = time.Now()
}

func (u *User) UpdateRole(role string) error {
	role = strings.TrimSpace(role)
	if role != "freelancer" && role != "project" && role != "undefind" {
		return errors.New("role must be 'freelancer', 'project' or 'undefind'")
	}
	u.Role = role
	u.UpdatedAt = time.Now()
	return nil
}
