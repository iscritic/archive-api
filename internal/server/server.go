package server

import (
	"fmt"
	"net/http"

	"github.com/iscritic/archive-api/internal/config"

	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg    *config.Config
	engine *gin.Engine
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg:    cfg,
		engine: gin.Default(),
	}
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.HTTPServer.Address, s.cfg.HTTPServer.Port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      s.engine,
		ReadTimeout:  s.cfg.HTTPServer.Timeout,
		WriteTimeout: s.cfg.HTTPServer.Timeout,
		IdleTimeout:  s.cfg.HTTPServer.IdleTimeout,
	}

	return srv.ListenAndServe()
}
