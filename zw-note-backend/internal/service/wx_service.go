package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"zw-note-backend/internal/config"
	"zw-note-backend/pkg/utils"
)

const wxCode2SessionURL = "https://api.weixin.qq.com/sns/jscode2session"

// WxSession is the response from the WeChat code2session API.
type WxSession struct {
	OpenID     string `json:"openid"`
	UnionID    string `json:"unionid"`
	SessionKey string `json:"session_key"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// WxService exchanges a wx.login() code for an OpenID via the WeChat API.
type WxService interface {
	Code2Session(ctx context.Context, appKey, code string) (*WxSession, error)
}

type wxService struct {
	cfg    config.WechatConfig
	client *http.Client
}

func NewWxService(cfg config.WechatConfig) WxService {
	return &wxService{
		cfg:    cfg,
		client: &http.Client{Timeout: 10 * 1e9}, // 10s
	}
}

// Code2Session calls the WeChat server to exchange the mini-program login code for
// an openid (and optional unionid). It uses the app credentials registered under appKey.
func (s *wxService) Code2Session(ctx context.Context, appKey, code string) (*WxSession, error) {
	app, err := s.cfg.GetApp(appKey)
	if err != nil {
		return nil, utils.ErrInvalidAppKey
	}

	params := url.Values{}
	params.Set("appid", app.AppID)
	params.Set("secret", app.AppSecret)
	params.Set("js_code", code)
	params.Set("grant_type", "authorization_code")

	reqURL := wxCode2SessionURL + "?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build wx request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call wx api: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read wx response: %w", err)
	}

	var session WxSession
	if err := json.Unmarshal(body, &session); err != nil {
		return nil, fmt.Errorf("unmarshal wx response: %w", err)
	}

	if session.ErrCode != 0 {
		return nil, fmt.Errorf("%w: errcode=%d errmsg=%s", utils.ErrWxLoginFailed, session.ErrCode, session.ErrMsg)
	}

	return &session, nil
}
