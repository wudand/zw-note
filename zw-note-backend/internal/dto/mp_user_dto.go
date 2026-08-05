package dto

// WxLoginRequest is the payload for POST /api/mp/v1/auth/wx-login.
type WxLoginRequest struct {
	AppKey string `json:"app_key" binding:"required"` // key in wechat.apps config
	Code   string `json:"code"    binding:"required"` // wx.login() code from mini-program
}

// MPLoginResponse carries the mini-program JWT and its expiry.
type MPLoginResponse struct {
	Token    string `json:"token"`
	ExpireAt int64  `json:"expire_at"`
	IsNew    bool   `json:"is_new"` // true if account was just created
}

// MPUserResponse is the public-facing representation of a mini-program user.
type MPUserResponse struct {
	ID        uint64 `json:"id"`
	OpenID    string `json:"openid"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Status    int8   `json:"status"`
	CreatedAt string `json:"created_at"`
}
