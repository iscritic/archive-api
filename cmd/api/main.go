package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iscritic/archive-api/internal/handlers"
)

func main() {
	router := gin.Default()

	api := router.Group("/api")
	{
		archive := api.Group("/archive")
		{
			archive.POST("/information", handlers.GetArchiveInformation)
			archive.POST("/files", handlers.CreateArchive)
		}

	}

	router.Run(":8080")
}
