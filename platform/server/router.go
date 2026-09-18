package server

import (
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	Path     string
	Handle   func(fiber.Router)
	Pipeline []fiber.Handler
}

func NewRouter(path string, routeSetup func(fiber.Router), pipeline ...fiber.Handler) Router {
	return Router{
		Path:     path,
		Handle:   routeSetup,
		Pipeline: pipeline,
	}
}

type Routers []Router

func (routers Routers) Register(app *fiber.App) {
	for _, r := range routers {
		r.Handle(app.Group(r.Path, r.Pipeline...))
	}
}

func Group(r fiber.Router, middlewares []fiber.Handler, setup func(fiber.Router)) {
	routes := CostumRouter{}
	setup(&routes)
	routes.registers(r, middlewares)
}

type CustomRoute struct {
	t        string
	args     []any
	path     string
	handlers []fiber.Handler
	method   string
	prefix   string
	root     string
	config   fiber.Static
	route    func(fiber.Router)
	names    []string
	group    *CostumRouter
	name     string
	fiber    *fiber.App
}

type CostumRouter struct {
	routes []CustomRoute
}

func (r *CostumRouter) registers(router fiber.Router, middlewares []fiber.Handler) {
	for _, route := range r.routes {
		switch route.t {
		case "use":
			router.Use(route.args...)
		case "get":
			router.Get(route.path, append(middlewares, route.handlers...)...)
		case "head":
			router.Head(route.path, append(middlewares, route.handlers...)...)
		case "post":
			router.Post(route.path, append(middlewares, route.handlers...)...)
		case "put":
			router.Put(route.path, append(middlewares, route.handlers...)...)
		case "delete":
			router.Delete(route.path, append(middlewares, route.handlers...)...)
		case "connect":
			router.Connect(route.path, append(middlewares, route.handlers...)...)
		case "options":
			router.Options(route.path, append(middlewares, route.handlers...)...)
		case "trace":
			router.Trace(route.path, append(middlewares, route.handlers...)...)
		case "patch":
			router.Patch(route.path, append(middlewares, route.handlers...)...)
		case "add":
			router.Add(route.method, route.path, append(middlewares, route.handlers...)...)
		case "static":
			router.Static(route.prefix, route.root, route.config)
		case "all":
			router.All(route.path, append(middlewares, route.handlers...)...)
		case "group":
			g := router.Group(route.prefix, append(middlewares, route.handlers...)...)
			if route.group != nil {
				route.group.registers(g, []fiber.Handler{})
			}
		case "route":
			router.Route(route.prefix, route.route, route.names...)
		case "mount":
			router.Mount(route.prefix, route.fiber)
		case "name":
			router.Name(route.name)
		}
	}
}

func (r *CostumRouter) Use(args ...any) fiber.Router {
	r.routes = append(r.routes, CustomRoute{
		t:    "use",
		args: args,
	})

	return r
}

func (r *CostumRouter) Get(path string, handlers ...fiber.Handler) fiber.Router {
	r.routes = append(r.routes, CustomRoute{
		t:        "get",
		path:     path,
		handlers: handlers,
	})

	return r
}

func (r *CostumRouter) Head(path string, handlers ...fiber.Handler) fiber.Router {
	r.routes = append(r.routes, CustomRoute{
		t:        "head",
		path:     path,
		handlers: handlers,
	})

	return r
}

func (r *CostumRouter) Post(path string, handlers ...fiber.Handler) fiber.Router {
	r.routes = append(r.routes, CustomRoute{
		t:        "post",
		path:     path,
		handlers: handlers,
	})

	return r
}

func (r *CostumRouter) Put(path string, handlers ...fiber.Handler) fiber.Router {
	r.routes = append(r.routes, CustomRoute{
		t:        "put",
		path:     path,
		handlers: handlers,
	})

	return r
}

func (r *CostumRouter) Delete(path string, handlers ...fiber.Handler) fiber.Router {
	r.routes = append(r.routes, CustomRoute{
		t:        "delete",
		path:     path,
		handlers: handlers,
	})

	return r
}

func (r *CostumRouter) Connect(path string, handlers ...fiber.Handler) fiber.Router {
	r.routes = append(r.routes, CustomRoute{
		t:        "connect",
		path:     path,
		handlers: handlers,
	})

	return r
}

func (r *CostumRouter) Options(path string, handlers ...fiber.Handler) fiber.Router {
	r.routes = append(r.routes, CustomRoute{
		t:        "options",
		path:     path,
		handlers: handlers,
	})

	return r
}

func (r *CostumRouter) Trace(path string, handlers ...fiber.Handler) fiber.Router {
	r.routes = append(r.routes, CustomRoute{
		t:        "trace",
		path:     path,
		handlers: handlers,
	})

	return r
}

func (r *CostumRouter) Patch(path string, handlers ...fiber.Handler) fiber.Router {
	r.routes = append(r.routes, CustomRoute{
		t:        "patch",
		path:     path,
		handlers: handlers,
	})

	return r
}

func (r *CostumRouter) Add(method, path string, handlers ...fiber.Handler) fiber.Router {
	r.routes = append(r.routes, CustomRoute{
		t:        "add",
		method:   method,
		path:     path,
		handlers: handlers,
	})

	return r
}

func (r *CostumRouter) Static(prefix, root string, config ...fiber.Static) fiber.Router {
	r.routes = append(r.routes, CustomRoute{
		t:      "static",
		prefix: prefix,
		root:   root,
		config: config[0],
	})

	return r
}

func (r *CostumRouter) All(path string, handlers ...fiber.Handler) fiber.Router {
	r.routes = append(r.routes, CustomRoute{
		t:        "all",
		path:     path,
		handlers: handlers,
	})

	return r
}

func (r *CostumRouter) Group(prefix string, handlers ...fiber.Handler) fiber.Router {
	newRouter := CostumRouter{
		routes: make([]CustomRoute, 0),
	}

	r.routes = append(r.routes, CustomRoute{
		t:        "group",
		prefix:   prefix,
		handlers: handlers,
		group:    &newRouter,
	})

	return &newRouter
}

func (r *CostumRouter) Route(prefix string, fn func(fiber.Router), name ...string) fiber.Router {
	r.routes = append(r.routes, CustomRoute{
		t:      "route",
		prefix: prefix,
		route:  fn,
		names:  name,
	})

	return r
}

func (r *CostumRouter) Mount(prefix string, fiber *fiber.App) fiber.Router {
	r.routes = append(r.routes, CustomRoute{
		t:      "mount",
		prefix: prefix,
		fiber:  fiber,
	})

	return r
}

func (r *CostumRouter) Name(name string) fiber.Router {
	r.routes = append(r.routes, CustomRoute{
		t:    "name",
		name: name,
	})

	return r
}
