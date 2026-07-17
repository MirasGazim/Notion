package users

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"notion/internal/handlers/middleware/ctx"

	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Deleter interface {
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

func NewDelete(log *slog.Logger, deleter Deleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contx := r.Context()
		const op = "handlers/http/workspace/GetWorkspace"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		ID, ok := r.Context().Value(ctx.UserIDKey).(uuid.UUID)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Info("userId got", slog.Any("UserID", ID))

		err := deleter.DeleteUser(contx, ID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			}
			log.Error("failed to delete user", "error", err, "workspace_id", ID)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
