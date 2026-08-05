package mp

import (
	"go-web-api/internal/dto"
	"go-web-api/internal/model"
	"go-web-api/internal/service"
	"go-web-api/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CarouselHandler groups mini-program carousel HTTP handlers (display only).
type CarouselHandler struct {
	svc service.CarouselService
	log *zap.Logger
}

func NewCarouselHandler(svc service.CarouselService, log *zap.Logger) *CarouselHandler {
	return &CarouselHandler{svc: svc, log: log}
}

// List godoc
// @Summary  List enabled carousels for home page (public, no auth)
// @Tags     mp-carousels
// @Produce  json
// @Success  200 {object} utils.Response{data=dto.CarouselListResponse}
// @Router   /api/mp/v1/carousels [get]
func (h *CarouselHandler) List(c *gin.Context) {
	list, err := h.svc.List(c.Request.Context(), true)
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
