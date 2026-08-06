package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AdminClaims is embedded in JWTs issued to management-console users.
// The role field enables stateless RBAC checks without a database round-trip.
type AdminClaims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"` // "super_admin" | "admin" | "viewer"
	jwt.RegisteredClaims
}

// MPClaims is embedded in JWTs issued to WeChat mini-program users.
type MPClaims struct {
	UserID uint64 `json:"user_id"`
	OpenID string `json:"openid"`
	AppID  string `json:"app_id"` // which mini-program
	jwt.RegisteredClaims
}

// GenerateAdminToken signs a new HS256 JWT for a management-console user.
// Returns the token string and its Unix expiry timestamp.
func GenerateAdminToken(userID uint64, username, role, secret string, expireHours int) (string, int64, error) {
	expireAt := time.Now().Add(time.Duration(expireHours) * time.Hour)
	claims := AdminClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "zw-note-backend:admin",
		},
	}
	return signToken(claims, secret, expireAt.Unix())
}

// GenerateMPToken signs a new HS256 JWT for a WeChat mini-program user.
func GenerateMPToken(userID uint64, openID, appID, secret string, expireHours int) (string, int64, error) {
	expireAt := time.Now().Add(time.Duration(expireHours) * time.Hour)
	claims := MPClaims{
		UserID: userID,
		OpenID: openID,
		AppID:  appID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "zw-note-backend:mp",
		},
	}
	return signToken(claims, secret, expireAt.Unix())
}

// ParseAdminToken validates an admin token and returns its claims.
func ParseAdminToken(tokenStr, secret string) (*AdminClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AdminClaims{}, keyFunc(secret))
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*AdminClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid admin token")
	}
	return claims, nil
}

// ParseMPToken validates a mini-program token and returns its claims.
func ParseMPToken(tokenStr, secret string) (*MPClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &MPClaims{}, keyFunc(secret))
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*MPClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid mp token")
	}
	return claims, nil
}

func signToken(claims jwt.Claims, secret string, expireUnix int64) (string, int64, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", 0, err
	}
	return tokenStr, expireUnix, nil
}

func keyFunc(secret string) jwt.Keyfunc {
	return func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	}
}
