package admin

import (
	"go-web-api/internal/dto"
	"go-web-api/internal/model"
	"go-web-api/internal/service"
	"go-web-api/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ProductHandler groups all admin product HTTP handlers.
type ProductHandler struct {
	svc service.ProductService
	log *zap.Logger
}

func NewProductHandler(svc service.ProductService, log *zap.Logger) *ProductHandler {
	return &ProductHandler{svc: svc, log: log}
}

// Create godoc
// @Summary  Create a new product
// @Tags     admin-products
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    body body dto.CreateProductRequest true "Create payload"
// @Success  201  {object} utils.Response{data=dto.ProductResponse}
// @Router   /api/admin/v1/products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, err.Error())
		return
	}

	p, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("create product", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.SuccessCreated(c, toProductResponse(p))
}

// GetByID godoc
// @Summary  Get product by ID
// @Tags     admin-products
// @Security BearerAuth
// @Produce  json
// @Param    id  path int true "Product ID"
// @Success  200 {object} utils.Response{data=dto.ProductResponse}
// @Router   /api/admin/v1/products/{id} [get]
func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.ParamError(c, "invalid product id")
		return
	}

	p, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("get product", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, toProductResponse(p))
}

// List godoc
// @Summary  List products with pagination
// @Tags     admin-products
// @Security BearerAuth
// @Produce  json
// @Param    page        query int false "Page (default 1)"
// @Param    page_size   query int false "Page size (default 10)"
// @Param    category_id query int false "Filter by category ID"
// @Success  200         {object} utils.Response{data=dto.ProductListResponse}
// @Router   /api/admin/v1/products [get]
func (h *ProductHandler) List(c *gin.Context) {
	var q dto.ProductListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		utils.ParamError(c, err.Error())
		return
	}
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}

	products, total, err := h.svc.List(c.Request.Context(), q.Page, q.PageSize, q.CategoryID, true)
	if err != nil {
		h.log.Error("list products", zap.Error(err))
		utils.ServerError(c)
		return
	}

	list := make([]*dto.ProductResponse, 0, len(products))
	for _, p := range products {
		list = append(list, toProductResponse(p))
	}

	utils.Success(c, dto.ProductListResponse{Total: total, List: list})
}

// Update godoc
// @Summary  Update a product
// @Tags     admin-products
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    id   path int true "Product ID"
// @Param    body body dto.UpdateProductRequest true "Update payload"
// @Success  200  {object} utils.Response{data=dto.ProductResponse}
// @Router   /api/admin/v1/products/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.ParamError(c, "invalid product id")
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, err.Error())
		return
	}

	p, err := h.svc.Update(c.Request.Context(), id, &req)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("update product", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, toProductResponse(p))
}

// Delete godoc
// @Summary  Delete a product
// @Tags     admin-products
// @Security BearerAuth
// @Produce  json
// @Param    id  path int true "Product ID"
// @Success  200 {object} utils.Response
// @Router   /api/admin/v1/products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.ParamError(c, "invalid product id")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("delete product", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, nil)
}

func toProductResponse(p *model.Product) *dto.ProductResponse {
	carousel := []string(p.CarouselImages)
	detailImgs := []string(p.DetailImages)
	if carousel == nil {
		carousel = []string{}
	}
	if detailImgs == nil {
		detailImgs = []string{}
	}
	return &dto.ProductResponse{
		ID:             p.ID,
		CategoryID:     p.CategoryID,
		Name:           p.Name,
		Ingredients:    p.Ingredients,
		Grade:          p.Grade,
		Storage:        p.Storage,
		Specification:  p.Specification,
		CoverImage:     p.CoverImage,
		CarouselImages: carousel,
		Detail:         p.Detail,
		DetailImages:   detailImgs,
		Status:         p.Status,
		CreatedAt:      p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
