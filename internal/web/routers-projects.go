package web

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/dto"
	"github.com/kitoyanok66/workk/internal/middleware"
	"github.com/kitoyanok66/workk/internal/projects"
	"github.com/kitoyanok66/workk/internal/web/oprojects"
)

type projectHandler struct {
	svc projects.ProjectService
}

func NewProjectHandler(svc projects.ProjectService) *projectHandler {
	return &projectHandler{svc: svc}
}

// GET /projects
func (h *projectHandler) GetProjects(ctx context.Context, _ oprojects.GetProjectsRequestObject) (oprojects.GetProjectsResponseObject, error) {
	projectsList, err := h.svc.GetAll(ctx)
	if err != nil {
		return oprojects.GetProjects500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, err.Error())), nil
	}

	dtos := make([]oprojects.ProjectDTO, 0, len(projectsList))
	for _, p := range projectsList {
		dtos = append(dtos, *dto.NewProjectDTO(p))
	}

	return oprojects.GetProjects200JSONResponse(dtos), nil
}

// GET /projects/{id}
func (h *projectHandler) GetProjectsId(ctx context.Context, request oprojects.GetProjectsIdRequestObject) (oprojects.GetProjectsIdResponseObject, error) {
	id, err := uuid.Parse(request.Id.String())
	if err != nil {
		return oprojects.GetProjectsId400JSONResponse(*dto.NewErrorDTO(http.StatusBadRequest, "invalid UUID")), nil
	}

	project, err := h.svc.GetByID(ctx, id)
	if err != nil {
		return oprojects.GetProjectsId500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, err.Error())), nil
	}
	if project == nil {
		return oprojects.GetProjectsId404JSONResponse(*dto.NewErrorDTO(http.StatusNotFound, "project not found")), nil
	}

	return oprojects.GetProjectsId200JSONResponse(*dto.NewProjectDTO(project)), nil
}

// GET /projects/by-user/{userID}
func (h *projectHandler) GetProjectsByUserUserID(ctx context.Context, request oprojects.GetProjectsByUserUserIDRequestObject) (oprojects.GetProjectsByUserUserIDResponseObject, error) {
	userID, err := uuid.Parse(request.UserID.String())
	if err != nil {
		return oprojects.GetProjectsByUserUserID400JSONResponse(*dto.NewErrorDTO(http.StatusBadRequest, "invalid UUID")), nil
	}

	project, err := h.svc.GetByUserID(ctx, userID)
	if err != nil {
		return oprojects.GetProjectsByUserUserID500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, err.Error())), nil
	}
	if project == nil {
		return oprojects.GetProjectsByUserUserID404JSONResponse(*dto.NewErrorDTO(http.StatusNotFound, "project not found")), nil
	}

	return oprojects.GetProjectsByUserUserID200JSONResponse(*dto.NewProjectDTO(project)), nil
}

// POST /projects
func (h *projectHandler) PostProjects(ctx context.Context, request oprojects.PostProjectsRequestObject) (oprojects.PostProjectsResponseObject, error) {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return oprojects.PostProjects401JSONResponse(*dto.NewErrorDTO(http.StatusUnauthorized, "unauthorized")), nil
	}

	body := request.Body

	createdProject, err := h.svc.Create(ctx, userID, body.Title, body.Description, body.Budget, body.Deadline, body.SkillIDs)
	if err != nil {
		return oprojects.PostProjects400JSONResponse(*dto.NewErrorDTO(http.StatusBadRequest, err.Error())), nil
	}

	return oprojects.PostProjects201JSONResponse(*dto.NewProjectDTO(createdProject)), nil
}

// PATCH /projects/{id}
func (h *projectHandler) PatchProjectsId(ctx context.Context, request oprojects.PatchProjectsIdRequestObject) (oprojects.PatchProjectsIdResponseObject, error) {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return oprojects.PatchProjectsId401JSONResponse(*dto.NewErrorDTO(http.StatusUnauthorized, "unauthorized")), nil
	}

	id, err := uuid.Parse(request.Id.String())
	if err != nil {
		return oprojects.PatchProjectsId400JSONResponse(*dto.NewErrorDTO(http.StatusBadRequest, "invalid UUID")), nil
	}

	project, err := h.svc.GetByID(ctx, id)
	if err != nil {
		return oprojects.PatchProjectsId500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, err.Error())), nil
	}
	if project == nil {
		return oprojects.PatchProjectsId404JSONResponse(*dto.NewErrorDTO(http.StatusNotFound, "freelancer not found")), nil
	}
	if project.UserID != userID {
		return oprojects.PatchProjectsId403JSONResponse(*dto.NewErrorDTO(http.StatusForbidden, "access denied")), nil
	}

	body := request.Body

	updatedProject, err := h.svc.Update(ctx, id, body.Title, body.Description, body.Budget, body.Deadline, body.SkillIDs)
	if err != nil {
		return oprojects.PatchProjectsId400JSONResponse(*dto.NewErrorDTO(http.StatusBadRequest, err.Error())), nil
	}

	return oprojects.PatchProjectsId200JSONResponse(*dto.NewProjectDTO(updatedProject)), nil

}

// DELETE /projects/{id}
func (h *projectHandler) DeleteProjectsId(ctx context.Context, request oprojects.DeleteProjectsIdRequestObject) (oprojects.DeleteProjectsIdResponseObject, error) {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return oprojects.DeleteProjectsId401JSONResponse(*dto.NewErrorDTO(http.StatusUnauthorized, "unauthorized")), nil
	}

	id, err := uuid.Parse(request.Id.String())
	if err != nil {
		return oprojects.DeleteProjectsId400JSONResponse(*dto.NewErrorDTO(http.StatusBadRequest, "invalid UUID")), nil
	}

	project, err := h.svc.GetByID(ctx, id)
	if err != nil {
		return oprojects.DeleteProjectsId500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, err.Error())), nil
	}
	if project == nil {
		return oprojects.DeleteProjectsId404JSONResponse(*dto.NewErrorDTO(http.StatusNotFound, "freelancer not found")), nil
	}
	if project.UserID != userID {
		return oprojects.DeleteProjectsId403JSONResponse(*dto.NewErrorDTO(http.StatusForbidden, "access denied")), nil
	}

	if err := h.svc.Delete(ctx, id); err != nil {
		return oprojects.DeleteProjectsId500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, err.Error())), nil
	}

	return oprojects.DeleteProjectsId204Response{}, nil
}
