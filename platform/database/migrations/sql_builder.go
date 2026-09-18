package migrations

import (
	"fmt"
	"strings"
)

// Execute executes the blueprint operations (table creation/modification)
func (b *Blueprint) Execute() error {
	if len(b.Commands) > 0 {
		return b.executeModifications()
	}

	return b.createTable()
}

// executeModifications executes table modification commands
func (b *Blueprint) executeModifications() error {
	dialect := b.DB.Name()

	for _, command := range b.Commands {
		sql := fmt.Sprintf("ALTER TABLE %s %s", quoteIdentifier(b.TableName, dialect), command)
		if err := b.DB.Exec(sql).Error; err != nil {
			return err
		}
	}

	return nil
}

// createTable creates a new table based on the blueprint
func (b *Blueprint) createTable() error {
	dialect := b.DB.Name()

	sql := b.buildCreateTableSQL(dialect)

	if err := b.DB.Exec(sql).Error; err != nil {
		return err
	}

	for _, index := range b.Indexes {
		indexSQL := b.buildIndexSQL(index, dialect)
		if err := b.DB.Exec(indexSQL).Error; err != nil {
			return err
		}
	}

	return nil
}

// buildCreateTableSQL builds a dialect-specific CREATE TABLE SQL statement
func (b *Blueprint) buildCreateTableSQL(dialect string) string {
	var parts []string

	parts = append(parts, fmt.Sprintf("CREATE TABLE %s (", quoteIdentifier(b.TableName, dialect)))

	var primaryCols []string
	for i, col := range b.Columns {
		if col.Primary {
			primaryCols = append(primaryCols, quoteIdentifier(col.Name, dialect))
		}

		colSQL := b.buildColumnSQL(col, dialect)
		if i < len(b.Columns)-1 {
			colSQL += ","
		}
		parts = append(parts, colSQL)
	}

	// Add a comma separator if there are columns and we're adding more constraints
	if len(b.Columns) > 0 && (len(b.Constraints) > 0) {
		parts = append(parts, ",")
	}

	// Add foreign key constraints
	for i, constraint := range b.Constraints {
		if constraint.Type == FOREIGN {
			// Generate the foreign key constraint definition from the constraint fields
			dialect := b.DB.Name()
			constraintSQL := generateForeignKeyConstraintSQL(constraint, dialect)
			if i < len(b.Constraints)-1 {
				parts = append(parts, constraintSQL+",")
			} else {
				parts = append(parts, constraintSQL)
			}
		}
	}

	if dialect != "sqlite" {
		if len(primaryCols) > 0 {
			// Add comma if there are other constraints and no foreign key constraints
			if len(b.Columns) > 0 {
				parts = append(parts, ",")
			}
			parts = append(parts, fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(primaryCols, ", ")))
		}
	}

	parts = append(parts, ")")

	if dialect == "mysql" {
		parts = append(parts, "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")
	}

	return strings.Join(parts, "\n")
}

// buildColumnSQL generates the SQL for a single column based on dialect
func (b *Blueprint) buildColumnSQL(col *Column, dialect string) string {
	var parts []string

	parts = append(parts, quoteIdentifier(col.Name, dialect))

	dataType := getDialectDataType(col, dialect)
	parts = append(parts, dataType)

	if col.Length > 0 && !isDataTypeLengthIgnored(col.DataType, dialect) {
		if col.DataType == ENUM {
			// Handle ENUM specially - values come from DefaultVal as a string slice
			if values, ok := col.DefaultVal.([]string); ok {
				enumValues := make([]string, len(values))
				for i, v := range values {
					enumValues[i] = fmt.Sprintf("'%s'", v)
				}
				parts[len(parts)-1] = fmt.Sprintf(
					"ENUM(%s)",
					strings.Join(enumValues, ", "),
				) // Replace the data type
			}
		} else {
			parts[len(parts)-1] += fmt.Sprintf("(%d)", col.Length)
		}
	}

	if col.DataType == DECIMAL && col.Precision > 0 {
		parts[len(parts)-1] = fmt.Sprintf("%s(%d,%d)", dataType, col.Precision, col.Scale)
	}

	if dialect != "sqlite" || !col.Primary || !col.AutoInc {
		if !col.IsNullable {
			parts = append(parts, "NOT NULL")
		} else if dialect != "postgres" || col.DataType == TIMESTAMP {
			parts = append(parts, "NULL")
		}
	}

	if col.AutoInc {
		switch dialect {
		case "mysql":
			parts = append(parts, "AUTO_INCREMENT")
		case "postgres":
			// For postgres, we change the data type to SERIAL or BIGSERIAL instead of adding AUTO_INCREMENT
			// This is handled in getDialectDataType
		case "sqlite":
			parts = append(parts, "PRIMARY KEY AUTOINCREMENT") // Only for primary keys in SQLite
		case "sqlserver":
			parts = append(parts, "IDENTITY(1,1)")
		}
	}

	if col.IsUnsigned && dialect == "mysql" {
		parts = append(parts, "UNSIGNED")
	}

	if col.DefaultVal != nil {
		defaultVal := col.DefaultVal
		if col.DataType != ENUM {
			if s, ok := defaultVal.(string); ok {
				if strings.ToUpper(s) == "CURRENT_TIMESTAMP" ||
					strings.Contains(strings.ToUpper(s), "NOW()") {
					parts = append(parts, fmt.Sprintf("DEFAULT %s", s))
				} else {
					parts = append(parts, fmt.Sprintf("DEFAULT '%s'", s))
				}
			} else {
				parts = append(parts, fmt.Sprintf("DEFAULT %v", defaultVal))
			}
		}
	}

	if col.ColComment != "" {
		switch dialect {
		case "mysql":
			parts = append(parts, fmt.Sprintf("COMMENT '%s'", col.ColComment))
		case "postgres":
			// Comment is added separately after column creation in PostgreSQL
		case "sqlite":
			// SQLite doesn't support column comments
		case "sqlserver":
			// Comments in SQL Server are added using extended properties
		}
	}

	return strings.Join(parts, " ")
}

// buildIndexSQL generates the SQL for creating an index
func (b *Blueprint) buildIndexSQL(index *Index, dialect string) string {
	var parts []string

	indexType := ""
	if index.Unique {
		indexType = "UNIQUE "
	}

	tableName := quoteIdentifier(b.TableName, dialect)
	indexName := quoteIdentifier(index.Name, dialect)

	columns := joinColumns(index.Columns, dialect)

	parts = append(
		parts,
		fmt.Sprintf("CREATE %sINDEX %s ON %s (%s)", indexType, indexName, tableName, columns),
	)

	return strings.Join(parts, " ")
}

// getDialectDataType maps the generic data type to dialect-specific types
func getDialectDataType(col *Column, dialect string) string {
	dataType := col.DataType

	// Handle auto-incrementing types differently
	if col.AutoInc {
		switch dialect {
		case "postgres":
			switch dataType {
			case BIGINT:
				return "BIGSERIAL"
			case INT:
				return "SERIAL"
			}
		case "sqlite":
			switch dataType {
			case BIGINT:
				return "INTEGER"
			case INT:
				return "INTEGER"
			}
		}
	}

	// Map generic types to dialect-specific types
	switch dataType {
	case BIGINT:
		if dialect == "postgres" && col.AutoInc {
			return "BIGSERIAL"
		}
		return dataType.String()
	case INT:
		if dialect == "postgres" && col.AutoInc {
			return "SERIAL"
		}
		if dialect == "postgres" {
			return "INTEGER"
		}
		return dataType.String()
	case VARCHAR:
		if dialect == "sqlserver" {
			return "NVARCHAR"
		}
		return dataType.String()
	case TEXT:
		if dialect == "sqlite" {
			return "TEXT" // SQLite's TEXT is equivalent
		}
		if dialect == "sqlserver" {
			return "NVARCHAR(MAX)"
		}
		return dataType.String()
	case BOOLEAN:
		if dialect == "sqlite" {
			return "INTEGER" // SQLite doesn't have a native boolean type
		}
		if dialect == "sqlserver" {
			return "BIT"
		}
		return dataType.String()
	case DATETIME:
		if dialect == "postgres" {
			return "TIMESTAMP"
		}
		if dialect == "sqlserver" {
			return "DATETIME2"
		}
		return dataType.String()
	case TIMESTAMP:
		if dialect == "postgres" {
			return "TIMESTAMP"
		}
		if dialect == "sqlserver" {
			return "DATETIME2"
		}
		return dataType.String()
	case JSON:
		if dialect == "postgres" {
			return "JSONB"
		}
		if dialect == "sqlite" {
			return "TEXT" // SQLite doesn't have native JSON support in older versions
		}
		if dialect == "sqlserver" {
			return "NVARCHAR(MAX)" // SQL Server represents JSON as NVARCHAR(MAX)
		}
		return dataType.String()
	case FLOAT:
		if dialect == "postgres" {
			return "REAL" // REAL is a 4-byte floating-point number in PostgreSQL
		}
		return dataType.String()
	case DOUBLE:
		if dialect == "postgres" {
			return "DOUBLE PRECISION"
		}
		return dataType.String()
	default:
		return dataType.String()
	}
}

// isDataTypeLengthIgnored checks if the length should be omitted for certain data types in specific dialects
func isDataTypeLengthIgnored(dataType ColumnType, dialect string) bool {
	switch {
	case dataType == TEXT:
		return true
	case dataType == BLOB && dialect == "sqlite":
		return true
	case dataType == CLOB && dialect == "sqlite":
		return true
	case dataType == JSON:
		return true
	case dataType == BOOLEAN && dialect == "sqlite":
		return true
	default:
		return false
	}
}

// quoteIdentifier adds appropriate quotes around identifiers based on the dialect
func quoteIdentifier(identifier, dialect string) string {
	switch dialect {
	case "mysql":
		return fmt.Sprintf("`%s`", identifier)
	case "postgres":
		return fmt.Sprintf(`"%s"`, identifier)
	case "sqlite":
		return fmt.Sprintf(`"%s"`, identifier)
	case "sqlserver":
		return fmt.Sprintf(`[%s]`, identifier)
	default:
		return fmt.Sprintf(`"%s"`, identifier)
	}
}

// joinColumns joins column names with proper quoting for the dialect
func joinColumns(cols []string, dialect string) string {
	quotedCols := make([]string, len(cols))
	for i, col := range cols {
		quotedCols[i] = quoteIdentifier(col, dialect)
	}
	return strings.Join(quotedCols, ", ")
}

// generateForeignKeyConstraintSQL creates the SQL for a foreign key constraint
func generateForeignKeyConstraintSQL(constraint *Constraint, dialect string) string {
	if constraint.Type != FOREIGN {
		return ""
	}

	var constraintSQL string
	switch dialect {
	case "mysql":
		constraintSQL = fmt.Sprintf(
			"CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)",
			quoteIdentifier(constraint.Name, dialect),
			quoteIdentifier(strings.Join(constraint.Columns, ", "), dialect),
			quoteIdentifier(constraint.Reference, dialect),
			quoteIdentifier(strings.Join(constraint.RefColumns, ", "), dialect),
		)
	case "postgres":
		constraintSQL = fmt.Sprintf(
			"CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)",
			quoteIdentifier(constraint.Name, dialect),
			quoteIdentifier(strings.Join(constraint.Columns, ", "), dialect),
			quoteIdentifier(constraint.Reference, dialect),
			quoteIdentifier(strings.Join(constraint.RefColumns, ", "), dialect),
		)
	case "sqlite":
		// For SQLite, foreign key constraints are defined inline with the column or as table constraints
		constraintSQL = fmt.Sprintf(
			"CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)",
			quoteIdentifier(constraint.Name, dialect),
			quoteIdentifier(strings.Join(constraint.Columns, ", "), dialect),
			quoteIdentifier(constraint.Reference, dialect),
			quoteIdentifier(strings.Join(constraint.RefColumns, ", "), dialect),
		)
	case "sqlserver":
		constraintSQL = fmt.Sprintf(
			"CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)",
			quoteIdentifier(constraint.Name, dialect),
			quoteIdentifier(strings.Join(constraint.Columns, ", "), dialect),
			quoteIdentifier(constraint.Reference, dialect),
			quoteIdentifier(strings.Join(constraint.RefColumns, ", "), dialect),
		)
	default:
		constraintSQL = fmt.Sprintf(
			"CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)",
			quoteIdentifier(constraint.Name, dialect),
			quoteIdentifier(strings.Join(constraint.Columns, ", "), dialect),
			quoteIdentifier(constraint.Reference, dialect),
			quoteIdentifier(strings.Join(constraint.RefColumns, ", "), dialect),
		)
	}

	// Add ON DELETE and ON UPDATE options if specified
	if constraint.OnDelete != "" {
		constraintSQL += fmt.Sprintf(" ON DELETE %s", constraint.OnDelete)
	}
	if constraint.OnUpdate != "" {
		constraintSQL += fmt.Sprintf(" ON UPDATE %s", constraint.OnUpdate)
	}

	return constraintSQL
}
