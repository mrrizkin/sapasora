package main

import (
	"sapasora/platform/support/console"
	"sapasora/platform/support/kuproy"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var KuproyCMD = &cobra.Command{
	Use:   "scaffold",
	Short: "Scaffold a new module",
}

var scaffoldModelCMD = &cobra.Command{
	Use:   "model module entity table_name [fields...]",
	Short: "Generate models",
	Args:  cobra.MinimumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {

		kp := kuproy.Kuproy{
			Task: "model",
			Args: args,
		}

		if err := kp.GenerateCode(); err != nil {
			console.Error("Error generating code: %s", err.Error())
			os.Exit(1)
		}

		if err := CreateMigration(fmt.Sprintf("create_%s_table", args[2]), args[3:]...); err != nil {
			console.Error("Error generate migration: %s", err.Error())
			os.Exit(1)
		}

		RunAutoWired()

		console.Info("Success!")
	},
}

var scaffoldContextCMD = &cobra.Command{
	Use:   "context module entity table_name [fields...]",
	Short: "Generate contexts",
	Args:  cobra.MinimumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {

		kp := kuproy.Kuproy{
			Task: "context",
			Args: args,
		}

		if err := kp.GenerateCode(); err != nil {
			console.Error("Error generating code: %s", err.Error())
			os.Exit(1)
		}

		if err := CreateMigration(fmt.Sprintf("create_%s_table", args[2]), args[3:]...); err != nil {
			console.Error("Error generate migration: %s", err.Error())
			os.Exit(1)
		}

		RunAutoWired()

		console.Info("Success!")
	},
}

var scaffoldJSONCMD = &cobra.Command{
	Use:   "json module entity table_name [field:type...]",
	Short: "Generate JSON api",
	Args:  cobra.MinimumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		kp := kuproy.Kuproy{
			Task: "json",
			Args: args,
		}

		if err := kp.GenerateCode(); err != nil {
			console.Error("Error generating code: %s", err.Error())
			os.Exit(1)
		}

		if err := CreateMigration(fmt.Sprintf("create_%s_table", args[2]), args[3:]...); err != nil {
			console.Error("Error generate migration: %s", err.Error())
			os.Exit(1)
		}

		GenerateSwagger()
		RunAutoWired()

		console.Info("Success!")
	},
}

var scaffoldPolicyCMD = &cobra.Command{
	Use:   "policy StructName policy_name",
	Short: "Generate policies",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		kp := kuproy.Kuproy{
			Task: "policy",
			Args: args,
		}

		if err := kp.GenerateCode(); err != nil {
			console.Error("Error generating code: %s", err.Error())
			os.Exit(1)
		}

		RunAutoWired()

		console.Info("Success!")
	},
}

var scaffoldAuthCMD = &cobra.Command{
	Use:   "auth module entity table_name [fields...]",
	Short: "Generate auth implementation",
	Args:  cobra.MinimumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		kp := kuproy.Kuproy{
			Task: "auth",
			Args: args,
		}

		if err := kp.GenerateCode(); err != nil {
			console.Error("Error generating code: %s", err.Error())
			os.Exit(1)
		}

		if err := CreateMigration(fmt.Sprintf("create_%s_table", args[2]), args[3:]...); err != nil {
			console.Error("Error generate migration: %s", err.Error())
			os.Exit(1)
		}

		RunAutoWired()

		console.Info("Success!")
	},
}

func init() {
	KuproyCMD.AddCommand(
		scaffoldAuthCMD,
		scaffoldModelCMD,
		scaffoldContextCMD,
		scaffoldJSONCMD,
		scaffoldPolicyCMD,
	)
}
