package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Freelancer struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	Title           string
	Description     string
	HourlyRate      float64
	PortfolioURL    string
	ExperienceYears int
	Rating          float64
	Skills          []Skill
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewFreelancer(
	userID uuid.UUID,
	title string,
	description string,
	hourlyRate float64,
	portfolioURL string,
	experienceYears int,
	skills []Skill,
) (*Freelancer, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}
	if userID == uuid.Nil {
		return nil, errors.New("user ID cannot be empty")
	}
	if hourlyRate < 0 {
		return nil, errors.New("hourly rate cannot be negative")
	}
	if experienceYears < 0 {
		return nil, errors.New("experience years cannot be negative")
	}

	now := time.Now()
	return &Freelancer{
		ID:              uuid.New(),
		UserID:          userID,
		Title:           title,
		Description:     strings.TrimSpace(description),
		HourlyRate:      hourlyRate,
		PortfolioURL:    strings.TrimSpace(portfolioURL),
		ExperienceYears: experienceYears,
		Rating:          0,
		Skills:          skills,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

func (f *Freelancer) UpdateFreelancer(
	newTitle string,
	newDescription string,
	newHourlyRate float64,
	newPortfolioURL string,
	newExperienceYears int,
	newSkills []Skill,
) error {
	newTitle = strings.TrimSpace(newTitle)
	if newTitle == "" {
		return errors.New("title cannot be empty")
	}
	if newHourlyRate < 0 {
		return errors.New("hourly rate cannot be negative")
	}
	if newExperienceYears < 0 {
		return errors.New("experience years cannot be negative")
	}

	f.Title = newTitle
	f.Description = strings.TrimSpace(newDescription)
	f.HourlyRate = newHourlyRate
	f.PortfolioURL = strings.TrimSpace(newPortfolioURL)
	f.ExperienceYears = newExperienceYears
	f.Skills = newSkills
	f.UpdatedAt = time.Now()

	return nil
}

func (f *Freelancer) UpdateRating(newRating float64) error {
	if newRating < 0 || newRating > 5 {
		return errors.New("rating must be between 0 and 5")
	}
	f.Rating = newRating
	f.UpdatedAt = time.Now()
	return nil
}
