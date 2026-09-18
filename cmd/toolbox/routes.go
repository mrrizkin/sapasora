package main

import (
	"sapasora/internal/app"
	"sapasora/platform/server"
	"sapasora/platform/support/arr"
	"sapasora/platform/support/console"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var RoutesCMD = &cobra.Command{
	Use:   "routes",
	Short: "Show app routes",
	Run: func(cmd *cobra.Command, args []string) {
		var shutdown fx.Shutdowner
		app.Bootstrap(
			fx.Populate(&shutdown),
			app.Module,
			app.FxLogger,
			app.Start(func(srv *server.Server) {
				httpMethodHierarcy := map[string]int{
					"GET":     0,
					"HEAD":    1,
					"POST":    2,
					"PUT":     3,
					"DELETE":  4,
					"CONNECT": 5,
					"OPTIONS": 6,
					"TRACE":   7,
					"PATCH":   8,
				}

				groupedRoutes := arr.GroupBy(
					srv.App().GetRoutes(true),
					func(r fiber.Route) string {
						parts := strings.Split(r.Name, ".")
						if len(parts) > 0 {
							parts = parts[:len(parts)-1]
							return strings.Join(parts, ".")
						}

						return r.Name
					})

				groups := arr.Keys(groupedRoutes)
				slices.SortFunc(groups, func(a, b string) int {
					if len(b) == 0 {
						return -1
					}

					if len(a) == 0 {
						return 1
					}

					if a > b {
						return 1
					}

					return -1
				})

				tables := make([][]string, 0)
				for _, group := range groups {
					routes := groupedRoutes[group]
					slices.SortFunc(routes, func(a, b fiber.Route) int {
						if len(a.Path) == len(b.Path) {
							aMethod, ok := httpMethodHierarcy[a.Method]
							if !ok {
								aMethod = -1
							}

							bMethod, ok := httpMethodHierarcy[b.Method]
							if !ok {
								bMethod = -1
							}

							if aMethod > bMethod {
								return 1
							}

							return -1
						}

						if a.Path > b.Path {
							return 1
						}

						return -1
					})

					for _, route := range routes {
						tables = append(
							tables,
							[]string{group, route.Path, route.Method, route.Name},
						)
					}
				}

				console.Table([]string{"Group", "Path", "Method", "Name"}, tables)

				if err := shutdown.Shutdown(); err != nil {
					console.Error("Error shutting down application: %s", err.Error())
				}
			}),
		)
	},
}
