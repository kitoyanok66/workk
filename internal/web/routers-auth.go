package web

import (
	"context"
	"net/http"

	"github.com/kitoyanok66/workk/dto"
	"github.com/kitoyanok66/workk/internal/auth"
	"github.com/kitoyanok66/workk/internal/middleware"
	"github.com/kitoyanok66/workk/internal/users"
	"github.com/kitoyanok66/workk/internal/web/oauth"
)

type authHandler struct {
	svc        auth.AuthService
	userSvc    users.UserService
	jwtManager *middleware.JWTManager
}

func NewAuthHandler(svc auth.AuthService, userSvc users.UserService, jwtManager *middleware.JWTManager) *authHandler {
	return &authHandler{
		svc:        svc,
		userSvc:    userSvc,
		jwtManager: jwtManager,
	}
}

// POST /auth
func (h *authHandler) PostAuth(ctx context.Context, request oauth.PostAuthRequestObject) (oauth.PostAuthResponseObject, error) {
	body := request.Body

	auth, err := h.svc.Authenticate(ctx, body.Provider, body.ExternalID, body.Username)
	if err != nil {
		return oauth.PostAuth400JSONResponse(*dto.NewErrorDTO(http.StatusBadRequest, err.Error())), nil
	}

	user, err := h.userSvc.GetByID(ctx, auth.UserID)
	if err != nil {
		return oauth.PostAuth500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, err.Error())), nil
	}

	token, err := h.jwtManager.GenerateToken(user.ID.String())
	if err != nil {
		return oauth.PostAuth500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, "failed to generate token")), nil
	}

	resp := dto.AuthResponse{
		Token: token,
		User:  dto.NewUserDTO(user),
	}

	return oauth.PostAuth200JSONResponse(resp), nil
}
