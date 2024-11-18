package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iscritic/archive-api/internal/config"
	"github.com/iscritic/archive-api/internal/server"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"os"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env")
	}

	cfg := config.MustLoad()
	slog.Info("Config is set...")

	gin.SetMode(gin.ReleaseMode)
	srv := server.NewServer(cfg)

	slog.Info("Starting API server",
		slog.String("address", cfg.HTTPServer.Address),
		slog.Int("port", cfg.HTTPServer.Port),
	)

	err = srv.Run()
	if err != nil {
		slog.Error("Failed to run server", slog.String("error", err.Error()))
		os.Exit(1)
	}

}
