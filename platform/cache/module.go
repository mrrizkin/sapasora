// Package cache provides cache for the application
package cache

import (
	"sapasora/platform/cache/config"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(config.NewConfig),
	fx.Provide(NewCache),
)
