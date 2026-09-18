// Package migrations provides a fluent, expressive database migration system for Go
package migrations

// Migration interface that all migrations must implement
type Migration interface {
	Name() string
	Up(schema *Schema) error
	Down(schema *Schema) error
}

type ColumnType int
type IndexType int
type ConstraintType int

const (
	BIGINT ColumnType = iota
	INT
	VARCHAR
	TEXT
	BOOLEAN
	DATETIME
	TIMESTAMP
	JSON
	ENUM
	DECIMAL
	FLOAT
	DOUBLE
	DATE
	JSONB
	BLOB
	TIME
	CLOB
)

const (
	// IndexType represents the type of index
	INDEX IndexType = iota
	UNIQUE
)

const (
	// ConstraintType represents the type of constraint
	PRIMARY ConstraintType = iota
	FOREIGN
)

func (ct ColumnType) String() string {
	switch ct {
	case BIGINT:
		return "BIGINT"
	case INT:
		return "INT"
	case VARCHAR:
		return "VARCHAR"
	case TEXT:
		return "TEXT"
	case BOOLEAN:
		return "BOOLEAN"
	case DATETIME:
		return "DATETIME"
	case TIMESTAMP:
		return "TIMESTAMP"
	case JSON:
		return "JSON"
	case ENUM:
		return "ENUM"
	case DECIMAL:
		return "DECIMAL"
	case FLOAT:
		return "FLOAT"
	case DOUBLE:
		return "DOUBLE"
	case DATE:
		return "DATE"
	case JSONB:
		return "JSONB"
	case BLOB:
		return "BLOB"
	case TIME:
		return "TIME"
	case CLOB:
		return "CLOB"
	}

	return "UNKNOWN"
}

func (it IndexType) String() string {
	switch it {
	case INDEX:
		return "INDEX"
	case UNIQUE:
		return "UNIQUE"
	default:
		return "UNKNOWN"
	}
}

func (ct ConstraintType) String() string {
	switch ct {
	case PRIMARY:
		return "PRIMARY"
	case FOREIGN:
		return "FOREIGN"
	default:
		return "UNKNOWN"
	}
}

// Column represents a database column definition
type Column struct {
	Name       string     // Column name
	DataType   ColumnType // Data type (e.g., VARCHAR, INT, etc.)
	Length     int        // Length for variable-length types
	Precision  int        // Precision for decimal types
	Scale      int        // Scale for decimal types
	IsNullable bool       // Whether the column allows NULL values
	IsUnique   bool       // Whether the column has a unique constraint
	Primary    bool       // Whether the column is part of primary key
	AutoInc    bool       // Whether the column is auto-incrementing
	DefaultVal any        // Default value for the column
	ColComment string     // Column comment
	IsUnsigned bool       // Whether the column is unsigned (MySQL only)
	PosAfter   string     // Column position (MySQL only)
	Charset    string     // Character set
	Collation  string     // Collation
}

// Index represents a database index
type Index struct {
	Name    string    // Index name
	Columns []string  // Columns in the index
	Type    IndexType // Index type (e.g., INDEX, UNIQUE, PRIMARY)
	Unique  bool      // Whether the index is unique
}

// Constraint represents a database constraint
type Constraint struct {
	Name       string         // Constraint name
	Type       ConstraintType // Constraint type (e.g., FOREIGN_KEY, CHECK)
	Table      string         // Table name
	Columns    []string
	Reference  string // Reference table
	RefColumns []string
	OnDelete   string
	OnUpdate   string
	Deferrable bool
	Initially  string
}

// ForeignKeyBuilder provides a fluent interface for defining foreign key constraints
type ForeignKeyBuilder struct {
	Column      string
	Reference   string
	Table       string
	RefOnDelete string
	RefOnUpdate string
	Blueprint   *Blueprint
}
