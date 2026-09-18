package wayfinder

import (
	"github.com/gofiber/fiber/v2"
)

type Wayfinder struct {
	log                  Logger
	templates            *Templates
	routes               map[string][]Route
	cleanOutputPath      bool
	routeOutputPath      string
	controllerOutputPath string
	wayfinderOutputPath  string
}

type Logger interface {
	Fatal(message string, fields ...any)
	Error(message string, fields ...any)
	Warn(message string, fields ...any)
	Info(message string, fields ...any)
	Debug(message string, fields ...any)
}

type Route struct {
	Name    string
	Route   fiber.Route
	Handler Handler
}

type Handler struct {
	Name          string
	Location      string
	Path          string
	StructPointer string
}

type RouteData struct {
	Location string
	Path     string
	Methods  []string
	Name     string
	Params   []string
}

type RoutesTemplateData struct {
	Routes          []RouteData
	Name            string
	IsDefaultExport bool
	IsForm          bool
}

type IndexTemplateData struct {
	Actions []string
	Name    string
}

type RouterGroup map[string][]Route
type StructPointerGroup map[string]RouterGroup
type PathGroup map[string]StructPointerGroup
