package web

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/dto"
	"github.com/kitoyanok66/workk/internal/skills"
	"github.com/kitoyanok66/workk/internal/web/oskills"
)

type skillHandler struct {
	svc skills.SkillService
}

func NewSkillHandler(svc skills.SkillService) *skillHandler {
	return &skillHandler{svc: svc}
}

// GET /skills
func (h *skillHandler) GetSkills(ctx context.Context, _ oskills.GetSkillsRequestObject) (oskills.GetSkillsResponseObject, error) {
	skillsList, err := h.svc.GetAll(ctx)
	if err != nil {
		return oskills.GetSkills500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, err.Error())), nil
	}

	dtos := make([]oskills.SkillDTO, 0, len(skillsList))
	for _, s := range skillsList {
		dtos = append(dtos, *dto.NewSkillDTO(s))
	}

	return oskills.GetSkills200JSONResponse(dtos), nil
}

// GET /skills/{id}
func (h *skillHandler) GetSkillsId(ctx context.Context, request oskills.GetSkillsIdRequestObject) (oskills.GetSkillsIdResponseObject, error) {
	id, err := uuid.Parse(request.Id.String())
	if err != nil {
		return oskills.GetSkillsId400JSONResponse(*dto.NewErrorDTO(http.StatusBadRequest, "invalid UUID")), nil
	}

	skill, err := h.svc.GetByID(ctx, id)
	if err != nil {
		return oskills.GetSkillsId500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, err.Error())), nil
	}
	if skill == nil {
		return oskills.GetSkillsId404JSONResponse(*dto.NewErrorDTO(http.StatusNotFound, "skill not found")), nil
	}

	return oskills.GetSkillsId200JSONResponse(*dto.NewSkillDTO(skill)), nil
}

// POST /skills
func (h *skillHandler) PostSkills(ctx context.Context, request oskills.PostSkillsRequestObject) (oskills.PostSkillsResponseObject, error) {
	body := request.Body

	createdSkill, err := h.svc.Create(ctx, body.Name, body.Category, body.Description)
	if err != nil {
		return oskills.PostSkills400JSONResponse(*dto.NewErrorDTO(http.StatusBadRequest, err.Error())), nil
	}

	return oskills.PostSkills201JSONResponse(*dto.NewSkillDTO(createdSkill)), nil
}

// PATCH /skills/{id}
func (h *skillHandler) PatchSkillsId(ctx context.Context, request oskills.PatchSkillsIdRequestObject) (oskills.PatchSkillsIdResponseObject, error) {
	id, err := uuid.Parse(request.Id.String())
	if err != nil {
		return oskills.PatchSkillsId400JSONResponse(*dto.NewErrorDTO(http.StatusBadRequest, "invalid UUID")), nil
	}

	body := request.Body

	updatedSkill, err := h.svc.Update(ctx, id, body.Name, body.Category, body.Description)
	if err != nil {
		return oskills.PatchSkillsId400JSONResponse(*dto.NewErrorDTO(http.StatusBadRequest, err.Error())), nil
	}

	return oskills.PatchSkillsId200JSONResponse(*dto.NewSkillDTO(updatedSkill)), nil
}

// DELETE /skills/{id}
func (h *skillHandler) DeleteSkillsId(ctx context.Context, request oskills.DeleteSkillsIdRequestObject) (oskills.DeleteSkillsIdResponseObject, error) {
	if err := h.svc.Delete(ctx, request.Id); err != nil {
		return oskills.DeleteSkillsId404JSONResponse(*dto.NewErrorDTO(http.StatusNotFound, err.Error())), nil
	}
	return oskills.DeleteSkillsId204Response{}, nil
}
