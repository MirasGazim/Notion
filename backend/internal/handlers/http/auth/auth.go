package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"notion/internal/lib/api/response"
	"notion/internal/lib/logger/sl"
	"notion/internal/models/user"
	"notion/internal/repository"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Creater interface {
	CreateUser(ctx context.Context, user user.SignUpRequest) (uuid.UUID, error)
}

type Response struct {
	response.Response
	UUID uuid.UUID `json:"id"`
}

func NewSignUp(log *slog.Logger, creater Creater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		const op = "handlers/http/auth/NewSignUp"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req_u user.SignUpRequest
		err := render.DecodeJSON(r.Body, &req_u)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req_u))

		if err := validator.New().Struct(req_u); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		id, err := creater.CreateUser(ctx, req_u)
		if errors.Is(err, repository.ErrUserExists) {
			log.Info("user already exists", slog.String("user", req_u.Username))

			render.Status(r, http.StatusConflict)
			render.JSON(w, r, response.Error("user already exists"))
			return
		}

		if err != nil {
			log.Error("failed to create user", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to create user"))
			return
		}

		log.Info("user created")

		render.Status(r, http.StatusCreated)
		responseOK(w, r, id)

	}
}

// func NewSignIn(log *slog.Logger) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		const op = "handlers/http/auth/NewSignIn"
// 		log = log.With(
// 			slog.String("op", op),
// 			slog.String("request_id", middleware.GetReqID(r.Context())),
// 		)
// 	}
// }

func responseOK(w http.ResponseWriter, r *http.Request, uuid uuid.UUID) {
	render.JSON(w, r, Response{
		Response: response.Ok(),
		UUID:     uuid,
	})
}
