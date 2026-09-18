// Package logger provides a logger for the application
package logger

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewLogger),
)
