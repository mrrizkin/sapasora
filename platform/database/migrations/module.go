package migrations

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewMigrationRunner),
)
