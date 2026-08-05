package admin

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

	"go-web-api/internal/config"
	"go-web-api/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var allowedImageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true,
}

// UploadHandler handles file uploads (admin only).
type UploadHandler struct {
	cfg config.UploadConfig
	log *zap.Logger
}

func NewUploadHandler(cfg config.UploadConfig, log *zap.Logger) *UploadHandler {
	return &UploadHandler{cfg: cfg, log: log}
}

// UploadImage godoc
// @Summary  Upload an image
// @Tags     admin-upload
// @Security BearerAuth
// @Accept   multipart/form-data
// @Produce  json
// @Param    file formData file true "Image file (jpg/png/webp/gif)"
// @Success  200  {object} utils.Response{data=object}
// @Router   /api/admin/v1/upload/image [post]
func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.ParamError(c, "missing or invalid file")
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedImageExts[ext] {
		utils.ParamError(c, "only jpg, png, webp, gif are allowed")
		return
	}

	maxMB := h.cfg.MaxSizeMB
	if maxMB <= 0 {
		maxMB = 5
	}
	maxBytes := int64(maxMB) * 1024 * 1024
	if file.Size > maxBytes {
		utils.ParamError(c, fmt.Sprintf("file size exceeds %dMB limit", maxMB))
		return
	}

	// Save as images/{date}/{random}{ext}
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		h.log.Error("generate random name", zap.Error(err))
		utils.ServerError(c)
		return
	}
	relPath := fmt.Sprintf("images/%s/%s%s", time.Now().Format("20060102"), hex.EncodeToString(b), ext)
	uploadDir := h.cfg.Dir
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	fullPath := filepath.Join(uploadDir, relPath)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		h.log.Error("create upload dir", zap.Error(err))
		utils.ServerError(c)
		return
	}

	if err := saveUploadedFile(file, fullPath); err != nil {
		h.log.Error("save uploaded file", zap.Error(err))
		utils.ServerError(c)
		return
	}

	h.log.Info("image uploaded", zap.String("path", relPath))
	utils.Success(c, gin.H{"path": relPath})
}

func saveUploadedFile(file *multipart.FileHeader, dest string) error {
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
