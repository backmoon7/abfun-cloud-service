package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

var (
	LocalStoragePath string
	DomainURL        string
	LocalStorage     *LocalStorageService
)

type LocalStorageService struct{}

// InitLocalStorage initializes the local storage directory
func InitLocalStorage(path, domain string) {
	LocalStoragePath = path
	DomainURL = domain
	LocalStorage = &LocalStorageService{}

	// Create directory if not exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			panic(fmt.Sprintf("Failed to create storage directory: %v", err))
		}
	}
}

// Upload saves the file to local disk and returns the access URL
// If the file is already at the path (e.g. uploaded by Gin), it just returns the URL
func (s *LocalStorageService) Upload(path string) (string, error) {
	filename := filepath.Base(path)
	// Ensure DomainURL doesn't have trailing slash
	domain := strings.TrimRight(DomainURL, "/")
	return fmt.Sprintf("%s/%s", domain, filename), nil
}

// SaveFile saves a multipart file to storage
func (s *LocalStorageService) SaveFile(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	filename := filepath.Base(file.Filename)
	dstPath := filepath.Join(LocalStoragePath, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return s.Upload(dstPath)
}
