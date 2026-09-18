// Package pubsub provides pubsub for the application
package pubsub

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(New),
)
