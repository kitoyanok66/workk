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
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func NewUser(telegramUserId int64, telegramUsername, fullname string) (*User, error) {
	if telegramUserId == 0 {
		return nil, errors.New("telegram user id cannot be empty")
	}
	if strings.TrimSpace(fullname) == "" {
		return nil, errors.New("full name cannot be empty")
	}

	return &User{
		ID:               uuid.New(),
		TelegramUserID:   telegramUserId,
		TelegramUsername: telegramUsername,
		Fullname:         fullname,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}, nil
}

func (u *User) ChangeFullName(newName string) error {
	newName = strings.TrimSpace(newName)
	if newName == "" {
		return errors.New("new name cannot be empty")
	}
	u.Fullname = newName
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) ChangeTelegramUsername(username string) {
	u.TelegramUsername = strings.TrimSpace(username)
	u.UpdatedAt = time.Now()
}
