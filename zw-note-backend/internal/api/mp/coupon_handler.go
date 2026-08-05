package mp

import (
	"strconv"

	"go-web-api/internal/dto"
	"go-web-api/internal/model"
	"go-web-api/internal/service"
	"go-web-api/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CouponHandler groups mini-program coupon HTTP handlers.
type CouponHandler struct {
	svc service.CouponService
	log *zap.Logger
}

func NewCouponHandler(svc service.CouponService, log *zap.Logger) *CouponHandler {
	return &CouponHandler{svc: svc, log: log}
}

// ListClaimable godoc
// @Summary  List claimable coupons (public, no auth)
// @Tags     mp-coupons
// @Produce  json
// @Success  200 {object} utils.Response{data=dto.CouponListResponse}
// @Router   /api/mp/v1/coupons [get]
func (h *CouponHandler) ListClaimable(c *gin.Context) {
	list, err := h.svc.ListClaimable(c.Request.Context())
	if err != nil {
		h.log.Error("list claimable coupons", zap.Error(err))
		utils.ServerError(c)
		return
	}

	resp := make([]*dto.CouponResponse, 0, len(list))
	for _, c := range list {
		resp = append(resp, toCouponResponse(c))
	}
	utils.Success(c, dto.CouponListResponse{Total: int64(len(resp)), List: resp})
}

// Claim godoc
// @Summary  Claim a coupon
// @Tags     mp-coupons
// @Security BearerAuth
// @Produce  json
// @Param    id  path int true "Coupon ID"
// @Success  201 {object} utils.Response{data=dto.UserCouponResponse}
// @Router   /api/mp/v1/coupons/{id}/claim [post]
func (h *CouponHandler) Claim(c *gin.Context) {
	mpUserID, ok := getMPUserID(c)
	if !ok {
		utils.Unauthorized(c)
		return
	}

	couponID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ParamError(c, "invalid coupon id")
		return
	}

	uc, err := h.svc.Claim(c.Request.Context(), mpUserID, couponID)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("claim coupon", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.SuccessCreated(c, toUserCouponResponse(uc))
}

// ListMy godoc
// @Summary  List current user's coupons
// @Tags     mp-coupons
// @Security BearerAuth
// @Produce  json
// @Param    status query string false "Filter: unused | used | expired"
// @Success  200    {object} utils.Response{data=dto.UserCouponListResponse}
// @Router   /api/mp/v1/coupons/my [get]
func (h *CouponHandler) ListMy(c *gin.Context) {
	mpUserID, ok := getMPUserID(c)
	if !ok {
		utils.Unauthorized(c)
		return
	}

	status := c.Query("status") // unused | used | expired, empty = all
	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	list, err := h.svc.ListMyCoupons(c.Request.Context(), mpUserID, statusPtr)
	if err != nil {
		h.log.Error("list my coupons", zap.Error(err))
		utils.ServerError(c)
		return
	}

	resp := make([]*dto.UserCouponResponse, 0, len(list))
	for _, uc := range list {
		resp = append(resp, toUserCouponResponse(uc))
	}
	utils.Success(c, dto.UserCouponListResponse{List: resp})
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

func toUserCouponResponse(uc *model.UserCoupon) *dto.UserCouponResponse {
	resp := &dto.UserCouponResponse{
		ID:            uc.ID,
		CouponID:      uc.CouponID,
		CouponName:    uc.CouponName,
		CouponType:    uc.CouponType,
		MinAmount:     uc.MinAmount,
		DiscountValue: uc.DiscountValue,
		Stackable:     uc.Stackable == 1,
		Status:        uc.Status,
		ClaimedAt:     uc.ClaimedAt.Format("2006-01-02T15:04:05Z07:00"),
		ExpiryAt:      uc.ExpiryAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if uc.UsedAt != nil {
		s := uc.UsedAt.Format("2006-01-02T15:04:05Z07:00")
		resp.UsedAt = &s
	}
	return resp
}
