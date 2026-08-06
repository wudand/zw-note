package mp

import (
	"zw-note-backend/internal/dto"
	"zw-note-backend/internal/model"
	"zw-note-backend/internal/service"
	"zw-note-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CategoryHandler groups all mini-program category HTTP handlers (read-only).
type CategoryHandler struct {
	svc service.CategoryService
	log *zap.Logger
}

func NewCategoryHandler(svc service.CategoryService, log *zap.Logger) *CategoryHandler {
	return &CategoryHandler{svc: svc, log: log}
}

// List godoc
// @Summary  List enabled categories (public, no auth)
// @Tags     mp-categories
// @Produce  json
// @Success  200 {object} utils.Response{data=[]dto.CategoryResponse}
// @Router   /api/mp/v1/categories [get]
func (h *CategoryHandler) List(c *gin.Context) {
	cats, _, err := h.svc.List(c.Request.Context(), 1, 1000, false)
	if err != nil {
		h.log.Error("list categories", zap.Error(err))
		utils.ServerError(c)
		return
	}

	list := make([]*dto.CategoryResponse, 0, len(cats))
	for _, cat := range cats {
		list = append(list, toCategoryResponse(cat))
	}

	utils.Success(c, list)
}

func toCategoryResponse(c *model.Category) *dto.CategoryResponse {
	return &dto.CategoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		SortOrder: c.SortOrder,
		Status:    c.Status,
		CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
