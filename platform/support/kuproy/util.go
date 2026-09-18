package kuproy

import (
	"fmt"
	"os"
	"runtime/debug"
)

func checkIfExists(check, file string) error {
	if info, err := os.Stat(file); err == nil {
		switch check {
		case "directory":
			if !info.IsDir() {
				return fmt.Errorf("'%s' exists but is not a directory", file)
			}

			entries, err := os.ReadDir(file)
			if err != nil {
				return fmt.Errorf("failed to read directory '%s': %w", file, err)
			}
			if len(entries) > 0 {
				return fmt.Errorf("output directory '%s' already exists and is not empty", file)
			}
		case "file":
			if info.IsDir() {
				return fmt.Errorf("'%s' exists but is a directory", file)
			}

			if info.Size() != 0 {
				return fmt.Errorf("'%s' exists but is not empty", file)
			}
		}

	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check output directory '%s': %w", file, err)
	}

	return nil
}

func prepareOutputDir(outputDir string) error {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory '%s': %w", outputDir, err)
	}

	return nil
}

func getModuleName() string {
	os, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	return os.Main.Path
}
