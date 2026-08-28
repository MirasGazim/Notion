package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"notion/internal/config"
	"notion/internal/database/postgres"
	"notion/internal/handlers/http/auth"
	"notion/internal/handlers/http/blocks"
	"notion/internal/handlers/http/users"
	"notion/internal/handlers/http/workspace"
	"notion/internal/handlers/middleware/jwt"
	"notion/internal/handlers/middleware/logger"
	"notion/internal/lib/logger/sl"
	"notion/internal/repository"
	"notion/internal/service"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
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
		"starting notion",
		slog.String("env", cfg.Env),
		slog.String("version", "123"),
	)
	log.Debug("Debug messages are enabled")

	storage, err := postgres.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	defer storage.DB.Close()

	Nclient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	defer Nclient.Close()

	if err := Nclient.Ping(context.Background()).Err(); err != nil {
		log.Error("failed to connect to redis", sl.Err(err))
		os.Exit(1)
	}

	repos := repository.NewRepository(storage.DB, Nclient)
	services := service.NewService(repos)

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Post("/sign-in", auth.NewSignIn(log, services))
	router.Post("/sign-up", auth.NewSignUp(log, services))

	router.Group(func(r chi.Router) {
		r.Use(jwt.AuthMiddleware(log))

		r.Post("/Workspace", workspace.NewCreateWorkspace(log, services))
		r.Get("/Workspaces", workspace.GetAllWorkspaces(log, services))
		r.Get("/Workspaces/{id}/blocks", workspace.NewGetWorkspaceBlocks(log, services))
		r.Patch("/Workspaces/{id}", workspace.UpdateWorkspace(log, services))
		r.Delete("/Workspaces/{id}", workspace.DeleteWorkspace(log, services))
		r.Delete("/Users/", users.NewDelete(log, services))
		r.Post("/workspaces/{workspace_id}/blocks", blocks.NewCreateBlockHandler(log, services))
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
