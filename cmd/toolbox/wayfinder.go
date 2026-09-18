package main

import (
	"sapasora/internal/app"
	"sapasora/platform/logger"
	"sapasora/platform/server"
	"sapasora/platform/support/console"
	"sapasora/platform/support/wayfinder"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// Command flags
var (
	routesOutputPath     string
	controllerOutputPath string
	wayfinderOutputPath  string
	cleanOutputPath      bool

	WayfinderCMD = &cobra.Command{
		Use:   "wayfinder",
		Short: "Generate TypeScript routes from Fiber routes",
		Run: func(cmd *cobra.Command, args []string) {
			var shutdown fx.Shutdowner
			app.Bootstrap(
				fx.Populate(&shutdown),
				app.Module,
				app.FxLogger,
				app.Start(func(srv *server.Server, log *logger.Logger) {
					wf := wayfinder.NewWayFinder(wayfinder.Config{
						Routes:               srv.App().GetRoutes(true),
						Log:                  log,
						CleanOutputPath:      cleanOutputPath,
						RouteOutputPath:      routesOutputPath,
						ControllerOutputPath: controllerOutputPath,
						WayfinderOutputPath:  wayfinderOutputPath,
					})
					wf.Generate()
					console.Info("Routes generated successfully")

					if err := shutdown.Shutdown(); err != nil {
						console.Error("Error shutting down application: %s", err.Error())
						os.Exit(1)
					}
				}),
			)
		},
	}
)

func init() {
	WayfinderCMD.Flags().StringVarP(
		&routesOutputPath,
		"routes-output",
		"r",
		"resources/js/lib/routes",
		"Output path for routes",
	)
	WayfinderCMD.Flags().StringVarP(
		&controllerOutputPath,
		"actions-output",
		"a",
		"resources/js/lib/actions",
		"Output path for controllers",
	)
	WayfinderCMD.Flags().StringVarP(
		&wayfinderOutputPath,
		"wayfinder-output",
		"w",
		"resources/js/lib",
		"Output path for core wayfinder helpers",
	)
	WayfinderCMD.Flags().BoolVarP(
		&cleanOutputPath,
		"clean",
		"c",
		true,
		"Clean output path before generating routes",
	)
}
