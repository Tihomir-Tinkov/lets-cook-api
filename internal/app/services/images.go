package services

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/models"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/ports"
	"github.com/google/uuid"
)

var allowedMimeTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

var (
	ErrEmptyFile       = errors.New("empty file")
	ErrUnsupportedMime = errors.New("unsupported image type")
)

type UploadRequest struct {
	Filename string
	Size     int64
	Reader   io.Reader
}

type DetectedFile struct {
	MimeType  string
	Extension string
	Reader    io.Reader
}

type ImageService struct {
	repository ports.ImageRepository
	storage    ports.FileStorage
}

func NewImageService(repository ports.ImageRepository, storage ports.FileStorage) *ImageService {
	return &ImageService{
		repository: repository,
		storage:    storage,
	}
}

func (s *ImageService) Upload(
	ctx context.Context,
	req UploadRequest,
) (uuid.UUID, error) {

	detectedFile, err := DetectFile(req.Reader)
	if err != nil {
		return uuid.Nil, err
	}

	image := &models.Image{
		FileName:  req.Filename,
		MimeType:  detectedFile.MimeType,
		Extension: detectedFile.Extension,
		Size:      req.Size,
	}

	err = s.repository.Create(ctx, image)
	if err != nil {
		return uuid.Nil, err
	}

	err = s.storage.Save(
		ctx,
		image.ID,
		detectedFile.Extension,
		detectedFile.Reader,
	)
	if err != nil {
		// cleanup DB row because storage failed
		_ = s.repository.Delete(ctx, image.ID)

		return uuid.Nil, err
	}

	return image.ID, nil
}

func (s *ImageService) Download(
	ctx context.Context,
	id uuid.UUID,
) (*models.Image, io.ReadCloser, error) {

	image, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	file, err := s.storage.Open(
		ctx,
		id,
		image.Extension,
	)
	if err != nil {
		return nil, nil, err
	}

	return image, file, nil
}

func (s *ImageService) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {

	image, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.storage.Delete(
		ctx,
		id,
		image.Extension,
	)
	if err != nil {
		return err
	}

	return s.repository.Delete(
		ctx,
		id,
	)
}

func DetectFile(r io.Reader) (*DetectedFile, error) {
	const sniffLen = 512

	buf := make([]byte, sniffLen)

	n, err := io.ReadFull(r, buf)
	switch err {
	case nil:
	case io.ErrUnexpectedEOF:
	case io.EOF:
		return nil, ErrEmptyFile
	default:
		return nil, err
	}

	mimeType := http.DetectContentType(buf[:n])

	extension, ok := allowedMimeTypes[mimeType]
	if !ok {
		return nil, ErrUnsupportedMime
	}

	return &DetectedFile{
		MimeType:  mimeType,
		Extension: extension,
		Reader: io.MultiReader(
			bytes.NewReader(buf[:n]),
			r,
		),
	}, nil
}
