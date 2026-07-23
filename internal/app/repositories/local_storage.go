package repositories

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type LocalStorage struct{ Root string }

func NewLocalStorage(root string) *LocalStorage {
	return &LocalStorage{Root: root}
}

func (s *LocalStorage) path(id uuid.UUID, extention string) string {
	return filepath.Join(s.Root, id.String()+extention)
}

func (s *LocalStorage) Save(_ context.Context, id uuid.UUID, extension string, r io.Reader) error {
	if err := os.MkdirAll(s.Root, 0755); err != nil {
		return err
	}
	f, err := os.Create(s.path(id, extension))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	return err
}

func (s *LocalStorage) Open(_ context.Context, id uuid.UUID, extension string) (io.ReadCloser, error) {
	return os.Open(s.path(id, extension))
}

func (s *LocalStorage) Delete(_ context.Context, id uuid.UUID, extension string) error {
	return os.Remove(s.path(id, extension))
}
