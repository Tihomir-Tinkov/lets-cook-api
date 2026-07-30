package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/controllers/middleware"
	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/models"
	"github.com/google/uuid"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

func requireAuth(r *http.Request) (*models.AuthContext, error) {
	auth, ok := middleware.GetAuthContext(r.Context())
	if !ok {
		return nil, ErrUnauthorized
	}

	return auth, nil
}

func parseUUIDParam(r *http.Request, name string) (uuid.UUID, error) {
	id := r.PathValue(name)

	if id == "" {
		return uuid.Nil, errors.New("missing path parameter: " + name)
	}

	parsed, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, errors.New("invalid " + name)
	}

	return parsed, nil
}

func decodeJSON(r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(nil, r.Body, 1<<20) // 1 MB

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return err
	}

	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		return errors.New("request body must contain only one JSON object")
	}

	return nil
}
