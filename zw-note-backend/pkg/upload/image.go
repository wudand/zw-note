// Package upload provides shared local-disk file storage used by the admin
// and notes upload endpoints. Files are saved under a configured directory
// and served back via the "/uploads" static route registered in
// bootstrap/router.go.
package upload

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"zw-note-backend/internal/config"
)

// AllowedImageExts is the set of accepted image file extensions (lower-case, with dot).
var AllowedImageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true,
}

const defaultMaxSizeMB = 5

// ValidationError indicates the upload itself is fine but the file failed a
// business rule (extension, size). Handlers should surface it as a 400.
type ValidationError struct {
	msg string
}

func (e *ValidationError) Error() string { return e.msg }

func newValidationError(format string, args ...any) *ValidationError {
	return &ValidationError{msg: fmt.Sprintf(format, args...)}
}

// SaveImage validates and persists an uploaded image under
// {cfg.Dir}/images/{yyyyMMdd}/{random}{ext}, returning the path relative to
// cfg.Dir (e.g. "images/20260811/ab12cd34.png"). That relative path is what
// gets stored (prefixed with "/uploads/") in Markdown content and served by
// the "/uploads" static route.
func SaveImage(cfg config.UploadConfig, file *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !AllowedImageExts[ext] {
		return "", newValidationError("only jpg, png, webp, gif are allowed")
	}

	maxMB := cfg.MaxSizeMB
	if maxMB <= 0 {
		maxMB = defaultMaxSizeMB
	}
	if file.Size > int64(maxMB)*1024*1024 {
		return "", newValidationError("file size exceeds %dMB limit", maxMB)
	}

	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate random name: %w", err)
	}
	relPath := fmt.Sprintf("images/%s/%s%s", time.Now().Format("20060102"), hex.EncodeToString(b), ext)

	uploadDir := cfg.Dir
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	fullPath := filepath.Join(uploadDir, relPath)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", fmt.Errorf("create upload dir: %w", err)
	}

	if err := saveMultipartFile(file, fullPath); err != nil {
		return "", fmt.Errorf("save uploaded file: %w", err)
	}

	return relPath, nil
}

func saveMultipartFile(file *multipart.FileHeader, dest string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
