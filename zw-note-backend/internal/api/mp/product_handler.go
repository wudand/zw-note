package mp

import (
	"strconv"

	"zw-note-backend/internal/dto"
	"zw-note-backend/internal/model"
	"zw-note-backend/internal/service"
	"zw-note-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func parseID(c *gin.Context) (uint64, error) {
	return strconv.ParseUint(c.Param("id"), 10, 64)
}

// ProductHandler groups all mini-program product HTTP handlers (read-only, public).
type ProductHandler struct {
	svc service.ProductService
	log *zap.Logger
}

func NewProductHandler(svc service.ProductService, log *zap.Logger) *ProductHandler {
	return &ProductHandler{svc: svc, log: log}
}

// GetByID godoc
// @Summary  Get product by ID (enabled only, public)
// @Tags     mp-products
// @Produce  json
// @Param    id  path int true "Product ID"
// @Success  200 {object} utils.Response{data=dto.ProductResponse}
// @Router   /api/mp/v1/products/{id} [get]
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
	if p.Status == 0 {
		utils.HandleAppError(c, utils.ErrProductNotFound)
		return
	}

	utils.Success(c, toProductResponse(p))
}

// List godoc
// @Summary  List products with pagination (enabled only, public)
// @Tags     mp-products
// @Produce  json
// @Param    page        query int false "Page (default 1)"
// @Param    page_size   query int false "Page size (default 10)"
// @Param    category_id query int false "Filter by category ID"
// @Success  200         {object} utils.Response{data=dto.ProductListResponse}
// @Router   /api/mp/v1/products [get]
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

	products, total, err := h.svc.List(c.Request.Context(), q.Page, q.PageSize, q.CategoryID, false)
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
