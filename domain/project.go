package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Title       string
	Description string
	Budget      float64
	Deadline    time.Time
	Skills      []Skill
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewProject(
	userID uuid.UUID,
	title string,
	description string,
	budget float64,
	deadline time.Time,
	skills []Skill,
) (*Project, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}
	if userID == uuid.Nil {
		return nil, errors.New("user ID cannot be empty")
	}

	now := time.Now()
	return &Project{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       title,
		Description: description,
		Budget:      budget,
		Deadline:    deadline,
		Skills:      skills,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (p *Project) UpdateProject(
	newTitle string,
	newDescription string,
	newBudget float64,
	newDeadline time.Time,
	newSkills []Skill,
) error {
	newTitle = strings.TrimSpace(newTitle)
	if newTitle == "" {
		return errors.New("title cannot be empty")
	}

	p.Title = newTitle
	p.Description = newDescription
	p.Budget = newBudget
	p.Deadline = newDeadline
	p.Skills = newSkills
	p.UpdatedAt = time.Now()

	return nil
}
