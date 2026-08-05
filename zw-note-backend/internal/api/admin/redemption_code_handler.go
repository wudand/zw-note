package admin

import (
	"go-web-api/internal/dto"
	"go-web-api/internal/model"
	"go-web-api/internal/service"
	"go-web-api/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RedemptionCodeHandler groups admin redemption code HTTP handlers.
type RedemptionCodeHandler struct {
	svc service.RedemptionCodeService
	log *zap.Logger
}

func NewRedemptionCodeHandler(svc service.RedemptionCodeService, log *zap.Logger) *RedemptionCodeHandler {
	return &RedemptionCodeHandler{svc: svc, log: log}
}

// CreateBatch godoc
// @Summary  Create redemption codes in batch (max 50)
// @Tags     admin-redemption-codes
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    body body dto.CreateRedemptionCodesRequest true "Create payload"
// @Success  201  {object} utils.Response{data=dto.RedemptionCodeListResponse}
// @Router   /api/admin/v1/redemption-codes [post]
func (h *RedemptionCodeHandler) CreateBatch(c *gin.Context) {
	var req dto.CreateRedemptionCodesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, err.Error())
		return
	}

	codes, err := h.svc.CreateBatch(c.Request.Context(), &req)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("create redemption codes", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	resp := make([]*dto.RedemptionCodeResponse, 0, len(codes))
	for _, rc := range codes {
		resp = append(resp, toRedemptionCodeResponse(rc, req.ProductIDs))
	}
	utils.SuccessCreated(c, dto.RedemptionCodeListResponse{Total: int64(len(resp)), List: resp})
}

// GetByID godoc
// @Summary  Get redemption code by ID
// @Tags     admin-redemption-codes
// @Security BearerAuth
// @Produce  json
// @Param    id  path int true "Redemption Code ID"
// @Success  200 {object} utils.Response{data=dto.RedemptionCodeResponse}
// @Router   /api/admin/v1/redemption-codes/{id} [get]
func (h *RedemptionCodeHandler) GetByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.ParamError(c, "invalid redemption code id")
		return
	}

	rc, pids, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("get redemption code", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, toRedemptionCodeResponse(rc, pids))
}

// List godoc
// @Summary  List redemption codes with pagination
// @Tags     admin-redemption-codes
// @Security BearerAuth
// @Produce  json
// @Param    page      query int false "Page (default 1)"
// @Param    page_size query int false "Page size (default 10)"
// @Param    status    query int false "Filter: 0 | 1"
// @Success  200       {object} utils.Response{data=dto.RedemptionCodeListResponse}
// @Router   /api/admin/v1/redemption-codes [get]
func (h *RedemptionCodeHandler) List(c *gin.Context) {
	var q dto.RedemptionCodeListQuery
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

	list, productMap, total, err := h.svc.List(c.Request.Context(), q.Page, q.PageSize, q.Status)
	if err != nil {
		h.log.Error("list redemption codes", zap.Error(err))
		utils.ServerError(c)
		return
	}

	resp := make([]*dto.RedemptionCodeResponse, 0, len(list))
	for _, rc := range list {
		pids := productMap[rc.ID]
		resp = append(resp, toRedemptionCodeResponse(rc, pids))
	}

	utils.Success(c, dto.RedemptionCodeListResponse{Total: total, List: resp})
}

// Update godoc
// @Summary  Update redemption code product bindings (only when unused)
// @Tags     admin-redemption-codes
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    id   path int true "Redemption Code ID"
// @Param    body body dto.UpdateRedemptionCodeRequest true "Update payload"
// @Success  200  {object} utils.Response{data=dto.RedemptionCodeResponse}
// @Router   /api/admin/v1/redemption-codes/{id} [put]
func (h *RedemptionCodeHandler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.ParamError(c, "invalid redemption code id")
		return
	}

	var req dto.UpdateRedemptionCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, err.Error())
		return
	}

	rc, err := h.svc.Update(c.Request.Context(), id, &req)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("update redemption code", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, toRedemptionCodeResponse(rc, req.ProductIDs))
}

// Delete godoc
// @Summary  Delete a redemption code (only when unused)
// @Tags     admin-redemption-codes
// @Security BearerAuth
// @Produce  json
// @Param    id  path int true "Redemption Code ID"
// @Success  200 {object} utils.Response
// @Router   /api/admin/v1/redemption-codes/{id} [delete]
func (h *RedemptionCodeHandler) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.ParamError(c, "invalid redemption code id")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("delete redemption code", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, nil)
}

func toRedemptionCodeResponse(rc *model.RedemptionCode, pids []uint64) *dto.RedemptionCodeResponse {
	resp := &dto.RedemptionCodeResponse{
		ID:        rc.ID,
		Code:      rc.Code,
		Status:    rc.Status,
		ProductIDs: pids,
		CreatedAt: rc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: rc.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if rc.MPUserID != nil {
		resp.MPUserID = rc.MPUserID
	}
	if rc.UsedProductID != nil {
		resp.UsedProductID = rc.UsedProductID
	}
	if rc.UsedAt != nil {
		s := rc.UsedAt.Format("2006-01-02T15:04:05Z07:00")
		resp.UsedAt = &s
	}
	return resp
}
