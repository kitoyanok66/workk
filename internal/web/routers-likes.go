package web

import (
	"context"
	"fmt"
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

func (h *likeHandler) resolvePair(ctx context.Context, fromUserID, toUserID uuid.UUID) (freelancerID, projectID uuid.UUID, err error) {
	fromUser, err := h.userSvc.GetByID(ctx, fromUserID)
	if err != nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("failed to get fromUser: %w", err)
	}
	if fromUser == nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("fromUser not found")
	}

	toUser, err := h.userSvc.GetByID(ctx, toUserID)
	if err != nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("failed to get toUser: %w", err)
	}
	if toUser == nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("toUser not found")
	}

	switch {
	case fromUser.Role == "freelancer" && toUser.Role == "project":
		freelancer, err := h.freelancerSvc.GetByUserID(ctx, fromUser.ID)
		if err != nil {
			return uuid.Nil, uuid.Nil, fmt.Errorf("failed to get freelancer by fromUserID: %w", err)
		}
		project, err := h.projectSvc.GetByUserID(ctx, toUser.ID)
		if err != nil {
			return uuid.Nil, uuid.Nil, fmt.Errorf("failed to get project by toUserID: %w", err)
		}
		return freelancer.ID, project.ID, nil

	case fromUser.Role == "project" && toUser.Role == "freelancer":
		freelancer, err := h.freelancerSvc.GetByUserID(ctx, toUser.ID)
		if err != nil {
			return uuid.Nil, uuid.Nil, fmt.Errorf("failed to get freelancer by toUserID: %w", err)
		}
		project, err := h.projectSvc.GetByUserID(ctx, fromUser.ID)
		if err != nil {
			return uuid.Nil, uuid.Nil, fmt.Errorf("failed to get project by fromUserID: %w", err)
		}
		return freelancer.ID, project.ID, nil

	default:
		return uuid.Nil, uuid.Nil, fmt.Errorf("invalid role combination: %s -> %s", fromUser.Role, toUser.Role)
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

	freelancerID, projectID, err := h.resolvePair(ctx, body.FromUserID, body.ToUserID)
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
	body := request.Body

	if err := h.svc.Dislike(ctx, body.FromUserID, body.ToUserID); err != nil {
		return olikes.PostLikesDislike500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, err.Error())), nil
	}

	next, err := h.svc.GetFeed(ctx, body.FromUserID)
	if err != nil {
		return olikes.PostLikesDislike500JSONResponse(*dto.NewErrorDTO(http.StatusInternalServerError, err.Error())), nil
	}

	resp := dto.DislikeResponse{Next: next}

	return olikes.PostLikesDislike200JSONResponse(resp), nil
}
