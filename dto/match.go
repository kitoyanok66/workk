package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
)

type MatchDTO struct {
	ID             uuid.UUID      `json:"id"`
	Freelancer     *FreelancerDTO `json:"freelancer"`
	Project        *ProjectDTO    `json:"project"`
	FreelancerUser *UserDTO       `json:"freelancer_user"`
	ProjectUser    *UserDTO       `json:"project_user"`
	CreatedAt      time.Time      `json:"created_at"`
}

func NewMatchDTO(m *domain.Match, freelancer *domain.Freelancer, project *domain.Project, freelancerUser *domain.User, projectUser *domain.User) *MatchDTO {
	if m == nil {
		return nil
	}
	return &MatchDTO{
		ID:             m.ID,
		Freelancer:     NewFreelancerDTO(freelancer),
		Project:        NewProjectDTO(project),
		FreelancerUser: NewUserDTO(freelancerUser),
		ProjectUser:    NewUserDTO(projectUser),
		CreatedAt:      m.CreatedAt,
	}
}

func (dto *MatchDTO) ToDomain() *domain.Match {
	if dto == nil {
		return nil
	}

	var freelancerID, projectID uuid.UUID
	if dto.Freelancer != nil {
		freelancerID = dto.Freelancer.ID
	}
	if dto.Project != nil {
		projectID = dto.Project.ID
	}

	return &domain.Match{
		ID:           dto.ID,
		FreelancerID: freelancerID,
		ProjectID:    projectID,
		CreatedAt:    dto.CreatedAt,
	}
}
