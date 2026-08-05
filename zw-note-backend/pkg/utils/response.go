package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is the unified envelope for every API response.
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Business error codes.
// 0        – success
// 1xxxx    – generic HTTP-level errors
// 2xxxx    – admin user domain errors
// 3xxxx    – mini-program domain errors
const (
	CodeSuccess      = 0
	CodeParamError   = 10001
	CodeUnauthorized = 10002
	CodeForbidden    = 10003
	CodeNotFound     = 10004
	CodeServerError  = 10005

	// Admin user
	CodeUserNotFound = 20001
	CodeUserExists   = 20002
	CodeEmailExists  = 20003
	CodeInvalidPass  = 20004
	CodeUserDisabled = 20005

	// Mini-program user
	CodeMPUserNotFound = 30001
	CodeWxLoginFailed  = 30002
	CodeInvalidAppKey  = 30003

	// Product category (shared)
	CodeCategoryNotFound = 40001

	// Product
	CodeProductNotFound = 40002

	// Address
	CodeAddressNotFound     = 40003
	CodeAddressLimitReached = 40004

	// Coupon
	CodeCouponNotFound     = 40005
	CodeCouponCancelled    = 40006
	CodeCouponClaimLimit   = 40007
	CodeUserCouponNotFound = 40008

	// Carousel
	CodeCarouselNotFound     = 40009
	CodeCarouselEnabledLimit = 40010

	// Redemption code
	CodeRedemptionCodeNotFound           = 40011
	CodeRedemptionCodeAlreadyUsed        = 40012
	CodeRedemptionCodeCooldown           = 40013
	CodeRedemptionCodeNotYours           = 40014
	CodeRedemptionCodeProductAlreadySelected = 40015
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: CodeSuccess, Message: "success", Data: data})
}

func SuccessCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{Code: CodeSuccess, Message: "success", Data: data})
}

func Fail(c *gin.Context, httpStatus, code int, message string) {
	c.JSON(httpStatus, Response{Code: code, Message: message, Data: nil})
}

func ParamError(c *gin.Context, message string) {
	if message == "" {
		message = "parameter error"
	}
	Fail(c, http.StatusBadRequest, CodeParamError, message)
}

func Unauthorized(c *gin.Context) {
	Fail(c, http.StatusUnauthorized, CodeUnauthorized, "unauthorized")
}

func Forbidden(c *gin.Context) {
	Fail(c, http.StatusForbidden, CodeForbidden, "forbidden")
}

func ServerError(c *gin.Context) {
	Fail(c, http.StatusInternalServerError, CodeServerError, "internal server error")
}

// HandleAppError maps an AppError to the appropriate HTTP response.
// Returns true if err was an AppError (response already written), false otherwise.
func HandleAppError(c *gin.Context, err error) bool {
	appErr, ok := AsAppError(err)
	if !ok {
		return false
	}

	httpStatus := http.StatusBadRequest
	switch appErr.Code {
	case CodeUnauthorized:
		httpStatus = http.StatusUnauthorized
	case CodeForbidden:
		httpStatus = http.StatusForbidden
	case CodeUserNotFound, CodeMPUserNotFound, CodeCategoryNotFound, CodeProductNotFound, CodeAddressNotFound, CodeCouponNotFound, CodeUserCouponNotFound, CodeCarouselNotFound, CodeRedemptionCodeNotFound, CodeRedemptionCodeNotYours, CodeNotFound:
		httpStatus = http.StatusNotFound
	}

	Fail(c, httpStatus, appErr.Code, appErr.Message)
	return true
}
