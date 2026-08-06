package admin

import (
	"strconv"

	"zw-note-backend/internal/dto"
	"zw-note-backend/internal/model"
	"zw-note-backend/internal/service"
	"zw-note-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AddressHandler groups admin address HTTP handlers (manage mp user addresses).
type AddressHandler struct {
	svc service.AddressService
	log *zap.Logger
}

func NewAddressHandler(svc service.AddressService, log *zap.Logger) *AddressHandler {
	return &AddressHandler{svc: svc, log: log}
}

func parseMPUserID(c *gin.Context) (uint64, error) {
	return strconv.ParseUint(c.Param("mp_user_id"), 10, 64)
}

// List godoc
// @Summary  List addresses for an MP user
// @Tags     admin-addresses
// @Security BearerAuth
// @Produce  json
// @Param    mp_user_id path int true "MP User ID"
// @Success  200        {object} utils.Response{data=dto.AddressListResponse}
// @Router   /api/admin/v1/mp-users/{mp_user_id}/addresses [get]
func (h *AddressHandler) List(c *gin.Context) {
	mpUserID, err := parseMPUserID(c)
	if err != nil {
		utils.ParamError(c, "invalid mp_user_id")
		return
	}

	list, err := h.svc.ListByUser(c.Request.Context(), mpUserID)
	if err != nil {
		h.log.Error("list addresses", zap.Error(err))
		utils.ServerError(c)
		return
	}

	resp := make([]*dto.AddressResponse, 0, len(list))
	for _, a := range list {
		resp = append(resp, toAddressResponse(a))
	}
	utils.Success(c, dto.AddressListResponse{List: resp})
}

// GetByID godoc
// @Summary  Get address by ID
// @Tags     admin-addresses
// @Security BearerAuth
// @Produce  json
// @Param    mp_user_id path int true "MP User ID"
// @Param    id         path int true "Address ID"
// @Success  200        {object} utils.Response{data=dto.AddressResponse}
// @Router   /api/admin/v1/mp-users/{mp_user_id}/addresses/{id} [get]
func (h *AddressHandler) GetByID(c *gin.Context) {
	mpUserID, err := parseMPUserID(c)
	if err != nil {
		utils.ParamError(c, "invalid mp_user_id")
		return
	}

	id, err := parseID(c)
	if err != nil {
		utils.ParamError(c, "invalid address id")
		return
	}

	a, err := h.svc.GetByID(c.Request.Context(), id, mpUserID)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("get address", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, toAddressResponse(a))
}

// Create godoc
// @Summary  Create address for an MP user
// @Tags     admin-addresses
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    mp_user_id path int true "MP User ID"
// @Param    body       body dto.CreateAddressRequest true "Create payload"
// @Success  201        {object} utils.Response{data=dto.AddressResponse}
// @Router   /api/admin/v1/mp-users/{mp_user_id}/addresses [post]
func (h *AddressHandler) Create(c *gin.Context) {
	mpUserID, err := parseMPUserID(c)
	if err != nil {
		utils.ParamError(c, "invalid mp_user_id")
		return
	}

	var req dto.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, err.Error())
		return
	}

	a, err := h.svc.Create(c.Request.Context(), mpUserID, &req)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("create address", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.SuccessCreated(c, toAddressResponse(a))
}

// Update godoc
// @Summary  Update an address
// @Tags     admin-addresses
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    mp_user_id path int true "MP User ID"
// @Param    id         path int true "Address ID"
// @Param    body       body dto.UpdateAddressRequest true "Update payload"
// @Success  200        {object} utils.Response{data=dto.AddressResponse}
// @Router   /api/admin/v1/mp-users/{mp_user_id}/addresses/{id} [put]
func (h *AddressHandler) Update(c *gin.Context) {
	mpUserID, err := parseMPUserID(c)
	if err != nil {
		utils.ParamError(c, "invalid mp_user_id")
		return
	}

	id, err := parseID(c)
	if err != nil {
		utils.ParamError(c, "invalid address id")
		return
	}

	var req dto.UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, err.Error())
		return
	}

	a, err := h.svc.Update(c.Request.Context(), id, mpUserID, &req)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("update address", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, toAddressResponse(a))
}

// SetDefault godoc
// @Summary  Set address as default
// @Tags     admin-addresses
// @Security BearerAuth
// @Produce  json
// @Param    mp_user_id path int true "MP User ID"
// @Param    id         path int true "Address ID"
// @Success  200        {object} utils.Response{data=dto.AddressResponse}
// @Router   /api/admin/v1/mp-users/{mp_user_id}/addresses/{id}/default [put]
func (h *AddressHandler) SetDefault(c *gin.Context) {
	mpUserID, err := parseMPUserID(c)
	if err != nil {
		utils.ParamError(c, "invalid mp_user_id")
		return
	}

	id, err := parseID(c)
	if err != nil {
		utils.ParamError(c, "invalid address id")
		return
	}

	a, err := h.svc.SetDefault(c.Request.Context(), id, mpUserID)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("set default address", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, toAddressResponse(a))
}

// Delete godoc
// @Summary  Delete an address
// @Tags     admin-addresses
// @Security BearerAuth
// @Produce  json
// @Param    mp_user_id path int true "MP User ID"
// @Param    id         path int true "Address ID"
// @Success  200        {object} utils.Response
// @Router   /api/admin/v1/mp-users/{mp_user_id}/addresses/{id} [delete]
func (h *AddressHandler) Delete(c *gin.Context) {
	mpUserID, err := parseMPUserID(c)
	if err != nil {
		utils.ParamError(c, "invalid mp_user_id")
		return
	}

	id, err := parseID(c)
	if err != nil {
		utils.ParamError(c, "invalid address id")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id, mpUserID); err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("delete address", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, nil)
}

func toAddressResponse(a *model.Address) *dto.AddressResponse {
	return &dto.AddressResponse{
		ID:        a.ID,
		MPUserID:  a.MPUserID,
		Receiver:  a.Receiver,
		Phone:     a.Phone,
		Province:  a.Province,
		City:      a.City,
		District:  a.District,
		Detail:    a.Detail,
		Tag:       a.Tag,
		IsDefault: a.IsDefault == model.AddressDefault,
		CreatedAt: a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: a.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
