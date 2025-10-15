package web

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/dto"
	"github.com/kitoyanok66/workk/internal/freelancers"
	"github.com/kitoyanok66/workk/internal/web/ofreelancers"
)

type freelancerHandler struct {
	svc freelancers.FreelancerService
}

func NewFreelancerHandler(svc freelancers.FreelancerService) *freelancerHandler {
	return &freelancerHandler{svc: svc}
}

// GET /freelancers
func (h *freelancerHandler) GetFreelancers(ctx context.Context, _ ofreelancers.GetFreelancersRequestObject) (ofreelancers.GetFreelancersResponseObject, error) {
	freelancersList, err := h.svc.GetAll(ctx)
	if err != nil {
		return ofreelancers.GetFreelancers500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, err.Error())), nil
	}

	dtos := make([]ofreelancers.FreelancerDTO, 0, len(freelancersList))
	for _, f := range freelancersList {
		dtos = append(dtos, *dto.NewFreelancerDTO(f))
	}

	return ofreelancers.GetFreelancers200JSONResponse(dtos), nil
}

// GET /freelancers/{id}
func (h *freelancerHandler) GetFreelancersId(ctx context.Context, request ofreelancers.GetFreelancersIdRequestObject) (ofreelancers.GetFreelancersIdResponseObject, error) {
	id, err := uuid.Parse(request.Id.String())
	if err != nil {
		return ofreelancers.GetFreelancersId400JSONResponse(*dto.NewErrorDTO(http.StatusBadRequest, "invalid UUID")), nil
	}

	freelancer, err := h.svc.GetByID(ctx, id)
	if err != nil {
		return ofreelancers.GetFreelancersId500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, err.Error())), nil
	}
	if freelancer == nil {
		return ofreelancers.GetFreelancersId404JSONResponse(*dto.NewErrorDTO(http.StatusNotFound, "freelancer not found")), nil
	}

	return ofreelancers.GetFreelancersId200JSONResponse(*dto.NewFreelancerDTO(freelancer)), nil
}

// POST /freelancers
func (h *freelancerHandler) PostFreelancers(ctx context.Context, request ofreelancers.PostFreelancersRequestObject) (ofreelancers.PostFreelancersResponseObject, error) {
	body := request.Body

	createdFreelancer, err := h.svc.Create(ctx, body.UserID, body.Title, body.Description, body.HourlyRate, body.PortfolioURL, body.ExperienceYears, body.SkillIDs)
	if err != nil {
		return ofreelancers.PostFreelancers400JSONResponse(*dto.NewErrorDTO(http.StatusBadRequest, err.Error())), nil
	}

	return ofreelancers.PostFreelancers201JSONResponse(*dto.NewFreelancerDTO(createdFreelancer)), nil
}

// PATCH /freelancers/{id}
func (h *freelancerHandler) PatchFreelancersId(ctx context.Context, request ofreelancers.PatchFreelancersIdRequestObject) (ofreelancers.PatchFreelancersIdResponseObject, error) {
	id, err := uuid.Parse(request.Id.String())
	if err != nil {
		return ofreelancers.PatchFreelancersId400JSONResponse(*dto.NewErrorDTO(http.StatusBadRequest, "invalid UUID")), nil
	}

	body := request.Body

	updatedFreelancer, err := h.svc.Update(ctx, id, body.Title, body.Description, body.HourlyRate, body.PortfolioURL, body.ExperienceYears, body.SkillIDs)
	if err != nil {
		return ofreelancers.PatchFreelancersId400JSONResponse(*dto.NewErrorDTO(http.StatusBadRequest, err.Error())), nil
	}

	return ofreelancers.PatchFreelancersId200JSONResponse(*dto.NewFreelancerDTO(updatedFreelancer)), nil
}

// DELETE /freelancers/{id}
func (h *freelancerHandler) DeleteFreelancersId(ctx context.Context, request ofreelancers.DeleteFreelancersIdRequestObject) (ofreelancers.DeleteFreelancersIdResponseObject, error) {
	if err := h.svc.Delete(ctx, request.Id); err != nil {
		return ofreelancers.DeleteFreelancersId404JSONResponse(*dto.NewErrorDTO(http.StatusNotFound, err.Error())), nil
	}
	return ofreelancers.DeleteFreelancersId204Response{}, nil
}
