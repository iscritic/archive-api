package handler

import (
	"bytes"
	"github.com/iscritic/archive-api/internal/utils"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iscritic/archive-api/internal/service"
	"log/slog"
)

type EmailHandler struct {
	EmailService *service.EmailService
}

func NewEmailHandler(emailService *service.EmailService) *EmailHandler {
	return &EmailHandler{
		EmailService: emailService,
	}
}

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func (h *EmailHandler) SendFileToEmails(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		slog.Error("Error retrieving file", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not received: " + err.Error()})
		return
	}
	defer file.Close()

	slog.Debug("File successfully received", "filename", header.Filename)

	emails := c.PostForm("emails")
	if emails == "" {
		slog.Warn("No email addresses provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No email addresses provided"})
		return
	}

	emailList := strings.Split(emails, ",")

	contentType := header.Header.Get("Content-Type")
	if !utils.IsValidEmailMIMEType(contentType) {
		slog.Warn("Invalid file type", "content_type", contentType)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only DOCX and PDF are supported."})
		return
	}

	fileBuffer := new(bytes.Buffer)
	_, err = fileBuffer.ReadFrom(file)
	if err != nil {
		slog.Error("Error reading file", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file: " + err.Error()})
		return
	}

	sentEmails := []string{}
	failedEmails := []string{}

	for _, email := range emailList {
		email = strings.TrimSpace(email)
		if email == "" {
			continue
		}

		if !isValidEmail(email) {
			slog.Warn("Invalid email format", "email", email)
			failedEmails = append(failedEmails, email)
			continue
		}

		err := h.EmailService.SendEmail(
			email,
			"File Attachment",
			"Hello, please find the attached document.",
			header.Filename,
			bytes.NewReader(fileBuffer.Bytes()),
			contentType,
		)

		if err != nil {
			slog.Error("Failed to send email", "email", email, "error", err)
			failedEmails = append(failedEmails, email)
		} else {
			slog.Info("Email successfully sent", "email", email)
			sentEmails = append(sentEmails, email)
		}
	}

	if len(failedEmails) > 0 {
		c.JSON(http.StatusPartialContent, gin.H{
			"message":      "Some emails could not be sent.",
			"sentEmails":   sentEmails,
			"failedEmails": failedEmails,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message":    "Emails successfully sent to all recipients",
			"sentEmails": sentEmails,
		})
	}
}
