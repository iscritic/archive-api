package handler

import (
	"github.com/iscritic/archive-api/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iscritic/archive-api/internal/service"
	"log/slog"
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
		slog.Warn("File retrieval failed",
			slog.String("error", err.Error()),
			slog.String("field", "file"),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a file in 'file' field."})
		return
	}

	if !utils.IsValidSize(fileHeader, utils.MaxArchiveSize) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size is exceed"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		slog.Error("Unable to open uploaded file",
			slog.String("filename", fileHeader.Filename),
			slog.Int64("size", fileHeader.Size),
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open the uploaded file."})
		return
	}
	defer file.Close()

	archiveInfo, err := h.archiveService.GetArchiveInfo(file, fileHeader.Filename, fileHeader.Size)
	if err != nil {
		slog.Warn("GetArchiveInfo failed",
			slog.String("filename", fileHeader.Filename),
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slog.Info("GetArchiveInformation succeeded",
		slog.String("filename", fileHeader.Filename),
		slog.Int64("size", fileHeader.Size),
	)
	c.JSON(http.StatusOK, archiveInfo)
}

func (h *ArchiveHandler) CreateArchive(c *gin.Context) {

	form, err := c.MultipartForm()
	if err != nil {
		slog.Warn("Multipart form retrieval failed",
			slog.String("error", err.Error()),
			slog.String("field", "files[]"),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide files in 'files[]' field."})
		return
	}

	files := form.File["files[]"]
	if len(files) == 0 {
		slog.Warn("No files provided in 'files[]' field")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files provided."})
		return
	}

	slog.Info("Creating zip archive",
		slog.Int("file_count", len(files)),
	)

	archiveData, err := h.archiveService.CreateZipArchive(files)
	if err != nil {
		slog.Error("CreateZipArchive failed",
			slog.String("error", err.Error()),
			slog.Int("file_count", len(files)),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Zip archive created successfully",
		slog.Int("archive_size", len(archiveData)),
	)

	c.Header("Content-Disposition", "attachment; filename=archive.zip")
	c.Data(http.StatusOK, "application/zip", archiveData)
}
