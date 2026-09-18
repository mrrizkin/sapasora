package main

import (
	"sapasora/internal/app"
	"sapasora/platform/database/migrations"
	"sapasora/platform/support/arr"
	"sapasora/platform/support/console"
	"sapasora/platform/support/kuproy"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var MigrateCMD = &cobra.Command{
	Use:   "migrate",
	Short: "Manage database migrations",
}

var migrateRunCMD = &cobra.Command{
	Use:   "run",
	Short: "Run pending migrations",
	Run: func(cmd *cobra.Command, args []string) {
		bootstrapAppAndRun(func(migration *migrations.MigrationRunner) error {
			return migration.Run()
		})
	},
}

var migrateRollbackCMD = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback last batch of migrations",
	Run: func(cmd *cobra.Command, args []string) {
		bootstrapAppAndRun(func(migration *migrations.MigrationRunner) error {
			return migration.Rollback()
		})
	},
}

var migrateResetCMD = &cobra.Command{
	Use:   "reset",
	Short: "Rollback all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		bootstrapAppAndRun(func(migration *migrations.MigrationRunner) error {
			return migration.Reset()
		})
	},
}

var migrateFreshCMD = &cobra.Command{
	Use:   "fresh",
	Short: "Drop all tables and re-run all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		bootstrapAppAndRun(func(migration *migrations.MigrationRunner) error {
			return migration.Fresh()
		})
	},
}

var migrateStatusCMD = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	Run: func(cmd *cobra.Command, args []string) {
		bootstrapAppAndRun(func(migration *migrations.MigrationRunner) error {
			return migration.Status()
		})
	},
}

func CreateMigration(migrationName string, args ...string) error {
	kp := kuproy.Kuproy{
		Task: "migration",
		Args: arr.Merge([]string{migrationName}, args),
	}

	if err := kp.GenerateCode(); err != nil {
		return err
	}

	RunAutoWired()

	return nil
}

var migrateMakeCMD = &cobra.Command{
	Use:   "make [name] [fields...]",
	Short: "Generate a new migration file",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if err := CreateMigration(args[0], args[1:]...); err != nil {
			console.Error("Error create migration: %s", err.Error())
			os.Exit(1)
		}

		RunAutoWired()

		console.Info("Success!")
	},
}

// bootstrapAppAndRun initializes the app and runs the migration function
func bootstrapAppAndRun(migrationFunc func(runner *migrations.MigrationRunner) error) {
	var shutdown fx.Shutdowner
	app.Bootstrap(
		fx.Populate(&shutdown),
		app.Module,
		app.FxLogger,
		app.Start(
			func(migration *migrations.MigrationRunner) {
				if err := migrationFunc(migration); err != nil {
					console.Error("Migration failed: %s", err)
					os.Exit(1)
				}

				if err := shutdown.Shutdown(); err != nil {
					console.Error("Error shutting down application: %s", err)
					os.Exit(1)
				}
			},
		),
	)
}

func init() {
	MigrateCMD.AddCommand(
		migrateRunCMD,
		migrateRollbackCMD,
		migrateResetCMD,
		migrateFreshCMD,
		migrateStatusCMD,
		migrateMakeCMD,
	)
}
