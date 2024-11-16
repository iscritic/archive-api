package handlers

import (
	"github.com/iscritic/archive-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetArchiveInformation(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a file in 'file' field."})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open the uploaded file."})
		return
	}
	defer file.Close()

	archiveInfo, err := services.GetArchiveInfo(file, fileHeader.Filename, fileHeader.Size)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, archiveInfo)
}

func CreateArchive(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide files in 'files[]' field."})
		return
	}

	files := form.File["files[]"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files provided."})
		return
	}

	archiveData, err := services.CreateZipArchive(files)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=archive.zip")
	c.Data(http.StatusOK, "application/zip", archiveData)
}
