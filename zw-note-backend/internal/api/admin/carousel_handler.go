package admin

import (
	"zw-note-backend/internal/dto"
	"zw-note-backend/internal/model"
	"zw-note-backend/internal/service"
	"zw-note-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CarouselHandler groups admin carousel HTTP handlers.
type CarouselHandler struct {
	svc service.CarouselService
	log *zap.Logger
}

func NewCarouselHandler(svc service.CarouselService, log *zap.Logger) *CarouselHandler {
	return &CarouselHandler{svc: svc, log: log}
}

// Create godoc
// @Summary  Create a new carousel
// @Tags     admin-carousels
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    body body dto.CreateCarouselRequest true "Create payload"
// @Success  201  {object} utils.Response{data=dto.CarouselResponse}
// @Router   /api/admin/v1/carousels [post]
func (h *CarouselHandler) Create(c *gin.Context) {
	var req dto.CreateCarouselRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, err.Error())
		return
	}

	carousel, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("create carousel", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.SuccessCreated(c, toCarouselResponse(carousel))
}

// GetByID godoc
// @Summary  Get carousel by ID
// @Tags     admin-carousels
// @Security BearerAuth
// @Produce  json
// @Param    id  path int true "Carousel ID"
// @Success  200 {object} utils.Response{data=dto.CarouselResponse}
// @Router   /api/admin/v1/carousels/{id} [get]
func (h *CarouselHandler) GetByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.ParamError(c, "invalid carousel id")
		return
	}

	carousel, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("get carousel", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, toCarouselResponse(carousel))
}

// List godoc
// @Summary  List all carousels
// @Tags     admin-carousels
// @Security BearerAuth
// @Produce  json
// @Success  200 {object} utils.Response{data=dto.CarouselListResponse}
// @Router   /api/admin/v1/carousels [get]
func (h *CarouselHandler) List(c *gin.Context) {
	list, err := h.svc.List(c.Request.Context(), false)
	if err != nil {
		h.log.Error("list carousels", zap.Error(err))
		utils.ServerError(c)
		return
	}

	resp := make([]*dto.CarouselResponse, 0, len(list))
	for _, carousel := range list {
		resp = append(resp, toCarouselResponse(carousel))
	}

	utils.Success(c, dto.CarouselListResponse{List: resp})
}

// Update godoc
// @Summary  Update a carousel
// @Tags     admin-carousels
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    id   path int true "Carousel ID"
// @Param    body body dto.UpdateCarouselRequest true "Update payload"
// @Success  200  {object} utils.Response{data=dto.CarouselResponse}
// @Router   /api/admin/v1/carousels/{id} [put]
func (h *CarouselHandler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.ParamError(c, "invalid carousel id")
		return
	}

	var req dto.UpdateCarouselRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, err.Error())
		return
	}

	carousel, err := h.svc.Update(c.Request.Context(), id, &req)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("update carousel", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, toCarouselResponse(carousel))
}

// Delete godoc
// @Summary  Delete a carousel
// @Tags     admin-carousels
// @Security BearerAuth
// @Produce  json
// @Param    id  path int true "Carousel ID"
// @Success  200 {object} utils.Response
// @Router   /api/admin/v1/carousels/{id} [delete]
func (h *CarouselHandler) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.ParamError(c, "invalid carousel id")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("delete carousel", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, nil)
}

func toCarouselResponse(c *model.Carousel) *dto.CarouselResponse {
	return &dto.CarouselResponse{
		ID:        c.ID,
		ImageURL:  c.ImageURL,
		Title:     c.Title,
		Link:      c.Link,
		SortOrder: c.SortOrder,
		Status:    c.Status,
		CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
