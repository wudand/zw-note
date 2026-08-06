package main

import (
	"fmt"
	"os"

	"zw-note-backend/bootstrap"
	"zw-note-backend/internal/config"

	"github.com/spf13/cobra"
)

// Version is injected at build time via -ldflags.
var Version = "dev"

// @title           zw-note-backend
// @version         1.0
// @description     RESTful API for PC admin and WeChat mini-program
// @host            localhost:8004
// @BasePath        /
// @schemes         http
func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	var cfgFile string

	root := &cobra.Command{
		Use:   "zw-note-backend",
		Short: "A production-ready RESTful API server",
		Long: `zw-note-backend is a RESTful API server built with Go and Gin.
It supports JWT authentication, Prometheus metrics, structured logging,
and a clean layered architecture.`,
		SilenceUsage: true,
	}

	root.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file path (default: ./configs/config.yaml)")

	root.AddCommand(newStartCmd(&cfgFile))
	root.AddCommand(newVersionCmd())

	return root
}

func newStartCmd(cfgFile *string) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the HTTP API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(*cfgFile)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			app, err := bootstrap.NewApp(cfg)
			if err != nil {
				return fmt.Errorf("init app: %w", err)
			}

			return app.Run()
		},
	}
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the application version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("zw-note-backend %s\n", Version)
		},
	}
}
