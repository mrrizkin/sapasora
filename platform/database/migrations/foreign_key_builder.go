package migrations

import (
	"sapasora/platform/support/console"
	"fmt"
)

// References specifies the referenced column in the foreign key relationship
func (fkb *ForeignKeyBuilder) References(column string) *ForeignKeyBuilder {
	fkb.Reference = column
	return fkb
}

// On specifies the referenced table in the foreign key relationship
func (fkb *ForeignKeyBuilder) On(table string) *ForeignKeyBuilder {
	fkb.Table = table
	return fkb
}

// OnDelete specifies the action to take when the referenced record is deleted
func (fkb *ForeignKeyBuilder) OnDelete(action string) *ForeignKeyBuilder {
	fkb.RefOnDelete = action
	return fkb
}

// OnUpdate specifies the action to take when the referenced record is updated
func (fkb *ForeignKeyBuilder) OnUpdate(action string) *ForeignKeyBuilder {
	fkb.RefOnUpdate = action
	return fkb
}

func (fkb *ForeignKeyBuilder) Build() error {
	if fkb.Reference == "" || fkb.Table == "" {
		return fmt.Errorf("foreign key must specify both reference column and table")
	}

	// Add the constraint to the blueprint's constraints list so it gets included in CREATE TABLE
	constraintName := fmt.Sprintf("%s_%s_foreign", fkb.Blueprint.TableName, fkb.Column)

	fkb.Blueprint.Constraints = append(fkb.Blueprint.Constraints, &Constraint{
		Name:       constraintName,
		Type:       FOREIGN,
		Table:      fkb.Blueprint.TableName,
		Columns:    []string{fkb.Column},
		Reference:  fkb.Table,
		RefColumns: []string{fkb.Reference},
		OnDelete:   fkb.RefOnDelete,
		OnUpdate:   fkb.RefOnUpdate,
	})

	// For SQLite, the foreign key needs to be handled differently since ALTER TABLE ADD CONSTRAINT
	// is not well supported. The constraint will be added in CREATE TABLE.
	if fkb.Blueprint.DB.Name() == "sqlite" {
		console.Warn(
			"SQLite has limited support for ALTER TABLE ADD CONSTRAINT, foreign key defined at table creation",
		)
		return nil
	}

	return nil
}

// ExecuteWithTableDefinition adds the foreign key as a constraint in CREATE TABLE statement
// This is used when creating tables with foreign keys rather than altering existing tables
func (fkb *ForeignKeyBuilder) ExecuteWithTableDefinition() string {
	if fkb.Reference == "" || fkb.Table == "" {
		return ""
	}

	// Generate the foreign key constraint name
	constraintName := fmt.Sprintf("%s_%s_foreign", fkb.Blueprint.TableName, fkb.Column)

	// Generate the foreign key clause for CREATE TABLE
	dialect := fkb.Blueprint.DB.Name()

	var constraint string
	switch dialect {
	case "mysql":
		constraint = fmt.Sprintf(
			"CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)",
			quoteIdentifier(constraintName, dialect),
			quoteIdentifier(fkb.Column, dialect),
			quoteIdentifier(fkb.Table, dialect),
			quoteIdentifier(fkb.Reference, dialect),
		)
	case "postgres":
		constraint = fmt.Sprintf(
			"CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)",
			quoteIdentifier(constraintName, dialect),
			quoteIdentifier(fkb.Column, dialect),
			quoteIdentifier(fkb.Table, dialect),
			quoteIdentifier(fkb.Reference, dialect),
		)
	case "sqlite":
		// For SQLite, foreign key constraints are defined inline with the column or as table constraints
		constraint = fmt.Sprintf(
			"CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)",
			quoteIdentifier(constraintName, dialect),
			quoteIdentifier(fkb.Column, dialect),
			quoteIdentifier(fkb.Table, dialect),
			quoteIdentifier(fkb.Reference, dialect),
		)
	case "sqlserver":
		constraint = fmt.Sprintf(
			"CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)",
			quoteIdentifier(constraintName, dialect),
			quoteIdentifier(fkb.Column, dialect),
			quoteIdentifier(fkb.Table, dialect),
			quoteIdentifier(fkb.Reference, dialect),
		)
	default:
		constraint = fmt.Sprintf(
			"CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)",
			quoteIdentifier(constraintName, dialect),
			quoteIdentifier(fkb.Column, dialect),
			quoteIdentifier(fkb.Table, dialect),
			quoteIdentifier(fkb.Reference, dialect),
		)
	}

	if fkb.RefOnDelete != "" {
		constraint += fmt.Sprintf(" ON DELETE %s", fkb.RefOnDelete)
	}
	if fkb.RefOnUpdate != "" {
		constraint += fmt.Sprintf(" ON UPDATE %s", fkb.RefOnUpdate)
	}

	return constraint
}
