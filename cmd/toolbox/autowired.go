package main

import (
	"sapasora/platform/support/autowired"
	"sapasora/platform/support/console"

	"github.com/spf13/cobra"
)

func RunAutoWired() {
	aw := autowired.NewAutoWired(autowired.Config{
		OutputFile:  "internal/app/bootstrap/autowired.go",
		PackageName: "bootstrap",
	})

	console.Info("Scanning directories for autowired dependencies")
	if err := aw.ScanDirectories(); err != nil {
		console.Error("Error scanning directories: %s", err.Error())
	}

	console.Info("Generating autowired dependencies")
	if err := aw.GenerateCode(); err != nil {
		console.Error("Error generating code: %s", err.Error())
	}

	console.Info("Autowired dependencies generated successfully")
}

var AutowiredCmd = &cobra.Command{
	Use:   "autowired",
	Short: "Generate autowired dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		RunAutoWired()
	},
}
