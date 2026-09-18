// Package session provides session for the application
package session

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewSession),
)
