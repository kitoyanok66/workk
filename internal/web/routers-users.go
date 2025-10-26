package web

import (
	"context"
	"net/http"

	"github.com/kitoyanok66/workk/dto"
	"github.com/kitoyanok66/workk/internal/middleware"
	"github.com/kitoyanok66/workk/internal/users"
	"github.com/kitoyanok66/workk/internal/web/ousers"
)

type userHandler struct {
	svc users.UserService
}

func NewUserHandler(svc users.UserService) *userHandler {
	return &userHandler{svc: svc}
}

// GET /users
func (h *userHandler) GetUsers(ctx context.Context, _ ousers.GetUsersRequestObject) (ousers.GetUsersResponseObject, error) {
	usersList, err := h.svc.GetAll(ctx)
	if err != nil {
		return ousers.GetUsers500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, err.Error())), nil
	}

	dtos := make([]ousers.UserDTO, 0, len(usersList))
	for _, u := range usersList {
		dtos = append(dtos, *dto.NewUserDTO(u))
	}

	return ousers.GetUsers200JSONResponse(dtos), nil
}

// GET /users/{id}
func (h *userHandler) GetUsersId(ctx context.Context, request ousers.GetUsersIdRequestObject) (ousers.GetUsersIdResponseObject, error) {
	user, err := h.svc.GetByID(ctx, request.Id)
	if err != nil {
		return ousers.GetUsersId500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, err.Error())), nil
	}
	if user == nil {
		return ousers.GetUsersId404JSONResponse(*dto.NewErrorDTO(http.StatusNotFound, "user not found")), nil
	}

	return ousers.GetUsersId200JSONResponse(*dto.NewUserDTO(user)), nil
}

// POST /users
func (h *userHandler) PostUsers(ctx context.Context, request ousers.PostUsersRequestObject) (ousers.PostUsersResponseObject, error) {
	body := request.Body

	createdUser, err := h.svc.Create(ctx, body.Fullname)
	if err != nil {
		return ousers.PostUsers400JSONResponse(*dto.NewErrorDTO(http.StatusBadRequest, err.Error())), nil
	}

	return ousers.PostUsers201JSONResponse(*dto.NewUserDTO(createdUser)), nil
}

// PATCH /users/{id}
func (h *userHandler) PatchUsersId(ctx context.Context, request ousers.PatchUsersIdRequestObject) (ousers.PatchUsersIdResponseObject, error) {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return ousers.PatchUsersId401JSONResponse(*dto.NewErrorDTO(http.StatusUnauthorized, "unauthorized")), nil
	}

	id := request.Id

	if userID != id {
		return ousers.PatchUsersId403JSONResponse(*dto.NewErrorDTO(http.StatusForbidden, "forbidden")), nil
	}

	body := request.Body

	updatedUser, err := h.svc.Update(ctx, id, body.Fullname, body.Role)
	if err != nil {
		return ousers.PatchUsersId400JSONResponse(*dto.NewErrorDTO(http.StatusBadRequest, err.Error())), nil
	}

	return ousers.PatchUsersId200JSONResponse(*dto.NewUserDTO(updatedUser)), nil
}

// DELETE /users/{id}
func (h *userHandler) DeleteUsersId(ctx context.Context, request ousers.DeleteUsersIdRequestObject) (ousers.DeleteUsersIdResponseObject, error) {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return ousers.DeleteUsersId401JSONResponse(*dto.NewErrorDTO(http.StatusUnauthorized, "unauthorized")), nil
	}

	id := request.Id

	if userID != id {
		return ousers.DeleteUsersId403JSONResponse(*dto.NewErrorDTO(http.StatusForbidden, "forbidden")), nil
	}

	if err := h.svc.Delete(ctx, request.Id); err != nil {
		return ousers.DeleteUsersId404JSONResponse(*dto.NewErrorDTO(http.StatusNotFound, err.Error())), nil
	}

	return ousers.DeleteUsersId204Response{}, nil
}
