package admin

import (
	"strconv"

	"go-web-api/internal/config"
	"go-web-api/internal/dto"
	"go-web-api/internal/model"
	"go-web-api/internal/service"
	"go-web-api/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserHandler groups all admin user HTTP handlers.
type UserHandler struct {
	svc    service.AdminUserService
	log    *zap.Logger
	jwtCfg config.JWTConfig
}

func NewUserHandler(svc service.AdminUserService, log *zap.Logger, jwtCfg config.JWTConfig) *UserHandler {
	return &UserHandler{svc: svc, log: log, jwtCfg: jwtCfg}
}

// Login godoc
// @Summary  Admin login
// @Tags     admin-auth
// @Accept   json
// @Produce  json
// @Param    body body dto.AdminLoginRequest true "Login payload"
// @Success  200  {object} utils.Response{data=dto.AdminLoginResponse}
// @Router   /api/admin/v1/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, err.Error())
		return
	}

	user, err := h.svc.Authenticate(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("admin login", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	token, expireAt, err := utils.GenerateAdminToken(
		user.ID, user.Username, user.Role,
		h.jwtCfg.AdminSecret, h.jwtCfg.AdminExpireHours,
	)
	if err != nil {
		h.log.Error("generate admin token", zap.Error(err))
		utils.ServerError(c)
		return
	}

	utils.Success(c, dto.AdminLoginResponse{Token: token, ExpireAt: expireAt, Role: user.Role})
}

// CreateUser godoc
// @Summary  Create a new admin user
// @Tags     admin-users
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    body body dto.CreateAdminUserRequest true "Create payload"
// @Success  201  {object} utils.Response{data=dto.AdminUserResponse}
// @Router   /api/admin/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateAdminUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, err.Error())
		return
	}

	user, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("create admin user", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.SuccessCreated(c, toAdminUserResponse(user))
}

// GetUser godoc
// @Summary  Get admin user by ID
// @Tags     admin-users
// @Security BearerAuth
// @Produce  json
// @Param    id  path     int true "Admin User ID"
// @Success  200 {object} utils.Response{data=dto.AdminUserResponse}
// @Router   /api/admin/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.ParamError(c, "invalid user id")
		return
	}

	user, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("get admin user", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, toAdminUserResponse(user))
}

// ListUsers godoc
// @Summary  List admin users with pagination
// @Tags     admin-users
// @Security BearerAuth
// @Produce  json
// @Param    page      query int false "Page (default 1)"
// @Param    page_size query int false "Page size (default 10)"
// @Success  200       {object} utils.Response{data=dto.AdminUserListResponse}
// @Router   /api/admin/v1/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	var q dto.PageQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		utils.ParamError(c, err.Error())
		return
	}
	q.Normalize()

	users, total, err := h.svc.List(c.Request.Context(), q.Page, q.PageSize)
	if err != nil {
		h.log.Error("list admin users", zap.Error(err))
		utils.ServerError(c)
		return
	}

	list := make([]*dto.AdminUserResponse, 0, len(users))
	for _, u := range users {
		list = append(list, toAdminUserResponse(u))
	}

	utils.Success(c, dto.AdminUserListResponse{Total: total, List: list})
}

// UpdateUser godoc
// @Summary  Update an admin user
// @Tags     admin-users
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    id   path     int                        true "Admin User ID"
// @Param    body body     dto.UpdateAdminUserRequest true "Update payload"
// @Success  200  {object} utils.Response{data=dto.AdminUserResponse}
// @Router   /api/admin/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.ParamError(c, "invalid user id")
		return
	}

	var req dto.UpdateAdminUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, err.Error())
		return
	}

	user, err := h.svc.Update(c.Request.Context(), id, &req)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("update admin user", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, toAdminUserResponse(user))
}

// DeleteUser godoc
// @Summary  Delete an admin user (super_admin only)
// @Tags     admin-users
// @Security BearerAuth
// @Produce  json
// @Param    id  path     int true "Admin User ID"
// @Success  200 {object} utils.Response
// @Router   /api/admin/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		utils.ParamError(c, "invalid user id")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("delete admin user", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, nil)
}

func toAdminUserResponse(u *model.AdminUser) *dto.AdminUserResponse {
	return &dto.AdminUserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Role:      u.Role,
		Status:    u.Status,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func parseID(c *gin.Context) (uint64, error) {
	return strconv.ParseUint(c.Param("id"), 10, 64)
}
