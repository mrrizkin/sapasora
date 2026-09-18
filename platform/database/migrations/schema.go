package migrations

import (
	"sapasora/platform/database"

	"gorm.io/gorm"
)

// Schema is the main entry point for migrations, providing a fluent interface for table operations
type Schema struct {
	db       *database.Database
	migrator gorm.Migrator
}

// NewSchema creates a new Schema instance with the provided database connection
func NewSchema(db *database.Database) *Schema {
	return &Schema{
		db:       db,
		migrator: db.Migrator(),
	}
}

// Create creates a new table using the provided callback function
func (s *Schema) Create(tableName string, callback func(*Blueprint)) error {
	blueprint := NewBlueprint(tableName, s.db)
	callback(blueprint)
	return blueprint.Execute()
}

// Table modifies an existing table using the provided callback function
func (s *Schema) Table(tableName string, callback func(*Blueprint)) error {
	blueprint := NewBlueprint(tableName, s.db)
	callback(blueprint)
	return blueprint.Execute()
}

// Drop drops a table
func (s *Schema) Drop(tableName string) error {
	return s.migrator.DropTable(tableName)
}

// DropIfExists drops a table if it exists
func (s *Schema) DropIfExists(tableName string) error {
	return s.migrator.DropTable(tableName)
}

// Rename renames a table from 'from' to 'to'
func (s *Schema) Rename(from, to string) error {
	return s.migrator.RenameTable(from, to)
}

// HasTable checks if a table exists
func (s *Schema) HasTable(tableName string) bool {
	return s.migrator.HasTable(tableName)
}

// HasColumn checks if a column exists in a table
func (s *Schema) HasColumn(tableName, columnName string) bool {
	return s.migrator.HasColumn(tableName, columnName)
}

// GetTables returns all table names in the database
func (s *Schema) GetTables() ([]string, error) {
	return s.migrator.GetTables()
}

// Exec executes a raw SQL query
func (s *Schema) Exec(query string, args ...any) error {
	return s.db.Exec(query, args...).Error
}

// GetDialector returns the DBMS type (mysql, postgres, sqlite, sqlserver)
func (s *Schema) GetDialector() string {
	return s.db.Name()
}

// IsMysql checks if database is MySQL
func (s *Schema) IsMysql() bool {
	return s.GetDialector() == "mysql"
}

// IsPostgres checks if database is PostgreSQL
func (s *Schema) IsPostgres() bool {
	return s.GetDialector() == "postgres"
}

// IsSqlite checks if database is SQLite
func (s *Schema) IsSqlite() bool {
	return s.GetDialector() == "sqlite"
}

// IsSQLServer checks if database is SQL Server
func (s *Schema) IsSQLServer() bool {
	return s.GetDialector() == "sqlserver"
}

// CurrentDatabase returns the current database name
func (s *Schema) CurrentDatabase() string {
	return s.migrator.CurrentDatabase()
}
