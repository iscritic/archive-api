package server

import (
	"github.com/gin-gonic/gin"
	"github.com/iscritic/archive-api/internal/handler"
	"github.com/iscritic/archive-api/internal/service"
	"net/http"
)

func (s *Server) registerRoutes() {
	archiveService := service.NewArchiveService()
	archiveHandler := handler.NewArchiveHandler(archiveService)

	emailSenderService := service.NewEmailService(&s.cfg.SMTP)
	emailSenderHandler := handler.NewEmailHandler(emailSenderService)

	api := s.engine.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})

		api.POST("/archive/information", archiveHandler.GetArchiveInformation)
		api.POST("/archive/files", archiveHandler.CreateArchive)

		api.POST("/mail/file", emailSenderHandler.SendFileToEmails)
	}
}
