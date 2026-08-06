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

// RedemptionCodeHandler groups mini-program redemption code HTTP handlers.
type RedemptionCodeHandler struct {
	svc service.RedemptionCodeService
	log *zap.Logger
}

func NewRedemptionCodeHandler(svc service.RedemptionCodeService, log *zap.Logger) *RedemptionCodeHandler {
	return &RedemptionCodeHandler{svc: svc, log: log}
}

// Validate godoc
// @Summary  Validate redemption code and lock to current user (returns selectable products)
// @Tags     mp-redemption-codes
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    body body dto.ValidateRedemptionCodeRequest true "Validate payload"
// @Success  200  {object} utils.Response{data=dto.ValidateRedemptionCodeResponse}
// @Router   /api/mp/v1/redemption-codes/validate [post]
func (h *RedemptionCodeHandler) Validate(c *gin.Context) {
	mpUserID, ok := getMPUserID(c)
	if !ok {
		utils.Unauthorized(c)
		return
	}

	var req dto.ValidateRedemptionCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, err.Error())
		return
	}

	rc, products, err := h.svc.Validate(c.Request.Context(), mpUserID, req.Code)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("validate redemption code", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	resp := &dto.ValidateRedemptionCodeResponse{
		RedemptionCodeID: rc.ID,
		Products:         toProductSummaries(products),
	}
	utils.Success(c, resp)
}

// SelectProduct godoc
// @Summary  Select product for redemption code (call once only)
// @Tags     mp-redemption-codes
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    id   path int true "Redemption Code ID"
// @Param    body body dto.SelectProductRequest true "Select payload"
// @Success  200  {object} utils.Response
// @Router   /api/mp/v1/redemption-codes/{id}/select-product [post]
func (h *RedemptionCodeHandler) SelectProduct(c *gin.Context) {
	mpUserID, ok := getMPUserID(c)
	if !ok {
		utils.Unauthorized(c)
		return
	}

	redemptionCodeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ParamError(c, "invalid redemption code id")
		return
	}

	var req dto.SelectProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, err.Error())
		return
	}

	if err := h.svc.SelectProduct(c.Request.Context(), mpUserID, redemptionCodeID, req.ProductID); err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("select product for redemption code", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, gin.H{"redemption_code_id": redemptionCodeID, "product_id": req.ProductID})
}

func toProductSummaries(products []*model.Product) []*dto.ProductSummary {
	if products == nil {
		return nil
	}
	resp := make([]*dto.ProductSummary, 0, len(products))
	for _, p := range products {
		resp = append(resp, &dto.ProductSummary{
			ID:            p.ID,
			Name:          p.Name,
			CoverImage:    p.CoverImage,
			Specification: p.Specification,
		})
	}
	return resp
}
