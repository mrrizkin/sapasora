// Package view provides a facade for rendering templates and Inertia.js pages.
package view

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewView),
)
