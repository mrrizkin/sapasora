// Package validator provides validator for the application
package validator

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewValidator),
)
