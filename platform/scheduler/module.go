// Package scheduler provides scheduler for the application
package scheduler

import (
	"go.uber.org/fx"
)

// Module provides the scheduler as an FX module
var Module = fx.Options(
	fx.Provide(NewScheduler),
	fx.Provide(NewAdvancedSchedulerFromLifecycle),
)
