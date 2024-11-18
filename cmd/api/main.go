package main

import (
	"fmt"
	"github.com/iscritic/archive-api/internal/config"
	"github.com/iscritic/archive-api/internal/server"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env")
	}

	cfg := config.MustLoad()

	fmt.Println(cfg)
	srv := server.NewServer(cfg)

	slog.Info("Starting API server...")
	err = srv.Run()
	if err != nil {
		log.Fatal(err)
	}

}
