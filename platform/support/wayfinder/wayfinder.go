// Package wayfinder provides generator to mimic the laravel/wayfinder package
package wayfinder

import (
	"github.com/gofiber/fiber/v2"
)

func DefaultConfig() Config {
	return Config{
		Routes:               []fiber.Route{},
		Log:                  &DefaultLogger{},
		CleanOutputPath:      false,
		RouteOutputPath:      "resources/js/lib/routes",
		ControllerOutputPath: "resources/js/lib/controllers",
		WayfinderOutputPath:  "resources/js/lib",
	}
}

func NewWayFinder(
	config ...Config,
) *Wayfinder {
	c := DefaultConfig()
	if len(config) > 0 {
		conf := config[0]
		if conf.Log != nil {
			c.Log = conf.Log
		}
		if conf.Routes != nil {
			c.Routes = conf.Routes
		}
		if conf.RouteOutputPath != "" {
			c.RouteOutputPath = conf.RouteOutputPath
		}
		if conf.ControllerOutputPath != "" {
			c.ControllerOutputPath = conf.ControllerOutputPath
		}
		if conf.WayfinderOutputPath != "" {
			c.WayfinderOutputPath = conf.WayfinderOutputPath
		}
		if conf.CleanOutputPath {
			c.CleanOutputPath = conf.CleanOutputPath
		}
	}

	if len(c.Routes) == 0 {
		c.Log.Fatal("No routes provided")
	}

	return &Wayfinder{
		log:                  c.Log,
		templates:            NewTemplates(),
		routes:               groupRoutes(c.Routes),
		cleanOutputPath:      c.CleanOutputPath,
		routeOutputPath:      c.RouteOutputPath,
		controllerOutputPath: c.ControllerOutputPath,
		wayfinderOutputPath:  c.WayfinderOutputPath,
	}
}
