package main

import (
	"log"
	"log/slog"
	"notion/internal/config"
	"notion/internal/database/postgres"
	"notion/internal/handlers/middleware/logger"
	"notion/internal/lib/logger/sl"
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

	_, err = postgres.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

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
