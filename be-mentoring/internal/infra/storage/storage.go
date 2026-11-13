package storage

import (
	"context"
	"io"
)

// StorageService adalah abstraksi untuk file storage (Local / AWS S3 / MinIO)
type StorageService interface {
	Upload(ctx context.Context, filename string, src io.Reader) (string, error)
	Delete(ctx context.Context, filename string) error
}

type localStorageService struct {
	basePath string
}

func NewLocalStorageService(basePath string) StorageService {
	return &localStorageService{basePath: basePath}
}

func (s *localStorageService) Upload(ctx context.Context, filename string, src io.Reader) (string, error) {
	// Contoh implementasi local storage sederhana
	return s.basePath + "/" + filename, nil
}

func (s *localStorageService) Delete(ctx context.Context, filename string) error {
	return nil
}
