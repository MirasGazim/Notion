package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"notion/internal/config"
	"notion/internal/database/postgres"
	"notion/internal/handlers/http/auth"
	"notion/internal/handlers/middleware/jwt"
	"notion/internal/handlers/middleware/logger"
	"notion/internal/lib/logger/sl"
	"notion/internal/repository"
	"notion/internal/service"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("file does not exist %s", err)
	}
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info(
		"starting url-shortener",
		slog.String("env", cfg.Env),
		slog.String("version", "123"),
	)
	log.Debug("Debug messages are enabled")

	storage, err := postgres.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	repos := repository.NewRepository(storage.DB)
	services := service.NewService(repos)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/Sign-p", auth.NewSignUp(log, services))
	router.Post("/Sign-In", auth.NewSignIn(log, services))

	router.Group(func(r chi.Router) {
		r.Use(jwt.AuthMiddleware(log))

		r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value("user_id")
			w.Write([]byte(fmt.Sprintf("user_id from token: %v", userID)))
		})
	})

	log.Info("starting server", slog.String("address", cfg.Address))
	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}
func setupLogger(env string) *slog.Logger {
	var level slog.Level
	switch env {
	case envLocal:
		level = slog.LevelDebug
	case envDev:
		level = slog.LevelDebug
	case envProd:
		level = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	return slog.New(handler)
}
