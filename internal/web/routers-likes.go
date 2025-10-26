package web

import (
	"context"
	"net/http"

	"github.com/kitoyanok66/workk/dto"
	"github.com/kitoyanok66/workk/internal/freelancers"
	"github.com/kitoyanok66/workk/internal/likes"
	"github.com/kitoyanok66/workk/internal/matches"
	"github.com/kitoyanok66/workk/internal/middleware"
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
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return olikes.PostLikesLike401JSONResponse(*dto.NewErrorDTO(http.StatusUnauthorized, "unauthorized")), nil
	}

	body := request.Body
	toUserID := body.ToUserID

	like, err := h.svc.Like(ctx, userID, toUserID)
	if err != nil {
		return olikes.PostLikesLike400JSONResponse(*dto.NewErrorDTO(http.StatusBadRequest, err.Error())), nil
	}

	nextCard, err := h.svc.GetFeed(ctx, userID)
	if err != nil {
		return olikes.PostLikesLike500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, "failed to get next card: "+err.Error())), nil
	}

	fromUser, err := h.userSvc.GetByID(ctx, userID)
	if err != nil || fromUser == nil {
		return olikes.PostLikesLike500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, "failed to get fromUser")), nil
	}
	toUser, err := h.userSvc.GetByID(ctx, toUserID)
	if err != nil || toUser == nil {
		return olikes.PostLikesLike500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, "failed to get toUser")), nil
	}

	freelancerID, projectID, err := h.svc.ResolvePair(ctx, fromUser, toUser)
	if err != nil {
		return olikes.PostLikesLike200JSONResponse(dto.LikeResponse{
			Like: dto.NewLikeDTO(like),
			Next: nextCard,
		}), nil
	}

	match, err := h.matchSvc.GetLastBetween(ctx, freelancerID, projectID)
	if err != nil {
		return olikes.PostLikesLike500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, "failed to check match: "+err.Error())), nil
	}

	var matchDTO *dto.MatchDTO
	if match != nil {
		freelancer, err := h.freelancerSvc.GetByID(ctx, match.FreelancerID)
		if err != nil {
			return olikes.PostLikesLike500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, "failed to get freelancer: "+err.Error())), nil
		}
		project, err := h.projectSvc.GetByID(ctx, match.ProjectID)
		if err != nil {
			return olikes.PostLikesLike500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, "failed to get project: "+err.Error())), nil
		}

		freelancerUser, err := h.userSvc.GetByID(ctx, freelancer.UserID)
		if err != nil {
			return olikes.PostLikesLike500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, "failed to get freelancer user: "+err.Error())), nil
		}
		projectUser, err := h.userSvc.GetByID(ctx, project.UserID)
		if err != nil {
			return olikes.PostLikesLike500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, "failed to get project user: "+err.Error())), nil
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
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return olikes.PostLikesDislike401JSONResponse(*dto.NewErrorDTO(http.StatusUnauthorized, "unauthorized")), nil
	}

	body := request.Body
	toUserID := body.ToUserID

	if err := h.svc.Dislike(ctx, userID, toUserID); err != nil {
		return olikes.PostLikesDislike500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, err.Error())), nil
	}

	next, err := h.svc.GetFeed(ctx, userID)
	if err != nil {
		return olikes.PostLikesDislike500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, err.Error())), nil
	}

	resp := dto.DislikeResponse{Next: next}

	return olikes.PostLikesDislike200JSONResponse(resp), nil
}
