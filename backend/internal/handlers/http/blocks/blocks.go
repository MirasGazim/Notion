package blocks

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"notion/internal/handlers/middleware/ctx"
	"notion/internal/lib/api/response"
	"notion/internal/lib/logger/sl"
	"notion/internal/models/blocks"
	"notion/internal/repository"
	"notion/internal/service"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type Creater interface {
	CreateBlock(ctx context.Context, req blocks.CreateBlockRequest, workspaceID, userID uuid.UUID) (*blocks.Block, error)
}

func NewCreateBlockHandler(log *slog.Logger, bl Creater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contx := r.Context()
		const op = "handlers/http/workspace/CreateBlocks"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var blocksReq blocks.CreateBlockRequest
		err := render.DecodeJSON(r.Body, &blocksReq)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}
		log.Info("request body decoded", slog.Any("request", blocksReq))

		userID, ok := r.Context().Value(ctx.UserIDKey).(uuid.UUID)
		if !ok {
			log.Error("failed to get user_id from context")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error("unauthorized"))
			return
		}
		log.Info("userId got", slog.Any("UserID", userID))

		urlid := chi.URLParam(r, "workspace_id")
		workspaceID, err := uuid.Parse(urlid)
		if err != nil {
			log.Error("failed to parse workspace_id", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid workspace_id"))
			return
		}

		block, err := bl.CreateBlock(contx, blocksReq, workspaceID, userID)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrWorkspaceNotFound):
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, response.Error("workspace not found"))
			case errors.Is(err, service.ErrInvalidContent):
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, response.Error(response.RootCause(err).Error()))
			default:
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, response.Error("failed to create block"))
			}
			return
		}
		log.Info("block created", slog.Any("block", block))
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, response.Created())
	}
}
