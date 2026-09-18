// Package kuproy provides a generator command to bootstrap code based on the domain provided
package kuproy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// “Kuproy” is short for “Kuli Proyek” — imagine “Bob the Builder,” but with flip-flops and no helmet.

type Kuproy struct {
	Task string
	Args []string
}

func (k *Kuproy) GenerateCode() error {
	if k.Task == "" {
		return fmt.Errorf("task is required")
	}

	if k.Args == nil {
		return fmt.Errorf("not enough parameters")
	} else if len(k.Args) == 0 {
		return fmt.Errorf("not enough parameters")
	}

	gen, err := parseArgs(k)
	if err != nil {
		return err
	}

	tmpl := NewTemplates()

	switch k.Task {
	case "auth":
		return nil
	case "model":
		if err := k.checkModelFiles(gen); err != nil {
			return err
		}
		if err := k.generateModel(gen, tmpl); err != nil {
			return err
		}
	case "context":
		if err := k.checkContextFiles(gen); err != nil {
			return err
		}
		if err := k.generateContext(gen, tmpl); err != nil {
			return err
		}
	case "json":
		if err := k.checkContextFiles(gen); err != nil {
			return err
		}
		if err := k.checkControllerFiles(gen); err != nil {
			return err
		}
		if err := k.generateContext(gen, tmpl); err != nil {
			return err
		}
		if err := k.generateController(gen, tmpl); err != nil {
			return err
		}
		return nil
	case "policy":
		if err := k.checkPolicyFiles(gen); err != nil {
			return err
		}
		if err := k.generatePolicy(gen, tmpl); err != nil {
			return err
		}
	case "migration":
		if err := k.checkMigrationFiles(gen); err != nil {
			return err
		}
		if err := k.generateMigration(gen, tmpl); err != nil {
			return err
		}
	default:
		fmt.Fprintf(os.Stderr, "Error: invalid task '%s'\n", k.Task)
		return fmt.Errorf("invalid task")
	}

	return nil
}

func (k *Kuproy) modelFiles(gen *Generator) string {
	outputDir := "internal/modules/"
	modelFilePath := filepath.Join(
		outputDir,
		strings.ToLower(gen.Module),
		strings.ToLower(gen.Entity)+".go",
	)

	return modelFilePath
}

func (k *Kuproy) checkModelFiles(gen *Generator) error {
	modelFilePath := k.modelFiles(gen)

	if err := k.checkFiles(
		modelFilePath,
	); err != nil {
		return err
	}
	return nil
}

func (k *Kuproy) generateModel(gen *Generator, tmpl *Templates) error {
	modelFilePath := k.modelFiles(gen)
	if err := gen.Model(modelFilePath, tmpl); err != nil {
		return err
	}
	return nil
}

func (k *Kuproy) policyFiles(gen *Generator) string {
	outputDir := "internal/app/policies"
	filePath := filepath.Join(outputDir, strings.ToLower(gen.Module)+".go")

	return filePath
}

func (k *Kuproy) checkPolicyFiles(gen *Generator) error {
	filePath := k.policyFiles(gen)

	if err := k.checkFiles(
		filePath,
	); err != nil {
		return err
	}
	return nil
}

func (k *Kuproy) generatePolicy(gen *Generator, tmpl *Templates) error {
	filePath := k.policyFiles(gen)
	if err := gen.Policy(filePath, tmpl); err != nil {
		return err
	}
	return nil
}

func (k *Kuproy) contextFiles(gen *Generator) (string, string, string, string) {
	outputDir := "internal/modules/"
	modelFilePath := k.modelFiles(gen)
	interfaceFilePath := filepath.Join(outputDir, strings.ToLower(gen.Module), "interface.go")
	repositoryFilePath := filepath.Join(outputDir, strings.ToLower(gen.Module), "repository.go")
	serviceFilePath := filepath.Join(outputDir, strings.ToLower(gen.Module), "service.go")

	return modelFilePath, interfaceFilePath, repositoryFilePath, serviceFilePath
}

func (k *Kuproy) checkContextFiles(gen *Generator) error {
	modelFilePath, interfaceFilePath, repositoryFilePath, serviceFilePath := k.contextFiles(gen)

	if err := k.checkFiles(
		modelFilePath,
		interfaceFilePath,
		repositoryFilePath,
		serviceFilePath,
	); err != nil {
		return err
	}

	return nil
}

func (k *Kuproy) generateContext(gen *Generator, tmpl *Templates) error {
	modelFilePath, interfaceFilePath, repositoryFilePath, serviceFilePath := k.contextFiles(gen)

	if err := gen.Model(modelFilePath, tmpl); err != nil {
		return err
	}
	if err := gen.Interface(interfaceFilePath, tmpl); err != nil {
		return err
	}
	if err := gen.Repository(repositoryFilePath, tmpl); err != nil {
		return err
	}
	if err := gen.Service(serviceFilePath, tmpl); err != nil {
		return err
	}

	return nil
}

func (k *Kuproy) controllerFiles(gen *Generator) (string, string) {
	outputDir := "internal/app/http/controllers"
	controllerFilePath := filepath.Join(
		outputDir,
		strings.ToLower(gen.Module),
		strings.ToLower(gen.Entity)+".go",
	)
	typeFilePath := filepath.Join(outputDir, strings.ToLower(gen.Module), "type.go")

	return controllerFilePath, typeFilePath
}

func (k *Kuproy) checkControllerFiles(gen *Generator) error {
	controllerFilePath, typeFilePath := k.controllerFiles(gen)

	if err := k.checkFiles(
		controllerFilePath,
		typeFilePath,
	); err != nil {
		return err
	}

	return nil
}

func (k *Kuproy) generateController(gen *Generator, tmpl *Templates) error {
	controllerFilePath, typeFilePath := k.controllerFiles(gen)

	if err := gen.Controller(controllerFilePath, tmpl); err != nil {
		return err
	}
	if err := gen.ControllerType(typeFilePath, tmpl); err != nil {
		return err
	}

	return nil
}

func (k *Kuproy) migrationFiles(gen *Generator) string {
	outputDir := "internal/migrations"
	filePath := filepath.Join(outputDir, gen.MigrationName)

	return filePath
}

func (k *Kuproy) checkMigrationFiles(gen *Generator) error {
	filePath := k.migrationFiles(gen)

	if err := k.checkFiles(
		filePath,
	); err != nil {
		return err
	}
	return nil
}

func (k *Kuproy) generateMigration(gen *Generator, tmpl *Templates) error {
	filePath := k.migrationFiles(gen)
	if err := gen.Migration(filePath, tmpl); err != nil {
		return err
	}
	return nil
}

func (k *Kuproy) checkFiles(filePaths ...string) error {
	for _, filePath := range filePaths {
		if err := checkIfExists("file", filePath); err != nil {
			return err
		}
	}
	return nil
}
