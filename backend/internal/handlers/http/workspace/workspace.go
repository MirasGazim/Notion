package workspace

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"notion/internal/handlers/middleware/ctx"
	"notion/internal/lib/api/response"
	"notion/internal/lib/logger/sl"
	"notion/internal/models/blocks"
	"notion/internal/models/workspace"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Creater interface {
	Create(ctx context.Context, req workspace.CreateWorkspaceRequest) (*workspace.Workspace, error)
}

type Getter interface {
	GetWorkspaces(ctx context.Context, id uuid.UUID) ([]workspace.Workspace, error)
}

type GetWorkspaceBlocks interface {
	GetByID(ctx context.Context, id uuid.UUID) (workspace.Workspace, error)
	GetByWorkspaceID(ctx context.Context, id uuid.UUID) ([]blocks.Block, error)
}

func NewCreateWorkspace(log *slog.Logger, creater Creater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contx := r.Context()
		const op = "handlers/http/workspace/CreateWorkspace"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var ws workspace.CreateWorkspaceRequest
		err := render.DecodeJSON(r.Body, &ws)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", ws))

		if err := validator.New().Struct(ws); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		userId, ok := r.Context().Value(ctx.UserIDKey).(uuid.UUID)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ws.ID = userId
		created, err := creater.Create(contx, ws)
		if err != nil {
			log.Error("failed to create workspace", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to create workspace"))

			return
		}
		log.Info("workspace created", slog.Any("workspace", created))

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, created)
	}
}

func GetAllWorkspaces(log *slog.Logger, getter Getter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contx := r.Context()
		const op = "handlers/http/workspace/GetWorkspace"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userId, ok := r.Context().Value(ctx.UserIDKey).(uuid.UUID)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return

		}
		log.Info("userId got", slog.Any("UserID", userId))
		workspace, err := getter.GetWorkspaces(contx, userId)
		if err != nil {
			log.Error("failed to get workspace", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get workspace"))

			return
		}

		log.Info("workspace got")

		render.Status(r, http.StatusOK)
		render.JSON(w, r, workspace)

	}

}

func NewGetWorkspaceBlocks(log *slog.Logger, getter GetWorkspaceBlocks) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contx := r.Context()
		const op = "handlers/http/workspace/GetWorkspace"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		UserID := chi.URLParam(r, "id")
		id, err := uuid.Parse(UserID)
		if err != nil {
			log.Error("invalid id", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid id"))

			return
		}

		var ws workspace.WorkspaceBlocks

		ws.Workspace, err = getter.GetByID(contx, id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, response.Error("workspace not found"))
				return
			}
			log.Error("failed to get workspace blocks", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get workspace blocks"))

			return
		}
		ws.Blocks, err = getter.GetByWorkspaceID(contx, id)
		if err != nil {
			log.Error("failed to get Blocks by id", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get blocks"))

			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, ws)
	}
}
