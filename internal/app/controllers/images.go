package controllers

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/controllers/responses"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/services"
	"github.com/google/uuid"
)

type ImageController struct {
	imageService *services.ImageService
}

func NewImageController(imageService *services.ImageService) *ImageController {
	return &ImageController{
		imageService: imageService,
	}
}

func (h *ImageController) Upload(w http.ResponseWriter, r *http.Request) {
	const maxUploadSize = 1024 * 1024 * 20 // 20 MB

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		responses.JSONError(w, r, errors.New("invalid multipart request"), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		responses.JSONError(w, r, errors.New("file is required"), http.StatusBadRequest)
		return
	}
	defer file.Close()

	image, err := h.imageService.Upload(
		r.Context(),
		services.UploadRequest{
			Filename: header.Filename,
			Size:     header.Size,
			Reader:   file,
		},
	)
	if err != nil {
		responses.JSONError(w, r, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	responses.JSONResponse(w, http.StatusCreated, image)
}

func (h *ImageController) Download(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(
		r.PathValue("id"),
	)

	if err != nil {
		responses.JSONError(w, r, errors.New("invalid_uuid"), http.StatusBadRequest)
		return
	}

	meta, file, err := h.imageService.Download(
		r.Context(),
		id,
	)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			responses.JSONError(w, r, err, http.StatusNotFound)
			return
		}
		responses.JSONError(w, r, err, http.StatusInternalServerError)
		return
	}

	defer file.Close()

	w.Header().Set("Content-Type", meta.MimeType)
	w.Header().Set("Content-Length", strconv.FormatInt(meta.Size, 10))

	w.Header().Set(
		"Cache-Control",
		"public,max-age=31536000,immutable",
	)

	w.Header().Set(
		"ETag",
		meta.ID.String(),
	)

	_, _ = io.Copy(w, file)
}

func (h *ImageController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(
		r.PathValue("id"),
	)

	if err != nil {
		responses.JSONError(w, r, errors.New("invalid_uuid"), http.StatusBadRequest)
		return
	}

	err = h.imageService.Delete(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			responses.JSONError(w, r, err, http.StatusNotFound)
			return
		}
		responses.JSONError(w, r, err, http.StatusInternalServerError)
		return
	}
	responses.JSONResponse(w, http.StatusOK, "image was deleted successfully")
}
