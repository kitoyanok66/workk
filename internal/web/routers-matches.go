package web

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/dto"
	"github.com/kitoyanok66/workk/internal/freelancers"
	"github.com/kitoyanok66/workk/internal/matches"
	"github.com/kitoyanok66/workk/internal/projects"
	"github.com/kitoyanok66/workk/internal/users"
	"github.com/kitoyanok66/workk/internal/web/omatches"
)

type matchHandler struct {
	svc           matches.MatchService
	freelancerSvc freelancers.FreelancerService
	projectSvc    projects.ProjectService
	userSvc       users.UserService
}

func NewMatchHandler(svc matches.MatchService, freelancerSvc freelancers.FreelancerService, projectSvc projects.ProjectService, userSvc users.UserService) *matchHandler {
	return &matchHandler{
		svc:           svc,
		freelancerSvc: freelancerSvc,
		projectSvc:    projectSvc,
		userSvc:       userSvc,
	}
}

// GET /matches
func (h *matchHandler) GetMatches(ctx context.Context, _ omatches.GetMatchesRequestObject) (omatches.GetMatchesResponseObject, error) {
	matchesList, err := h.svc.GetAll(ctx)
	if err != nil {
		return omatches.GetMatches500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, err.Error())), nil
	}

	dtos := make([]omatches.MatchDTO, 0, len(matchesList))
	for _, m := range matchesList {
		freelancer, err := h.freelancerSvc.GetByID(ctx, m.FreelancerID)
		if err != nil || freelancer == nil {
			return omatches.GetMatches500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, "failed to fetch freelancer")), nil
		}

		project, err := h.projectSvc.GetByID(ctx, m.ProjectID)
		if err != nil || project == nil {
			return omatches.GetMatches500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, "failed to fetch project")), nil
		}

		freelancerUser, err := h.userSvc.GetByID(ctx, freelancer.UserID)
		if err != nil || freelancerUser == nil {
			return omatches.GetMatches500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, "failed to fetch freelancer user")), nil
		}

		projectUser, err := h.userSvc.GetByID(ctx, project.UserID)
		if err != nil || projectUser == nil {
			return omatches.GetMatches500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, "failed to fetch project user")), nil
		}

		dtos = append(dtos, *dto.NewMatchDTO(m, freelancer, project, freelancerUser, projectUser))
	}

	return omatches.GetMatches200JSONResponse(dtos), nil
}

// GET /matches/{id}
func (h *matchHandler) GetMatchesId(ctx context.Context, request omatches.GetMatchesIdRequestObject) (omatches.GetMatchesIdResponseObject, error) {
	id, err := uuid.Parse(request.Id.String())
	if err != nil {
		return omatches.GetMatchesId400JSONResponse(*dto.NewErrorDTO(http.StatusBadRequest, "invalid UUID")), nil
	}

	match, err := h.svc.GetByID(ctx, id)
	if err != nil {
		return omatches.GetMatchesId500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, err.Error())), nil
	}
	if match == nil {
		return omatches.GetMatchesId404JSONResponse(*dto.NewErrorDTO(http.StatusNotFound, "match not found")), nil
	}

	freelancer, err := h.freelancerSvc.GetByID(ctx, match.FreelancerID)
	if err != nil || freelancer == nil {
		return omatches.GetMatchesId500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, "failed to fetch freelancer")), nil
	}
	project, err := h.projectSvc.GetByID(ctx, match.ProjectID)
	if err != nil || project == nil {
		return omatches.GetMatchesId500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, "failed to fetch project")), nil
	}
	freelancerUser, err := h.userSvc.GetByID(ctx, freelancer.UserID)
	if err != nil || freelancerUser == nil {
		return omatches.GetMatchesId500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, "failed to fetch freelancer user")), nil
	}
	projectUser, err := h.userSvc.GetByID(ctx, project.UserID)
	if err != nil || projectUser == nil {
		return omatches.GetMatchesId500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, "failed to fetch project user")), nil
	}

	return omatches.GetMatchesId200JSONResponse(*dto.NewMatchDTO(match, freelancer, project, freelancerUser, projectUser)), nil
}
