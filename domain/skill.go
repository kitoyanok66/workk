package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Skill struct {
	ID          uuid.UUID
	Name        string
	Category    string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewSkill(name, category, description string) (*Skill, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("name cannot be empty")
	}
	if strings.TrimSpace(category) == "" {
		return nil, errors.New("category cannot be empty")
	}

	now := time.Now()
	return &Skill{
		ID:          uuid.New(),
		Name:        name,
		Category:    category,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (s *Skill) UpdateSkill(newName, newCategory, newDescription string) error {
	newName = strings.TrimSpace(newName)
	if newName == "" {
		return errors.New("new name cannot be empty")
	}
	newCategory = strings.TrimSpace(newCategory)
	if newCategory == "" {
		return errors.New("new category cannot be empty")
	}

	s.Name = newName
	s.Category = newCategory
	s.Description = newDescription
	s.UpdatedAt = time.Now()
	return nil
}

func (s *Skill) ChangeName(newName string) error {
	newName = strings.TrimSpace(newName)
	if newName == "" {
		return errors.New("new name cannot be empty")
	}
	s.Name = newName
	s.UpdatedAt = time.Now()
	return nil
}

func (s *Skill) ChangeCategory(newCategory string) error {
	newCategory = strings.TrimSpace(newCategory)
	if newCategory == "" {
		return errors.New("new category cannot be empty")
	}
	s.Category = newCategory
	s.UpdatedAt = time.Now()
	return nil
}

func (s *Skill) ChangeDescription(newDescription string) error {
	newDescription = strings.TrimSpace(newDescription)
	if newDescription == "" {
		return errors.New("new description cannot be empty")
	}
	s.Category = newDescription
	s.UpdatedAt = time.Now()
	return nil
}
