package utils

import "errors"

// AppError is a structured application-level error that carries an HTTP-friendly
// business error code alongside a human-readable message.
type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string { return e.Message }

// Admin user error codes (2xxxx range)
var (
	ErrUserNotFound    = &AppError{Code: CodeUserNotFound, Message: "user not found"}
	ErrUserExists      = &AppError{Code: CodeUserExists, Message: "username already exists"}
	ErrEmailExists     = &AppError{Code: CodeEmailExists, Message: "email already exists"}
	ErrInvalidPassword = &AppError{Code: CodeInvalidPass, Message: "invalid password"}
	ErrUserDisabled    = &AppError{Code: CodeUserDisabled, Message: "user is disabled"}
)

// Mini-program user error codes (3xxxx range)
var (
	ErrMPUserNotFound = &AppError{Code: CodeMPUserNotFound, Message: "mp user not found"}
	ErrWxLoginFailed  = &AppError{Code: CodeWxLoginFailed, Message: "wechat login failed"}
	ErrInvalidAppKey  = &AppError{Code: CodeInvalidAppKey, Message: "invalid wechat app key"}
)

// Product category error codes (4xxxx range)
var (
	ErrCategoryNotFound = &AppError{Code: CodeCategoryNotFound, Message: "category not found"}
	ErrProductNotFound  = &AppError{Code: CodeProductNotFound, Message: "product not found"}
)

// Address error codes (4xxxx range)
var (
	ErrAddressNotFound     = &AppError{Code: CodeAddressNotFound, Message: "address not found"}
	ErrAddressLimitReached = &AppError{Code: CodeAddressLimitReached, Message: "address limit reached (max 10)"}
)

// Coupon error codes (4xxxx range)
var (
	ErrCouponNotFound     = &AppError{Code: CodeCouponNotFound, Message: "coupon not found"}
	ErrCouponCancelled    = &AppError{Code: CodeCouponCancelled, Message: "coupon is cancelled"}
	ErrCouponClaimLimit   = &AppError{Code: CodeCouponClaimLimit, Message: "claim limit reached"}
	ErrUserCouponNotFound = &AppError{Code: CodeUserCouponNotFound, Message: "user coupon not found"}
)

// Carousel error codes (4xxxx range)
var (
	ErrCarouselNotFound     = &AppError{Code: CodeCarouselNotFound, Message: "carousel not found"}
	ErrCarouselEnabledLimit = &AppError{Code: CodeCarouselEnabledLimit, Message: "max 3 enabled carousels"}
)

// Redemption code error codes (4xxxx range)
var (
	ErrRedemptionCodeNotFound           = &AppError{Code: CodeRedemptionCodeNotFound, Message: "redemption code not found"}
	ErrRedemptionCodeAlreadyUsed        = &AppError{Code: CodeRedemptionCodeAlreadyUsed, Message: "redemption code already used"}
	ErrRedemptionCodeCooldown           = &AppError{Code: CodeRedemptionCodeCooldown, Message: "you have used a redemption code recently, please try again later"}
	ErrRedemptionCodeNotYours           = &AppError{Code: CodeRedemptionCodeNotYours, Message: "redemption code does not belong to you"}
	ErrRedemptionCodeProductAlreadySelected = &AppError{Code: CodeRedemptionCodeProductAlreadySelected, Message: "product already selected for this redemption code"}
)

// Document error codes (5xxxx range)
var (
	ErrDocumentNotFound     = &AppError{Code: CodeDocumentNotFound, Message: "document not found"}
	ErrOutlineNotFound      = &AppError{Code: CodeOutlineNotFound, Message: "outline not found"}
	ErrOutlineInvalidParent = &AppError{Code: CodeOutlineInvalidParent, Message: "invalid outline parent"}
)

// NewAppError creates an AppError with the given code and message.
func NewAppError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// AsAppError unwraps err looking for an *AppError anywhere in the chain.
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
