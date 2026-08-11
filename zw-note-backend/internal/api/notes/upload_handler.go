package notes

import (
	"errors"

	"zw-note-backend/internal/config"
	"zw-note-backend/pkg/upload"
	"zw-note-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UploadHandler handles image uploads used by the Markdown editor's "insert
// image" toolbar action. It reuses the same local-disk storage as the admin
// upload endpoint (see pkg/upload); files are served back via the "/uploads"
// static route registered in bootstrap/router.go.
type UploadHandler struct {
	cfg config.UploadConfig
	log *zap.Logger
}

func NewUploadHandler(cfg config.UploadConfig, log *zap.Logger) *UploadHandler {
	return &UploadHandler{cfg: cfg, log: log}
}

// UploadImage godoc
// @Summary  Upload an image for use in note content
// @Tags     notes-upload
// @Accept   multipart/form-data
// @Produce  json
// @Param    file formData file true "Image file (jpg/png/webp/gif)"
// @Success  200  {object} utils.Response{data=object}
// @Router   /api/notes/v1/upload/image [post]
func (h *UploadHandler) UploadImage(c *gin.Context) {
	if _, ok := getNotesUserID(c); !ok {
		utils.Unauthorized(c)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		utils.ParamError(c, "missing or invalid file")
		return
	}

	relPath, err := upload.SaveImage(h.cfg, file)
	if err != nil {
		var verr *upload.ValidationError
		if errors.As(err, &verr) {
			utils.ParamError(c, verr.Error())
			return
		}
		h.log.Error("save uploaded image", zap.Error(err))
		utils.ServerError(c)
		return
	}

	h.log.Info("note image uploaded", zap.String("path", relPath))
	utils.Success(c, gin.H{
		"path": relPath,
		"url":  "/uploads/" + relPath,
	})
}
