package web

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/kitoyanok66/workk/dto"
	"github.com/kitoyanok66/workk/internal/freelancers"
	"github.com/kitoyanok66/workk/internal/likes"
	"github.com/kitoyanok66/workk/internal/matches"
	"github.com/kitoyanok66/workk/internal/projects"
	"github.com/kitoyanok66/workk/internal/users"
	"github.com/kitoyanok66/workk/internal/web/olikes"
)

type likeHandler struct {
	svc           likes.LikeService
	matchSvc      matches.MatchService
	freelancerSvc freelancers.FreelancerService
	projectSvc    projects.ProjectService
	userSvc       users.UserService
}

func NewLikeHandler(svc likes.LikeService, matchSvc matches.MatchService, freelancerSvc freelancers.FreelancerService, projectSvc projects.ProjectService, userSvc users.UserService) *likeHandler {
	return &likeHandler{
		svc:           svc,
		matchSvc:      matchSvc,
		freelancerSvc: freelancerSvc,
		projectSvc:    projectSvc,
		userSvc:       userSvc,
	}
}

// POST /likes/like
func (h *likeHandler) PostLikesLike(ctx context.Context, request olikes.PostLikesLikeRequestObject) (olikes.PostLikesLikeResponseObject, error) {
	body := request.Body

	like, err := h.svc.Like(ctx, body.FromUserID, body.ToUserID)
	if err != nil {
		return olikes.PostLikesLike400JSONResponse(*dto.NewErrorDTO(http.StatusBadRequest, err.Error())), nil
	}

	nextCard, err := h.svc.GetFeed(ctx, body.FromUserID)
	if err != nil {
		return olikes.PostLikesLike500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, "failed to get next card: "+err.Error())), nil
	}

	var matchDTO *dto.MatchDTO
	match, err := h.matchSvc.GetLastBetween(ctx, body.FromUserID, body.ToUserID)
	if err != nil {
		return olikes.PostLikesLike500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, "failed to check match: "+err.Error())), nil
	}

	if match != nil {
		freelancer, err := h.freelancerSvc.GetByID(ctx, match.FreelancerID)
		if err != nil {
			return olikes.PostLikesLike500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, "failed to get freelancer: "+err.Error())), nil
		}
		if freelancer == nil {
			return olikes.PostLikesLike404JSONResponse(*dto.NewErrorDTO(http.StatusNotFound, "freelancer not found")), nil
		}

		project, err := h.projectSvc.GetByID(ctx, match.ProjectID)
		if err != nil {
			return olikes.PostLikesLike500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, "failed to get project: "+err.Error())), nil
		}
		if project == nil {
			return olikes.PostLikesLike404JSONResponse(*dto.NewErrorDTO(http.StatusNotFound, "project not found")), nil
		}

		freelancerUser, err := h.userSvc.GetByID(ctx, freelancer.UserID)
		if err != nil {
			return olikes.PostLikesLike500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, "failed to get freelancer user: "+err.Error())), nil
		}
		if freelancerUser == nil {
			return olikes.PostLikesLike404JSONResponse(*dto.NewErrorDTO(http.StatusNotFound, "freelancer user not found")), nil
		}

		projectUser, err := h.userSvc.GetByID(ctx, project.UserID)
		if err != nil {
			return olikes.PostLikesLike500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, "failed to get project user: "+err.Error())), nil
		}
		if projectUser == nil {
			return olikes.PostLikesLike404JSONResponse(*dto.NewErrorDTO(http.StatusNotFound, "project user not found")), nil
		}

		matchDTO = dto.NewMatchDTO(match, freelancer, project, freelancerUser, projectUser)
	}

	resp := dto.LikeResponse{
		Like:  dto.NewLikeDTO(like),
		Match: matchDTO,
		Next:  nextCard,
	}

	return olikes.PostLikesLike200JSONResponse(resp), nil
}

// POST /likes/dislike
func (h *likeHandler) PostLikesDislike(ctx context.Context, request olikes.PostLikesDislikeRequestObject) (olikes.PostLikesDislikeResponseObject, error) {
	body := request.Body

	fromID, err := uuid.Parse(body.FromUserID.String())
	if err != nil {
		return olikes.PostLikesDislike400JSONResponse(*dto.NewErrorDTO(http.StatusBadRequest, "invalid from_user_id")), nil
	}

	toID, err := uuid.Parse(body.ToUserID.String())
	if err != nil {
		return olikes.PostLikesDislike400JSONResponse(*dto.NewErrorDTO(http.StatusBadRequest, "invalid to_user_id")), nil
	}

	if err := h.svc.Dislike(ctx, fromID, toID); err != nil {
		return olikes.PostLikesDislike500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, err.Error())), nil
	}

	next, err := h.svc.GetFeed(ctx, fromID)
	if err != nil {
		return olikes.PostLikesDislike500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, err.Error())), nil
	}

	resp := dto.DislikeResponse{Next: next}

	return olikes.PostLikesDislike200JSONResponse(resp), nil
}
