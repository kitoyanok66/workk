package users

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/domain"
)

type UserService interface {
	GetAll(ctx context.Context) ([]*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Create(ctx context.Context, telegramID int64, username, fullname string) (*domain.User, error)
	Update(ctx context.Context, id uuid.UUID, fullname, username, role string) (*domain.User, error)
	UpdateRole(ctx context.Context, id uuid.UUID, role string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetAll(ctx context.Context) ([]*domain.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) Create(ctx context.Context, telegramID int64, username, fullname string) (*domain.User, error) {
	existing, err := s.repo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("user with telegram ID %d already exists", telegramID)
	}

	user, err := domain.NewUser(telegramID, username, fullname)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Update(ctx context.Context, id uuid.UUID, fullname, username, role string) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	if fullname != "" {
		if err := user.UpdateFullName(fullname); err != nil {
			return nil, err
		}
	}
	if username != "" {
		user.UpdateTelegramUsername(username)
	}
	if role != "" {
		if err := user.UpdateRole(role); err != nil {
			return nil, err
		}
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateRole(ctx context.Context, id uuid.UUID, role string) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	if role != "" {
		if err := user.UpdateRole(role); err != nil {
			return err
		}
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
