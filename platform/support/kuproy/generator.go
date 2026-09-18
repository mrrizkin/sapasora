package kuproy

import (
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func (g *Generator) Policy(filePath string, t *Templates) error {
	outputDir := filepath.Dir(filePath)
	packageName := strings.ToLower(filepath.Base(outputDir))
	prepareOutputDir(outputDir)
	return t.Write(filePath, "policy.gotmpl", KuproyMap{
		"Name":        g.Module,
		"StructName":  g.Entity,
		"PackageName": packageName,
		"ModulePath":  g.ModulePath,
	})
}

func (g *Generator) Model(filePath string, t *Templates) error {
	outputDir := filepath.Dir(filePath)
	packageName := strings.ToLower(filepath.Base(outputDir))
	prepareOutputDir(outputDir)
	return t.Write(filePath, "model.gotmpl", KuproyMap{
		"Entity":      g.Entity,
		"TableName":   g.TableName,
		"Fields":      g.Fields,
		"PackageName": packageName,
	})
}

func (g *Generator) Controller(filePath string, t *Templates) error {
	outputDir := filepath.Dir(filePath)
	packageName := strings.ToLower(filepath.Base(outputDir))
	prepareOutputDir(outputDir)
	return t.Write(filePath, "controller.gotmpl", KuproyMap{
		"Entity":      g.Entity,
		"PackageName": packageName,
		"Fields":      g.Fields,
		"ModulePath":  g.ModulePath,
	})
}

func (g *Generator) ControllerType(filePath string, t *Templates) error {
	outputDir := filepath.Dir(filePath)
	packageName := strings.ToLower(filepath.Base(outputDir))
	prepareOutputDir(outputDir)
	return t.Write(filePath, "controller_type.gotmpl", KuproyMap{
		"Entity":      g.Entity,
		"PackageName": packageName,
		"Fields":      g.Fields,
		"ModulePath":  g.ModulePath,
	})
}

func (g *Generator) Interface(filePath string, t *Templates) error {
	outputDir := filepath.Dir(filePath)
	packageName := strings.ToLower(filepath.Base(outputDir))
	prepareOutputDir(outputDir)
	return t.Write(filePath, "interface.gotmpl", KuproyMap{
		"Entity":      g.Entity,
		"PackageName": packageName,
	})
}

func (g *Generator) Repository(filePath string, t *Templates) error {
	outputDir := filepath.Dir(filePath)
	packageName := strings.ToLower(filepath.Base(outputDir))
	prepareOutputDir(outputDir)
	return t.Write(filePath, "repository.gotmpl", KuproyMap{
		"Entity":      g.Entity,
		"PackageName": packageName,
		"ModulePath":  g.ModulePath,
	})
}

func (g *Generator) Service(filePath string, t *Templates) error {
	outputDir := filepath.Dir(filePath)
	packageName := strings.ToLower(filepath.Base(outputDir))
	prepareOutputDir(outputDir)
	return t.Write(filePath, "service.gotmpl", KuproyMap{
		"Entity":      g.Entity,
		"PackageName": packageName,
	})
}

func (g *Generator) Migration(filePath string, t *Templates) error {
	outputDir := filepath.Dir(filePath)
	packageName := strings.ToLower(filepath.Base(outputDir))
	migrationType := g.determineMigrationType(g.Entity)
	tableName := g.extractTableName(g.Entity, migrationType)
	prepareOutputDir(outputDir)
	return t.Write(filePath, "migration.gotmpl", KuproyMap{
		"StructName":    g.convertToStructName(g.Entity),
		"Name":          g.MigrationName,
		"TableName":     tableName,
		"PackageName":   packageName,
		"Fields":        g.Fields,
		"MigrationType": migrationType,
		"ModulePath":    g.ModulePath,
	})
}

func (g *Generator) determineMigrationType(name string) string {
	lowerName := strings.ToLower(name)

	if strings.HasPrefix(lowerName, "create_") && strings.HasSuffix(lowerName, "_table") {
		return "create_table"
	} else if strings.HasPrefix(lowerName, "add_") && strings.Contains(lowerName, "_to_") {
		return "add_columns"
	} else if strings.HasPrefix(lowerName, "drop_") && strings.HasSuffix(lowerName, "_table") {
		return "drop_table"
	} else if strings.HasPrefix(lowerName, "remove_") && strings.Contains(lowerName, "_from_") {
		return "remove_columns"
	} else {
		return "generic"
	}
}

func (g *Generator) convertToStructName(name string) string {
	parts := strings.Split(name, "_")
	title := cases.Title(language.English)
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = title.String(part)
		}
	}

	return strings.Join(parts, "")
}

func (g *Generator) extractTableName(migrationName, migrationType string) string {
	switch migrationType {
	case "create_table":
		// "create_users_table" -> "users"
		parts := strings.Split(migrationName, "_")
		if len(parts) >= 2 {
			return strings.Join(parts[1:len(parts)-1], "_")
		}
	case "add_columns", "remove_columns":
		// "add_name_to_users" -> "users"
		parts := strings.Split(migrationName, "_")
		if len(parts) >= 3 && parts[len(parts)-2] == "to" {
			return parts[len(parts)-1]
		}
	case "drop_table":
		// "drop_users_table" -> "users"
		parts := strings.Split(migrationName, "_")
		if len(parts) >= 2 {
			return strings.Join(parts[1:len(parts)-1], "_")
		}
	}

	parts := strings.Split(migrationName, "_")
	if len(parts) > 0 {
		for i, part := range parts {
			if part == "table" && i > 0 {
				return parts[i-1]
			}
			if part == "to" && i > 0 && i < len(parts)-1 {
				return parts[i+1]
			}
		}
	}

	return "table_name"
}
