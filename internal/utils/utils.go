package utils

import (
	"mime"
	"mime/multipart"
	"path/filepath"
)

var allowedMIMETypes = map[string]bool{
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/xml": true,
	"image/jpeg":      true,
	"image/png":       true,
}

var allowedEmailMIMETypes = map[string]bool{
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/pdf": true,
}

const (
	MaxArchiveSize        = 100 * 1024 * 1024
	MaxIndividualFileSize = 20 * 1024 * 1024
)

func IsValidMIMEType(mimeType string) bool {
	return allowedMIMETypes[mimeType]
}

func IsValidEmailMIMEType(mimeType string) bool {
	return allowedEmailMIMETypes[mimeType]
}

func GetMIMETypeFromFilename(filename string) string {
	ext := filepath.Ext(filename)
	return mime.TypeByExtension(ext)
}

func IsArchive(filename string) bool {
	ext := filepath.Ext(filename)
	return ext == ".zip"
}

func IsValidSize(fileHeader *multipart.FileHeader, maxSize int64) bool {
	return fileHeader.Size <= maxSize
}
