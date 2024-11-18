package handler

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/iscritic/archive-api/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestGetArchiveInformation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	archiveService := service.NewArchiveService()
	handler := NewArchiveHandler(archiveService)

	var buf bytes.Buffer
	zipWriter := multipart.NewWriter(&buf)
	part, _ := zipWriter.CreateFormFile("file", "test.zip")
	part.Write([]byte("fake zip content"))
	zipWriter.Close()

	req, _ := http.NewRequest(http.MethodPost, "/api/archive/information", &buf)
	req.Header.Set("Content-Type", zipWriter.FormDataContentType())

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetArchiveInformation(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
