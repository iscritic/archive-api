package service

import (
	"archive/zip"
	"bytes"
	"errors"
	"github.com/iscritic/archive-api/internal/utils"
	"github.com/iscritic/archive-api/models"
	"io"
	"log/slog"
	"mime/multipart"
)

type ArchiveService interface {
	GetArchiveInfo(file multipart.File, filename string, fileSize int64) (*models.ArchiveInfo, error)
	CreateZipArchive(files []*multipart.FileHeader) ([]byte, error)
}

type archiveService struct{}

func NewArchiveService() ArchiveService {
	return &archiveService{}
}

func (s *archiveService) GetArchiveInfo(file multipart.File, filename string, fileSize int64) (*models.ArchiveInfo, error) {
	if !utils.IsArchive(filename) {
		return nil, errors.New("Oops! The file you uploaded is not an archive. Please try again with a valid archive file.")
	}

	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, file)
	if err != nil {
		return nil, err
	}

	reader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), fileSize)
	if err != nil {
		return nil, errors.New("Failed to read the archive. Make sure it's a valid ZIP file.")
	}

	var totalSize float64
	var files []models.FileInfo

	for _, f := range reader.File {
		fileInfo := models.FileInfo{
			FilePath: f.Name,
			Size:     float64(f.UncompressedSize64),
			MIMEType: utils.GetMIMETypeFromFilename(f.Name),
		}
		totalSize += fileInfo.Size
		files = append(files, fileInfo)
	}

	archiveInfo := &models.ArchiveInfo{
		Filename:    filename,
		ArchiveSize: float64(fileSize),
		TotalSize:   totalSize,
		TotalFiles:  len(files),
		Files:       files,
	}

	return archiveInfo, nil
}

func (s *archiveService) CreateZipArchive(files []*multipart.FileHeader) ([]byte, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for _, fileHeader := range files {
		if !utils.IsValidMIMEType(fileHeader.Header.Get("Content-Type")) {
			return nil, errors.New("One or more files have unsupported formats. Please upload files with allowed MIME types.")
		}

		if !utils.IsValidSize(fileHeader, utils.MaxIndividualFileSize) {
			slog.Warn("File size exceeds the limit",
				slog.String("filename", fileHeader.Filename),
				slog.Int64("size", fileHeader.Size),
			)
			return nil, errors.New("File size exceeds the limit (20MB) per file")
		}

		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}

		w, err := zipWriter.Create(fileHeader.Filename)
		if err != nil {
			file.Close()
			return nil, err
		}

		_, err = io.Copy(w, file)
		file.Close()
		if err != nil {
			return nil, err
		}
	}

	err := zipWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
