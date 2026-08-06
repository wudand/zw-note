package admin

import (
	"zw-note-backend/internal/dto"
	"zw-note-backend/internal/model"
	"zw-note-backend/internal/service"
	"zw-note-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CouponHandler groups admin coupon HTTP handlers.
type CouponHandler struct {
	svc service.CouponService
	log *zap.Logger
}

func NewCouponHandler(svc service.CouponService, log *zap.Logger) *CouponHandler {
	return &CouponHandler{svc: svc, log: log}
}

// Create godoc
// @Summary  Create a new coupon
// @Tags     admin-coupons
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    body body dto.CreateCouponRequest true "Create payload"
// @Success  201  {object} utils.Response{data=dto.CouponResponse}
// @Router   /api/admin/v1/coupons [post]
func (h *CouponHandler) Create(c *gin.Context) {
	var req dto.CreateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, err.Error())
		return
	}

	coupon, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("create coupon", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.SuccessCreated(c, toCouponResponse(coupon))
}

// GetByID returns a coupon by ID.
func (h *CouponHandler) GetByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.ParamError(c, "invalid coupon id")
		return
	}

	coupon, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("get coupon", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, toCouponResponse(coupon))
}

// List godoc
// @Summary  List coupons with pagination
// @Tags     admin-coupons
// @Security BearerAuth
// @Produce  json
// @Param    page      query int    false "Page (default 1)"
// @Param    page_size query int    false "Page size (default 10)"
// @Param    type      query string false "Filter: new_user | spend_reduce"
// @Param    status    query int    false "Filter: 0 | 1"
// @Success  200       {object} utils.Response{data=dto.CouponListResponse}
// @Router   /api/admin/v1/coupons [get]
func (h *CouponHandler) List(c *gin.Context) {
	var q dto.CouponListQuery
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

	list, total, err := h.svc.List(c.Request.Context(), q.Page, q.PageSize, q.Type, q.Status)
	if err != nil {
		h.log.Error("list coupons", zap.Error(err))
		utils.ServerError(c)
		return
	}

	resp := make([]*dto.CouponResponse, 0, len(list))
	for _, c := range list {
		resp = append(resp, toCouponResponse(c))
	}

	utils.Success(c, dto.CouponListResponse{Total: total, List: resp})
}

// Update godoc
// @Summary  Update a coupon
// @Tags     admin-coupons
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    id   path int true "Coupon ID"
// @Param    body body dto.UpdateCouponRequest true "Update payload"
// @Success  200  {object} utils.Response{data=dto.CouponResponse}
// @Router   /api/admin/v1/coupons/{id} [put]
func (h *CouponHandler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.ParamError(c, "invalid coupon id")
		return
	}

	var req dto.UpdateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, err.Error())
		return
	}

	coupon, err := h.svc.Update(c.Request.Context(), id, &req)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("update coupon", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, toCouponResponse(coupon))
}

// Cancel godoc
// @Summary  Cancel a coupon (no longer claimable)
// @Tags     admin-coupons
// @Security BearerAuth
// @Produce  json
// @Param    id  path int true "Coupon ID"
// @Success  200 {object} utils.Response{data=dto.CouponResponse}
// @Router   /api/admin/v1/coupons/{id}/cancel [put]
func (h *CouponHandler) Cancel(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.ParamError(c, "invalid coupon id")
		return
	}

	coupon, err := h.svc.Cancel(c.Request.Context(), id)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("cancel coupon", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, toCouponResponse(coupon))
}

// Delete godoc
// @Summary  Delete a coupon
// @Tags     admin-coupons
// @Security BearerAuth
// @Produce  json
// @Param    id  path int true "Coupon ID"
// @Success  200 {object} utils.Response
// @Router   /api/admin/v1/coupons/{id} [delete]
func (h *CouponHandler) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.ParamError(c, "invalid coupon id")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("delete coupon", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, nil)
}

func toCouponResponse(c *model.Coupon) *dto.CouponResponse {
	return &dto.CouponResponse{
		ID:            c.ID,
		Name:          c.Name,
		Type:          c.Type,
		MinAmount:     c.MinAmount,
		DiscountValue: c.DiscountValue,
		ValidDays:     c.ValidDays,
		Stackable:     c.Stackable == 1,
		Status:        c.Status,
		CreatedAt:     c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
