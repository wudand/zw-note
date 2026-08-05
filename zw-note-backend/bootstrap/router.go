package bootstrap

import (
	"net/http"
	"path/filepath"
	"time"

	adminAPI "go-web-api/internal/api/admin"
	mpAPI "go-web-api/internal/api/mp"
	"go-web-api/internal/config"
	"go-web-api/internal/middleware"
	"go-web-api/internal/repository"
	"go-web-api/internal/service"
	"go-web-api/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	_ "go-web-api/docs"
)

// SetupRouter builds the gin engine with all global middleware, observability
// endpoints, and both the admin and mini-program route groups wired up.
func SetupRouter(cfg *config.Config, db *sqlx.DB, log *zap.Logger) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()

	// Global middleware – order: recovery → logging → cors → metrics
	r.Use(middleware.Recovery(log))
	r.Use(middleware.RequestLogger(log))
	r.Use(middleware.CORS(cfg.CORS))
	r.Use(middleware.Metrics())

	// Observability (no auth, excluded from business metrics by full-path)
	r.GET("/health", healthCheck(db))
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Static files: /uploads/* -> upload dir (for image access)
	uploadDir := cfg.Upload.Dir
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	if abs, err := filepath.Abs(uploadDir); err == nil {
		uploadDir = abs
	}
	r.Static("/uploads", uploadDir)

	// ── Admin API ─────────────────────────────────────────────────────────────
	adminUserRepo := repository.NewAdminUserRepository(db)
	adminUserSvc := service.NewAdminUserService(adminUserRepo, log)
	adminUserHandler := adminAPI.NewUserHandler(adminUserSvc, log, cfg.JWT)

	categoryRepo := repository.NewCategoryRepository(db)
	categorySvc := service.NewCategoryService(categoryRepo, log)
	categoryHandler := adminAPI.NewCategoryHandler(categorySvc, log)
	uploadHandler := adminAPI.NewUploadHandler(cfg.Upload, log)

	productRepo := repository.NewProductRepository(db)
	productSvc := service.NewProductService(productRepo, categoryRepo, log)
	productHandler := adminAPI.NewProductHandler(productSvc, log)

	txMgr := database.NewTransactionManager(db)

	addressRepo := repository.NewAddressRepository(db)
	addressSvc := service.NewAddressService(txMgr, addressRepo, log)
	addressHandler := adminAPI.NewAddressHandler(addressSvc, log)

	couponRepo := repository.NewCouponRepository(db)
	userCouponRepo := repository.NewUserCouponRepository(db)
	couponSvc := service.NewCouponService(couponRepo, userCouponRepo, log)
	couponHandler := adminAPI.NewCouponHandler(couponSvc, log)

	carouselRepo := repository.NewCarouselRepository(db)
	carouselSvc := service.NewCarouselService(carouselRepo, log)
	carouselHandler := adminAPI.NewCarouselHandler(carouselSvc, log)

	redemptionCodeRepo := repository.NewRedemptionCodeRepository(db)
	redemptionCodeProductRepo := repository.NewRedemptionCodeProductRepository(db)
	cooldownDays := cfg.Redemption.CooldownDays
	if cooldownDays <= 0 {
		cooldownDays = 30
	}
	redemptionCodeSvc := service.NewRedemptionCodeService(txMgr, redemptionCodeRepo, redemptionCodeProductRepo, productRepo, cooldownDays, log)
	redemptionCodeHandler := adminAPI.NewRedemptionCodeHandler(redemptionCodeSvc, log)

	adminAPI.RegisterRoutes(r, adminUserHandler, categoryHandler, productHandler, uploadHandler, addressHandler, couponHandler, carouselHandler, redemptionCodeHandler, cfg.JWT)

	// ── Mini-Program API ──────────────────────────────────────────────────────
	mpUserRepo := repository.NewMPUserRepository(db)
	wxSvc := service.NewWxService(cfg.Wechat)
	mpUserSvc := service.NewMPUserService(mpUserRepo, wxSvc, cfg.Wechat, log)
	mpAuthHandler := mpAPI.NewAuthHandler(mpUserSvc, log, cfg.JWT)

	mpCategoryHandler := mpAPI.NewCategoryHandler(categorySvc, log)
	mpProductHandler := mpAPI.NewProductHandler(productSvc, log)
	mpAddressHandler := mpAPI.NewAddressHandler(addressSvc, log)
	mpCouponHandler := mpAPI.NewCouponHandler(couponSvc, log)
	mpCarouselHandler := mpAPI.NewCarouselHandler(carouselSvc, log)
	mpRedemptionCodeHandler := mpAPI.NewRedemptionCodeHandler(redemptionCodeSvc, log)

	mpAPI.RegisterRoutes(r, mpAuthHandler, mpCategoryHandler, mpProductHandler, mpAddressHandler, mpCouponHandler, mpCarouselHandler, mpRedemptionCodeHandler, cfg.JWT)

	return r
}

func healthCheck(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := "ok"
		httpStatus := http.StatusOK

		if err := db.PingContext(c.Request.Context()); err != nil {
			status = "degraded"
			httpStatus = http.StatusServiceUnavailable
		}

		c.JSON(httpStatus, gin.H{
			"status": status,
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	}
}
