package admin

import (
	"go-web-api/internal/dto"
	"go-web-api/internal/model"
	"go-web-api/internal/service"
	"go-web-api/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CategoryHandler groups all admin category HTTP handlers.
type CategoryHandler struct {
	svc service.CategoryService
	log *zap.Logger
}

func NewCategoryHandler(svc service.CategoryService, log *zap.Logger) *CategoryHandler {
	return &CategoryHandler{svc: svc, log: log}
}

// Create godoc
// @Summary  Create a new product category
// @Tags     admin-categories
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    body body dto.CreateCategoryRequest true "Create payload"
// @Success  201  {object} utils.Response{data=dto.CategoryResponse}
// @Router   /api/admin/v1/categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, err.Error())
		return
	}

	cat, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("create category", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.SuccessCreated(c, toCategoryResponse(cat))
}

// GetByID godoc
// @Summary  Get category by ID
// @Tags     admin-categories
// @Security BearerAuth
// @Produce  json
// @Param    id  path int true "Category ID"
// @Success  200 {object} utils.Response{data=dto.CategoryResponse}
// @Router   /api/admin/v1/categories/{id} [get]
func (h *CategoryHandler) GetByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.ParamError(c, "invalid category id")
		return
	}

	cat, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("get category", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, toCategoryResponse(cat))
}

// List godoc
// @Summary  List categories with pagination
// @Tags     admin-categories
// @Security BearerAuth
// @Produce  json
// @Param    page      query int false "Page (default 1)"
// @Param    page_size query int false "Page size (default 10)"
// @Success  200       {object} utils.Response{data=dto.CategoryListResponse}
// @Router   /api/admin/v1/categories [get]
func (h *CategoryHandler) List(c *gin.Context) {
	var q dto.PageQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		utils.ParamError(c, err.Error())
		return
	}
	q.Normalize()

	cats, total, err := h.svc.List(c.Request.Context(), q.Page, q.PageSize, true)
	if err != nil {
		h.log.Error("list categories", zap.Error(err))
		utils.ServerError(c)
		return
	}

	list := make([]*dto.CategoryResponse, 0, len(cats))
	for _, cat := range cats {
		list = append(list, toCategoryResponse(cat))
	}

	utils.Success(c, dto.CategoryListResponse{Total: total, List: list})
}

// Update godoc
// @Summary  Update a category
// @Tags     admin-categories
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    id   path int true "Category ID"
// @Param    body body dto.UpdateCategoryRequest true "Update payload"
// @Success  200  {object} utils.Response{data=dto.CategoryResponse}
// @Router   /api/admin/v1/categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.ParamError(c, "invalid category id")
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, err.Error())
		return
	}

	cat, err := h.svc.Update(c.Request.Context(), id, &req)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("update category", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, toCategoryResponse(cat))
}

// Delete godoc
// @Summary  Delete a category
// @Tags     admin-categories
// @Security BearerAuth
// @Produce  json
// @Param    id  path int true "Category ID"
// @Success  200 {object} utils.Response
// @Router   /api/admin/v1/categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.ParamError(c, "invalid category id")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("delete category", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, nil)
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
