package storage

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewStorageConfig),
	fx.Provide(NewMinio),
	fx.Provide(NewS3),
	fx.Provide(NewLocal),
)
