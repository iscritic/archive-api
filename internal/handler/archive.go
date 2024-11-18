package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iscritic/archive-api/internal/service"
)

type ArchiveHandler struct {
	archiveService service.ArchiveService
}

func NewArchiveHandler(archiveService service.ArchiveService) *ArchiveHandler {
	return &ArchiveHandler{
		archiveService: archiveService,
	}
}

func (h *ArchiveHandler) GetArchiveInformation(c *gin.Context) {
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

	archiveInfo, err := h.archiveService.GetArchiveInfo(file, fileHeader.Filename, fileHeader.Size)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, archiveInfo)
}

func (h *ArchiveHandler) CreateArchive(c *gin.Context) {
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

	archiveData, err := h.archiveService.CreateZipArchive(files)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=archive.zip")
	c.Data(http.StatusOK, "application/zip", archiveData)
}
