package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"sapasora/platform/support/console"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func GenerateSwagger() {
	outputDir := "public/docs/v3"
	_, err := exec.Command("go", "tool", "swag", "--version").Output()
	if err != nil {
		console.Error("swag Go tool is unavailable: %s", err.Error())
		os.Exit(1)
	}

	cmd := exec.Command(
		"go",
		"tool",
		"swag",
		"init",
		"-g",
		"./cmd/main/main.go",
		"-ot",
		"json",
		"-o",
		"/tmp",
		"--pd",
	)
	cmd.Stdout = swaggerInfoLogger{}
	cmd.Stderr = swaggerInfoLogger{}
	err = cmd.Run()
	if err != nil {
		console.Error("Error running swag init: %s", err)
		os.Exit(1)
	}

	if outputDir == "" {
		console.Error("Output directory is required")
		os.Exit(1)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		console.Error("Error creating output directory: %s", err.Error())
		os.Exit(1)
	}

	// Read the swagger file
	swaggerFilePath := filepath.Join("/tmp", "swagger.json")
	swaggerFileBytes, err := os.ReadFile(swaggerFilePath)
	if err != nil {
		console.Error("Error reading swagger file: %s", err.Error())
		os.Exit(1)
	}

	var doc openapi2.T
	if err := json.Unmarshal(swaggerFileBytes, &doc); err != nil {
		console.Error("Error unmarshalling swagger file: %s", err.Error())
		os.Exit(1)
	}

	// Convert the swagger file to OpenAPI 3.0 specification
	spec, err := openapi2conv.ToV3(&doc)
	if err != nil {
		console.Error(
			"Error converting swagger file to OpenAPI 3.0 specification: %s",
			err.Error(),
		)
		os.Exit(1)
	}

	specJSON, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		console.Error("Error marshalling OpenAPI 3.0 specification: %s", err.Error())
		os.Exit(1)
	}

	specYaml, err := yaml.Marshal(spec)
	if err != nil {
		console.Error("Error marshalling OpenAPI 3.0 specification: %s", err.Error())
		os.Exit(1)
	}

	err = os.WriteFile(filepath.Join(outputDir, "openapi.json"), specJSON, 0644)
	if err != nil {
		console.Error(
			"Error writing OpenAPI 3.0 specification to JSON file: %s",
			err.Error(),
		)
		os.Exit(1)
	}

	err = os.WriteFile(filepath.Join(outputDir, "openapi.yaml"), specYaml, 0644)
	if err != nil {
		console.Error(
			"Error writing OpenAPI 3.0 specification to YAML file: %s",
			err.Error(),
		)
		os.Exit(1)
	}

	console.Info("OpenAPI 3.0 specification written to: %s", outputDir)
}

var (
	SwaggerCMD = &cobra.Command{
		Use:   "swagger",
		Short: "Generate Swagger file from Go code",
		Long:  "Generate Swagger file with OAS3 specification from Go code",
		Run: func(cobra *cobra.Command, args []string) {
			GenerateSwagger()
		},
	}
)

type swaggerInfoLogger struct{}

func (l swaggerInfoLogger) Write(p []byte) (n int, err error) {
	console.Info("%s", p)
	return len(p), nil
}

type swaggerErrorLogger struct{}

func (l swaggerErrorLogger) Write(p []byte) (n int, err error) {
	console.Error("%s", p)
	return len(p), nil
}
