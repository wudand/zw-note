package model

import "time"

// MPUser maps to the `mp_users` table (WeChat mini-program users).
// One row per (app_id, openid) pair to support multiple mini-programs.
type MPUser struct {
	ID        uint64    `db:"id"`
	AppID     string    `db:"app_id"`  // WeChat AppID of the mini-program
	OpenID    string    `db:"openid"`  // unique per user per AppID
	UnionID   string    `db:"unionid"` // cross-app unique ID (may be empty)
	Nickname  string    `db:"nickname"`
	Avatar    string    `db:"avatar"`
	Status    int8      `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

const (
	MPUserStatusActive   int8 = 1
	MPUserStatusInactive int8 = 0
)
