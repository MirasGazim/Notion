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
	"notion/internal/service"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Creater interface {
	CreateUser(ctx context.Context, user user.SignUpRequest) (uuid.UUID, error)
}

type Getter interface {
	GetUser(ctx context.Context, u user.SignInRequest) (user.AuthUser, error)
}

const (
	salt       = "dfhgsdfhgidu1224"
	signingKey = "grkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 12 * time.Hour
)

type Response struct {
	response.Response
	UUID uuid.UUID `json:"id"`
}

type TokenResponse struct {
	response.Response
	Token string `json:"token"`
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

		jwtToken, err := GenerateToken(ctx, id)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{
				"error": err.Error(),
			})
			return
		}

		render.Status(r, http.StatusCreated)
		responseOK(w, r, jwtToken)
	}
}

func GenerateToken(ctx context.Context, id uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &service.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: id,
	})

	return token.SignedString([]byte(signingKey))
}

func NewSignIn(log *slog.Logger, getter Getter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		const op = "handlers/http/auth/NewSignIn"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var signIn user.SignInRequest
		err := render.DecodeJSON(r.Body, &signIn)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", signIn))

		if err := validator.New().Struct(signIn); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		auth, err := getter.GetUser(ctx, signIn)
		if err != nil {
			if errors.Is(err, service.ErrInvalidCredentials) {
				log.Info("User Authentication failed", slog.String("user", signIn.Username))

				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, response.Error("Authentication failed"))

				return
			}
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("Internal Server Error"))

			return
		}

		jwttoken, err := GenerateToken(ctx, auth.ID)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{
				"error": err.Error(),
			})
			return
		}

		render.Status(r, http.StatusOK)
		responseOK(w, r, jwttoken)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, token string) {
	render.JSON(w, r, TokenResponse{
		Response: response.Ok(),
		Token:    token,
	})
}
