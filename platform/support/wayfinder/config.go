package wayfinder

import "github.com/gofiber/fiber/v2"

type Config struct {
	Routes               []fiber.Route
	Log                  Logger
	CleanOutputPath      bool
	RouteOutputPath      string
	ControllerOutputPath string
	WayfinderOutputPath  string
}
