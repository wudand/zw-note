package service

import (
	"context"
	"fmt"

	"zw-note-backend/internal/config"
	"zw-note-backend/internal/model"
	"zw-note-backend/internal/repository"
	"zw-note-backend/pkg/utils"

	"go.uber.org/zap"
)

//go:generate mockgen -destination=mock/mp_user_service_mock.go -package=mock zw-note-backend/internal/service MPUserService

// MPUserService defines business logic for WeChat mini-program user operations.
type MPUserService interface {
	// WxLogin performs the full WeChat login flow:
	//   1. Exchange code for openid via WeChat API
	//   2. Upsert the mp_user record
	//   3. Return the user and whether the account was just created
	WxLogin(ctx context.Context, appKey, code string) (user *model.MPUser, isNew bool, err error)
	GetByID(ctx context.Context, id uint64) (*model.MPUser, error)
}

type mpUserService struct {
	repo  repository.MPUserRepository
	wxSvc WxService
	wxCfg config.WechatConfig
	log   *zap.Logger
}

func NewMPUserService(
	repo repository.MPUserRepository,
	wxSvc WxService,
	wxCfg config.WechatConfig,
	log *zap.Logger,
) MPUserService {
	return &mpUserService{repo: repo, wxSvc: wxSvc, wxCfg: wxCfg, log: log}
}

func (s *mpUserService) WxLogin(ctx context.Context, appKey, code string) (*model.MPUser, bool, error) {
	// Resolve AppID from appKey
	app, err := s.wxCfg.GetApp(appKey)
	if err != nil {
		return nil, false, utils.ErrInvalidAppKey
	}

	// Exchange code for openid
	session, err := s.wxSvc.Code2Session(ctx, appKey, code)
	if err != nil {
		return nil, false, fmt.Errorf("code2session: %w", err)
	}

	// Upsert user
	user, err := s.repo.GetByOpenID(ctx, app.AppID, session.OpenID)
	if err != nil {
		return nil, false, fmt.Errorf("get mp user: %w", err)
	}

	if user != nil {
		if user.Status == MPUserStatusInactive {
			return nil, false, utils.ErrUserDisabled
		}
		return user, false, nil
	}

	// First-time login – auto-register
	user = &model.MPUser{
		AppID:   app.AppID,
		OpenID:  session.OpenID,
		UnionID: session.UnionID,
		Status:  model.MPUserStatusActive,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, false, err
	}

	s.log.Info("mp user registered", zap.String("openid", user.OpenID), zap.String("app_id", user.AppID))
	return user, true, nil
}

func (s *mpUserService) GetByID(ctx context.Context, id uint64) (*model.MPUser, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, utils.ErrMPUserNotFound
	}
	return user, nil
}

// MPUserStatusInactive is a local alias to avoid cross-package import from the service layer.
const MPUserStatusInactive int8 = 0
