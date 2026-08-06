package admin

import (
	"zw-note-backend/internal/config"
	"zw-note-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes mounts all /api/admin/v1 routes.
//
// Route permission matrix:
//
//	POST   /auth/login              – public
//	POST   /users                   – super_admin
//	GET    /users                   – admin, super_admin
//	GET    /users/:id               – admin, super_admin
//	PUT    /users/:id               – super_admin
//	DELETE /users/:id               – super_admin
//
//	GET    /categories              – admin, super_admin, viewer
//	GET    /categories/:id          – admin, super_admin, viewer
//	POST   /categories              – admin, super_admin
//	PUT    /categories/:id          – admin, super_admin
//	DELETE /categories/:id          – admin, super_admin
//
//	POST   /upload/image             – admin, super_admin
//
//	GET    /products                 – admin, super_admin, viewer
//	GET    /products/:id             – admin, super_admin, viewer
//	POST   /products                 – admin, super_admin
//	PUT    /products/:id             – admin, super_admin
//	DELETE /products/:id             – admin, super_admin
//
//	GET    /mp-users/:mp_user_id/addresses                – admin, super_admin, viewer
//	GET    /mp-users/:mp_user_id/addresses/:id            – admin, super_admin, viewer
//	POST   /mp-users/:mp_user_id/addresses                – admin, super_admin
//	PUT    /mp-users/:mp_user_id/addresses/:id            – admin, super_admin
//	PUT    /mp-users/:mp_user_id/addresses/:id/default    – admin, super_admin
//	DELETE /mp-users/:mp_user_id/addresses/:id            – admin, super_admin
//
//	GET    /coupons                  – admin, super_admin, viewer
//	GET    /coupons/:id              – admin, super_admin, viewer
//	POST   /coupons                  – admin, super_admin
//	PUT    /coupons/:id              – admin, super_admin
//	PUT    /coupons/:id/cancel       – admin, super_admin
//	DELETE /coupons/:id              – admin, super_admin
//
//	GET    /carousels               – admin, super_admin, viewer
//	GET    /carousels/:id           – admin, super_admin, viewer
//	POST   /carousels               – admin, super_admin
//	PUT    /carousels/:id           – admin, super_admin
//	DELETE /carousels/:id           – admin, super_admin
//
//	GET    /redemption-codes        – admin, super_admin, viewer
//	GET    /redemption-codes/:id    – admin, super_admin, viewer
//	POST   /redemption-codes        – admin, super_admin
//	PUT    /redemption-codes/:id    – admin, super_admin
//	DELETE /redemption-codes/:id    – admin, super_admin
func RegisterRoutes(r *gin.Engine, userHandler *UserHandler, categoryHandler *CategoryHandler, productHandler *ProductHandler, uploadHandler *UploadHandler, addressHandler *AddressHandler, couponHandler *CouponHandler, carouselHandler *CarouselHandler, redemptionCodeHandler *RedemptionCodeHandler, jwtCfg config.JWTConfig) {
	adminV1 := r.Group("/api/admin/v1")

	// Public – no auth required
	auth := adminV1.Group("/auth")
	{
		auth.POST("/login", userHandler.Login)
	}

	// All routes below require a valid admin JWT
	protected := adminV1.Group("")
	protected.Use(middleware.AdminJWTAuth(jwtCfg))
	{
		users := protected.Group("/users")

		// Readers: admin + super_admin
		readers := middleware.RequireRole("admin", "super_admin")
		users.GET("", readers, userHandler.ListUsers)
		users.GET("/:id", readers, userHandler.GetUser)

		// Writers: super_admin only
		writers := middleware.RequireRole("super_admin")
		users.POST("", writers, userHandler.CreateUser)
		users.PUT("/:id", writers, userHandler.UpdateUser)
		users.DELETE("/:id", writers, userHandler.DeleteUser)

		// Categories
		categories := protected.Group("/categories")
		readersAll := middleware.RequireRole("admin", "super_admin", "viewer")
		categories.GET("", readersAll, categoryHandler.List)
		categories.GET("/:id", readersAll, categoryHandler.GetByID)
		writersAdmin := middleware.RequireRole("admin", "super_admin")
		categories.POST("", writersAdmin, categoryHandler.Create)
		categories.PUT("/:id", writersAdmin, categoryHandler.Update)
		categories.DELETE("/:id", writersAdmin, categoryHandler.Delete)

		// Upload (admin, super_admin)
		protected.POST("/upload/image", writersAdmin, uploadHandler.UploadImage)

		// Products
		products := protected.Group("/products")
		products.GET("", readersAll, productHandler.List)
		products.GET("/:id", readersAll, productHandler.GetByID)
		products.POST("", writersAdmin, productHandler.Create)
		products.PUT("/:id", writersAdmin, productHandler.Update)
		products.DELETE("/:id", writersAdmin, productHandler.Delete)

		// MP user addresses
		mpUsers := protected.Group("/mp-users")
		mpUsersWithID := mpUsers.Group("/:mp_user_id")
		addresses := mpUsersWithID.Group("/addresses")
		addresses.GET("", readersAll, addressHandler.List)
		addresses.GET("/:id", readersAll, addressHandler.GetByID)
		addresses.POST("", writersAdmin, addressHandler.Create)
		addresses.PUT("/:id/default", writersAdmin, addressHandler.SetDefault) // before /:id
		addresses.PUT("/:id", writersAdmin, addressHandler.Update)
		addresses.DELETE("/:id", writersAdmin, addressHandler.Delete)

		// Coupons
		coupons := protected.Group("/coupons")
		coupons.GET("", readersAll, couponHandler.List)
		coupons.GET("/:id", readersAll, couponHandler.GetByID)
		coupons.POST("", writersAdmin, couponHandler.Create)
		coupons.PUT("/:id/cancel", writersAdmin, couponHandler.Cancel) // before /:id
		coupons.PUT("/:id", writersAdmin, couponHandler.Update)
		coupons.DELETE("/:id", writersAdmin, couponHandler.Delete)

		// Carousels
		carousels := protected.Group("/carousels")
		carousels.GET("", readersAll, carouselHandler.List)
		carousels.GET("/:id", readersAll, carouselHandler.GetByID)
		carousels.POST("", writersAdmin, carouselHandler.Create)
		carousels.PUT("/:id", writersAdmin, carouselHandler.Update)
		carousels.DELETE("/:id", writersAdmin, carouselHandler.Delete)

		// Redemption codes
		redemptionCodes := protected.Group("/redemption-codes")
		redemptionCodes.GET("", readersAll, redemptionCodeHandler.List)
		redemptionCodes.GET("/:id", readersAll, redemptionCodeHandler.GetByID)
		redemptionCodes.POST("", writersAdmin, redemptionCodeHandler.CreateBatch)
		redemptionCodes.PUT("/:id", writersAdmin, redemptionCodeHandler.Update)
		redemptionCodes.DELETE("/:id", writersAdmin, redemptionCodeHandler.Delete)
	}
}
