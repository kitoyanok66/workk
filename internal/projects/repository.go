package projects

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	GetAll(ctx context.Context) ([]*domain.Project, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Project, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Project, error)
	Create(ctx context.Context, project *domain.Project) error
	Update(ctx context.Context, project *domain.Project) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetBySkillIDs(ctx context.Context, skillIDs []uuid.UUID, currentUserID uuid.UUID) (*domain.Project, error)
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) GetAll(ctx context.Context) ([]*domain.Project, error) {
	var ormProjects []ProjectORM
	if err := r.db.WithContext(ctx).Preload("Skills").Find(&ormProjects).Error; err != nil {
		return nil, err
	}
	result := make([]*domain.Project, len(ormProjects))
	for i, p := range ormProjects {
		result[i] = p.ToDomain()
	}
	return result, nil
}

func (r *projectRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	var ormProject ProjectORM
	if err := r.db.WithContext(ctx).Preload("Skills").First(&ormProject, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return ormProject.ToDomain(), nil
}

func (r *projectRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Project, error) {
	var ormProject ProjectORM
	if err := r.db.WithContext(ctx).Preload("Skills").First(&ormProject, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return ormProject.ToDomain(), nil
}

func (r *projectRepository) Create(ctx context.Context, project *domain.Project) error {
	if project == nil {
		return errors.New("project is nil")
	}
	ormProject := FromDomain(project)
	return r.db.WithContext(ctx).Create(&ormProject).Error
}

func (r *projectRepository) Update(ctx context.Context, project *domain.Project) error {
	if project == nil {
		return errors.New("project is nil")
	}
	ormProject := FromDomain(project)
	tx := r.db.WithContext(ctx).Model(&ProjectORM{}).Where("id = ?", project.ID).Updates(&ormProject)
	if tx.Error != nil {
		return tx.Error
	}
	if err := r.db.WithContext(ctx).Model(&ormProject).Association("Skills").Replace(ormProject.Skills); err != nil {
		return err
	}
	return nil
}

func (r *projectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&ProjectORM{}, "id = ?", id).Error
}

func (r *projectRepository) GetBySkillIDs(ctx context.Context, skillIDs []uuid.UUID, currentUserID uuid.UUID) (*domain.Project, error) {
	if len(skillIDs) == 0 {
		return nil, nil
	}

	query := `
        SELECT p.*
        FROM projects p
        JOIN project_skills ps ON p.id = ps.project_id
        WHERE ps.skill_id IN ?
        AND p.user_id NOT IN (
            SELECT to_user_id
            FROM likes
            WHERE from_user_id = ?
        )
        GROUP BY p.id
        ORDER BY COUNT(DISTINCT ps.skill_id) DESC
        LIMIT 1;
    `

	var ormProject ProjectORM
	if err := r.db.WithContext(ctx).Raw(query, skillIDs, currentUserID).Scan(&ormProject).Error; err != nil {
		return nil, err
	}

	if ormProject.ID != uuid.Nil {
		if err := r.db.WithContext(ctx).Preload("Skills").First(&ormProject, "id = ?", ormProject.ID).Error; err != nil {
			return nil, err
		}
	}

	return ormProject.ToDomain(), nil
}
