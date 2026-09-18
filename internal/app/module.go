// Package app provides the application's root module
package app

import (
	"sapasora/internal/app/bootstrap"
	"sapasora/platform/config"
	"sapasora/platform/database"
	"sapasora/platform/database/migrations"
	"sapasora/platform/logger"
	"sapasora/platform/pubsub"
	"sapasora/platform/scheduler"
	"sapasora/platform/server"
	"sapasora/platform/session"
	"sapasora/platform/swagger"
	"sapasora/platform/ui/bundler"
	"sapasora/platform/ui/inertia"
	"sapasora/platform/ui/view"
	"sapasora/platform/validator"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

var (
	// Add manually here for module that doesn't have @wired: tag
	modules = loadModules(
		bundler.Module,
		config.Module,
		database.Module,
		migrations.Module,
		inertia.Module,
		logger.Module,
		pubsub.Module,
		scheduler.Module,
		server.Module,
		swagger.Module,
		session.Module,
		validator.Module,
		view.Module,
	)

	Module = appModule(
		modules,
		bootstrap.AutowiredModule,
	)

	FxLogger = fx.WithLogger(func(log *logger.Logger) fxevent.Logger {
		return log.GetFxLogger()
	})
)

func Bootstrap(modules ...any) {
	fx.New(loadModules(modules...)).Run()
}

func Start(fn any) fx.Option {
	return fx.Invoke(fn)
}

func appModule(modules ...any) fx.Option {
	var options []fx.Option
	for _, module := range modules {
		switch m := module.(type) {
		case []fx.Option:
			options = append(options, m...)
		case fx.Option:
			options = append(options, m)
		}
	}

	return fx.Module("app", options...)
}

func loadModules(modules ...any) fx.Option {
	var options []fx.Option
	for _, module := range modules {
		switch m := module.(type) {
		case []fx.Option:
			options = append(options, m...)
		case fx.Option:
			options = append(options, m)
		}
	}

	return fx.Options(options...)
}
