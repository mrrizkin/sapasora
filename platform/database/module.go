// Package database provides database for the application
package database

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewConfig),
	fx.Provide(NewDatabase),
)
