package main

import (
	"context"
	"sapasora/internal/app"
	"sapasora/internal/modules/telegram"
	"sapasora/internal/modules/whatsapp"
	"sapasora/platform/config"
	"sapasora/platform/logger"
	"sapasora/platform/scheduler"
	"sapasora/platform/server"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
)

func serve() {
	var shutdown fx.Shutdowner
	app.Bootstrap(
		fx.Populate(&shutdown),
		app.Module,
		app.FxLogger,
		app.Start(func(
			srv *server.Server,
			log *logger.Logger,
			cfg config.Config,
			cron *scheduler.CRON,
			whatsmeow *whatsapp.Whatsmeow,
			tdlib *telegram.TDLib,
		) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			g, gCtx := errgroup.WithContext(ctx)

			cron.Start()
			whatsmeow.ConnectDevices(gCtx)
			tdlib.ConnectDevices(gCtx)

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

			g.Go(func() error {
				port := cfg.GetInt("app.port", 3000)
				url := cfg.GetString("app.url", "http://localhost:3000")
				addr := fmt.Sprintf(":%d", port)
				log.Info("Server starting", "port", port, "address", url)

				if err := srv.App().Listen(addr); err != nil {
					log.Error("Server failed to start", "error", err.Error())
					return err
				}
				return nil
			})

			g.Go(func() error {
				select {
				case <-quit:
					log.Info("Shutdown signal received, starting graceful shutdown...")

					shutdownCtx, shutdownCancel := context.WithTimeout(
						context.Background(),
						30*time.Second,
					)
					defer shutdownCancel()

					if err := srv.App().ShutdownWithContext(shutdownCtx); err != nil {
						log.Error("Server forced to shutdown", "error", err.Error())
						return err
					}

					log.Info("Server shutdown completed")
					if err := shutdown.Shutdown(); err != nil {
						log.Error("Server forced to shutdown", "error", err.Error())
						return err
					}
					return nil
				case <-gCtx.Done():
					return gCtx.Err()
				}
			})

			if err := g.Wait(); err != nil {
				log.Error("Application error", "error", err.Error())
			}
		}),
	)
}
