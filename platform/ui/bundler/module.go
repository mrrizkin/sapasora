// Package bundler provides bundler for the application.
package bundler

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewBundler),
)
