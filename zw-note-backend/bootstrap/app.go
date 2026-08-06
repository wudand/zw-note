package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"zw-note-backend/internal/config"
	"zw-note-backend/pkg/database"
	"zw-note-backend/pkg/logger"
	"go.uber.org/zap"
)

// App holds all long-lived application resources and owns the server lifecycle.
type App struct {
	cfg    *config.Config
	db     *sqlx.DB
	log    *zap.Logger
	server *http.Server
}

// NewApp initialises all infrastructure components and wires them together.
// It returns a ready-to-run App or an error if any component fails to start.
func NewApp(cfg *config.Config) (*App, error) {
	log, err := logger.Init(cfg.Log)
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}

	db, err := database.NewPostgres(database.PostgresConfig{
		DSN:             cfg.Database.DSN,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	})
	if err != nil {
		return nil, fmt.Errorf("init database: %w", err)
	}
	log.Info("database connected", zap.String("driver", cfg.Database.Driver))

	router := SetupRouter(cfg, db, log)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	return &App{cfg: cfg, db: db, log: log, server: srv}, nil
}

// Run starts the HTTP server in a goroutine and blocks until a termination
// signal is received, then performs a graceful shutdown.
func (a *App) Run() error {
	go func() {
		a.log.Info("server starting", zap.String("addr", a.server.Addr))
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.log.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	a.log.Info("shutdown signal received", zap.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	if err := a.db.Close(); err != nil {
		a.log.Error("close database", zap.Error(err))
	}

	logger.Sync()
	a.log.Info("server exited cleanly")
	return nil
}
