package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Fullname  string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(fullname string) (*User, error) {
	if strings.TrimSpace(fullname) == "" {
		return nil, errors.New("full name cannot be empty")
	}

	id := uuid.New()

	now := time.Now()
	return &User{
		ID:        id,
		Fullname:  fullname,
		Role:      "undefined",
		CreatedAt: now,
		UpdatedAt: now,
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

func (u *User) UpdateRole(role string) error {
	role = strings.TrimSpace(role)
	if role != "freelancer" && role != "project" && role != "undefined" {
		return errors.New("role must be 'freelancer', 'project' or 'undefined'")
	}
	u.Role = role
	u.UpdatedAt = time.Now()
	return nil
}
