package workspace

import (
	"context"
	"log/slog"
	"net/http"
	"notion/internal/handlers/middleware/ctx"
	"notion/internal/lib/api/response"
	"notion/internal/lib/logger/sl"
	"notion/internal/models/workspace"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Creater interface {
	Create(ctx context.Context, req workspace.CreateWorkspaceRequest) (*workspace.Workspace, error)
}

func NewCreateWorkspace(log *slog.Logger, creater Creater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contx := r.Context()
		const op = "handlers/http/workspace/CreateCreateWorkspace"
		log = log.With(
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
