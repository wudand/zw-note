package mp

import (
	"go-web-api/internal/config"
	"go-web-api/internal/dto"
	"go-web-api/internal/model"
	"go-web-api/internal/service"
	"go-web-api/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthHandler handles WeChat mini-program authentication.
type AuthHandler struct {
	mpSvc  service.MPUserService
	log    *zap.Logger
	jwtCfg config.JWTConfig
}

func NewAuthHandler(mpSvc service.MPUserService, log *zap.Logger, jwtCfg config.JWTConfig) *AuthHandler {
	return &AuthHandler{mpSvc: mpSvc, log: log, jwtCfg: jwtCfg}
}

// WxLogin godoc
// @Summary  WeChat mini-program login
// @Description  Exchange a wx.login() code for a project JWT.
//
//	If the user is new, an mp_user record is created automatically.
//
// @Tags     mp-auth
// @Accept   json
// @Produce  json
// @Param    body body dto.WxLoginRequest true "WeChat login payload"
// @Success  200  {object} utils.Response{data=dto.MPLoginResponse}
// @Router   /api/mp/v1/auth/wx-login [post]
func (h *AuthHandler) WxLogin(c *gin.Context) {
	var req dto.WxLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, err.Error())
		return
	}

	user, isNew, err := h.mpSvc.WxLogin(c.Request.Context(), req.AppKey, req.Code)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("wx login", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	token, expireAt, err := utils.GenerateMPToken(
		user.ID, user.OpenID, user.AppID,
		h.jwtCfg.MPSecret, h.jwtCfg.MPExpireHours,
	)
	if err != nil {
		h.log.Error("generate mp token", zap.Error(err))
		utils.ServerError(c)
		return
	}

	utils.Success(c, dto.MPLoginResponse{Token: token, ExpireAt: expireAt, IsNew: isNew})
}

// GetProfile godoc
// @Summary  Get the current mini-program user's profile
// @Tags     mp-user
// @Security BearerAuth
// @Produce  json
// @Success  200 {object} utils.Response{data=dto.MPUserResponse}
// @Router   /api/mp/v1/user/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, ok := userID.(uint64)
	if !ok {
		utils.Unauthorized(c)
		return
	}

	user, err := h.mpSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("get mp user profile", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, toMPUserResponse(user))
}

func toMPUserResponse(u *model.MPUser) *dto.MPUserResponse {
	return &dto.MPUserResponse{
		ID:        u.ID,
		OpenID:    u.OpenID,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Status:    u.Status,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
