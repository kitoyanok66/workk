package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/kitoyanok66/workk/domain"
	"github.com/kitoyanok66/workk/internal/users"
	"gorm.io/gorm"
)

type AuthService interface {
	Authenticate(ctx context.Context, provider, externalID, username string) (*domain.Auth, error)
}

type authService struct {
	repo    AuthRepository
	userSvc users.UserService
}

func NewAuthService(repo AuthRepository, userSvc users.UserService) AuthService {
	return &authService{
		repo:    repo,
		userSvc: userSvc,
	}
}

func (s *authService) Authenticate(ctx context.Context, provider, externalID, username string) (*domain.Auth, error) {
	auth, err := s.repo.GetByProviderAndExternalID(ctx, provider, externalID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if auth != nil {
		return auth, nil
	}

	user, err := s.userSvc.Create(ctx, fmt.Sprintf("%s_user", username))
	if err != nil {
		return nil, err
	}

	auth, err = domain.NewAuth(user.ID, provider, externalID, username)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, auth); err != nil {
		return nil, err
	}

	return auth, nil
}
