// Package inertia provides Inertia.js for the application.
package inertia

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewInertiaConfig),
	fx.Provide(NewInertia),
)
