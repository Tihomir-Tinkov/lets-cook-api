package models

import (
	"time"

	"github.com/google/uuid"
)

type Image struct {
	ID        uuid.UUID `json:"id" db:"id"`
	FileName  string    `json:"-" db:"filename"`
	MimeType  string    `json:"-" db:"mime_type"`
	Extension string    `json:"-" db:"extension"`
	Size      int64     `json:"-" db:"size"`
	CreatedAt time.Time `json:"-" db:"created_at"`
}

func (Image) Table() string {
	return "images"
}
