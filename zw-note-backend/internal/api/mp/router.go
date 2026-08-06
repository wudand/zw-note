package mp

import (
	"zw-note-backend/internal/config"
	"zw-note-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes mounts all /api/mp/v1 routes.
//
//	POST   /auth/wx-login          – public (WeChat code exchange)
//	GET    /categories            – public (enabled categories only, no auth)
//
//	GET    /products              – public (enabled only, paginated, category_id filter)
//	GET    /products/:id          – public (enabled only)
//
//	GET    /coupons               – public (claimable list)
//
//	GET    /carousels             – public (home carousel)
//
//	GET    /user/profile          – JWT required
//
//	GET    /addresses             – JWT required
//	GET    /addresses/:id         – JWT required
//	POST   /addresses             – JWT required
//	PUT    /addresses/:id         – JWT required
//	PUT    /addresses/:id/default – JWT required
//	DELETE /addresses/:id         – JWT required
//
// GET    /coupons/my            – JWT required
// POST   /coupons/:id/claim     – JWT required
//
// GET    /redemption-codes/validate – JWT required
// POST   /redemption-codes/:id/select-product – JWT required
func RegisterRoutes(r *gin.Engine, authHandler *AuthHandler, categoryHandler *CategoryHandler, productHandler *ProductHandler, addressHandler *AddressHandler, couponHandler *CouponHandler, carouselHandler *CarouselHandler, redemptionCodeHandler *RedemptionCodeHandler, jwtCfg config.JWTConfig) {
	mpV1 := r.Group("/api/mp/v1")

	// Public – no auth required
	auth := mpV1.Group("/auth")
	{
		auth.POST("/wx-login", authHandler.WxLogin)
	}
	mpV1.GET("/categories", categoryHandler.List)
	mpV1.GET("/products", productHandler.List)
	mpV1.GET("/products/:id", productHandler.GetByID)
	mpV1.GET("/coupons", couponHandler.ListClaimable) // public: claimable list
	mpV1.GET("/carousels", carouselHandler.List)      // public: home carousel

	// Protected – mini-program JWT required
	protected := mpV1.Group("")
	protected.Use(middleware.MPJWTAuth(jwtCfg))
	{
		user := protected.Group("/user")
		user.GET("/profile", authHandler.GetProfile)

		addresses := protected.Group("/addresses")
		addresses.GET("", addressHandler.List)
		addresses.GET("/:id", addressHandler.GetByID)
		addresses.POST("", addressHandler.Create)
		addresses.PUT("/:id/default", addressHandler.SetDefault) // must be before /:id
		addresses.PUT("/:id", addressHandler.Update)
		addresses.DELETE("/:id", addressHandler.Delete)

		coupons := protected.Group("/coupons")
		coupons.GET("/my", couponHandler.ListMy)
		coupons.POST("/:id/claim", couponHandler.Claim)

		redemptionCodes := protected.Group("/redemption-codes")
		redemptionCodes.POST("/validate", redemptionCodeHandler.Validate)
		redemptionCodes.POST("/:id/select-product", redemptionCodeHandler.SelectProduct)
	}
}
