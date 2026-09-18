// Package toolbox provides the tool root command for the application
package main

import (
	"sapasora/cmd/toolbox/apikeys"
	"sapasora/platform/support/console"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "toolbox",
		Short: "Application Toolbox",
	}

	cmd.AddCommand(
		AutowiredCmd,
		KuproyCMD,
		MigrateCMD,
		RoutesCMD,
		SwaggerCMD,
		WayfinderCMD,
		AdminUserCmd,
		apikeys.APIKeysCmd,
	)

	if err := cmd.Execute(); err != nil {
		console.Error("Error executing toolbox: %s", err.Error())
		os.Exit(1)
	}
}
