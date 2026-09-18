package migrations

import (
	"sapasora/platform/database"
	"fmt"
)

type Blueprint struct {
	TableName     string
	DB            *database.Database
	Columns       []*Column
	Indexes       []*Index
	Constraints   []*Constraint
	Commands      []string
	Modifications []string
}

func NewBlueprint(tableName string, db *database.Database) *Blueprint {
	return &Blueprint{
		TableName: tableName,
		DB:        db,
		Columns:   []*Column{},
		Indexes:   []*Index{},
	}
}

// ID creates an auto-incrementing BIGINT UNSIGNED primary key column
func (b *Blueprint) ID(name ...string) *Column {
	colName := "id"
	if len(name) > 0 && name[0] != "" {
		colName = name[0]
	}

	col := &Column{
		Name:       colName,
		DataType:   BIGINT,
		AutoInc:    true,
		Primary:    true,
		IsUnsigned: true, // For MySQL
	}

	publicIDCol := &Column{
		Name:       "public_id",
		DataType:   TEXT,
		Length:     255,
		IsUnique:   true,
		IsNullable: false,
		Primary:    false,
	}

	publicIDIndex := &Index{
		Name:    fmt.Sprintf("%s_public_id_index", b.TableName),
		Columns: []string{"public_id"},
		Type:    INDEX,
		Unique:  true,
	}

	b.Columns = append(b.Columns, col, publicIDCol)
	b.Indexes = append(b.Indexes, publicIDIndex)

	return col
}

// String creates a VARCHAR column
func (b *Blueprint) String(name string, length ...int) *Column {
	colLength := 255
	if len(length) > 0 {
		colLength = length[0]
	}

	col := &Column{
		Name:     name,
		DataType: VARCHAR,
		Length:   colLength,
	}

	b.Columns = append(b.Columns, col)
	return col
}

// Text creates a TEXT column
func (b *Blueprint) Text(name string) *Column {
	col := &Column{
		Name:     name,
		DataType: TEXT,
	}

	b.Columns = append(b.Columns, col)
	return col
}

// Integer creates an INT column
func (b *Blueprint) Integer(name string) *Column {
	col := &Column{
		Name:     name,
		DataType: INT,
	}

	b.Columns = append(b.Columns, col)
	return col
}

// BigInteger creates a BIGINT column
func (b *Blueprint) BigInteger(name string) *Column {
	col := &Column{
		Name:     name,
		DataType: BIGINT,
	}

	b.Columns = append(b.Columns, col)
	return col
}

// Boolean creates a BOOLEAN column
func (b *Blueprint) Boolean(name string) *Column {
	col := &Column{
		Name:     name,
		DataType: BOOLEAN,
	}

	b.Columns = append(b.Columns, col)
	return col
}

// Decimal creates a DECIMAL column with precision and scale
func (b *Blueprint) Decimal(name string, precision, scale int) *Column {
	col := &Column{
		Name:      name,
		DataType:  DECIMAL,
		Precision: precision,
		Scale:     scale,
	}

	b.Columns = append(b.Columns, col)
	return col
}

// Float creates a FLOAT column
func (b *Blueprint) Float(name string) *Column {
	col := &Column{
		Name:     name,
		DataType: FLOAT,
	}

	b.Columns = append(b.Columns, col)
	return col
}

// Double creates a DOUBLE column
func (b *Blueprint) Double(name string) *Column {
	col := &Column{
		Name:     name,
		DataType: DOUBLE,
	}

	b.Columns = append(b.Columns, col)
	return col
}

// Date creates a DATE column
func (b *Blueprint) Date(name string) *Column {
	col := &Column{
		Name:     name,
		DataType: DATE,
	}

	b.Columns = append(b.Columns, col)
	return col
}

// DateTime creates a DATETIME column
func (b *Blueprint) DateTime(name string) *Column {
	col := &Column{
		Name:     name,
		DataType: DATETIME,
	}

	b.Columns = append(b.Columns, col)
	return col
}

// Timestamp creates a TIMESTAMP column
func (b *Blueprint) Timestamp(name string) *Column {
	col := &Column{
		Name:     name,
		DataType: TIMESTAMP,
	}

	b.Columns = append(b.Columns, col)
	return col
}

// Timestamps creates created_at and updated_at columns (nullable timestamps)
func (b *Blueprint) Timestamps() {
	createdAtCol := &Column{
		Name:       "created_at",
		DataType:   TIMESTAMP,
		IsNullable: true,
	}
	updatedAtCol := &Column{
		Name:       "updated_at",
		DataType:   TIMESTAMP,
		IsNullable: true,
	}

	b.Columns = append(b.Columns, createdAtCol, updatedAtCol)
}

// SoftDeletes creates deleted_at column (nullable timestamp)
func (b *Blueprint) SoftDeletes() *Column {
	col := &Column{
		Name:       "deleted_at",
		DataType:   TIMESTAMP,
		IsNullable: true,
	}

	index := &Index{
		Name:    fmt.Sprintf("%s_deleted_at_index", b.TableName),
		Columns: []string{"deleted_at"},
		Type:    INDEX,
		Unique:  false,
	}

	b.Columns = append(b.Columns, col)
	b.Indexes = append(b.Indexes, index)
	return col
}

// JSON creates a JSON column
func (b *Blueprint) JSON(name string) *Column {
	dataType := JSON
	if b.DB.Name() == "postgres" {
		dataType = JSONB // Use JSONB for PostgreSQL
	}

	col := &Column{
		Name:     name,
		DataType: dataType,
	}

	b.Columns = append(b.Columns, col)
	return col
}

// Enum creates an ENUM column with allowed values
func (b *Blueprint) Enum(name string, values []string) *Column {
	col := &Column{
		Name:     name,
		DataType: ENUM,
		// Store the enum values in a way that can be accessed during SQL generation
		DefaultVal: values, // Store enum values as default value temporarily
	}

	b.Columns = append(b.Columns, col)
	return col
}

// Index creates a regular index
func (b *Blueprint) Index(columns ...string) *Blueprint {
	indexName := b.TableName + "_" + columns[0] + "_index"
	if len(columns) > 1 {
		indexName = b.TableName + "_" + columns[0]
		for i := 1; i < len(columns); i++ {
			indexName += "_" + columns[i]
		}
		indexName += "_index"
	}

	index := &Index{
		Name:    indexName,
		Columns: columns,
		Type:    INDEX,
		Unique:  false,
	}

	b.Indexes = append(b.Indexes, index)
	return b
}

// Unique creates a unique index
func (b *Blueprint) Unique(columns ...string) *Blueprint {
	indexName := b.TableName + "_" + columns[0] + "_unique"
	if len(columns) > 1 {
		indexName = b.TableName + "_" + columns[0]
		for i := 1; i < len(columns); i++ {
			indexName += "_" + columns[i]
		}
		indexName += "_unique"
	}

	index := &Index{
		Name:    indexName,
		Columns: columns,
		Type:    UNIQUE,
		Unique:  true,
	}

	b.Indexes = append(b.Indexes, index)
	return b
}

// Foreign starts defining a foreign key
func (b *Blueprint) Foreign(column string) *ForeignKeyBuilder {
	return &ForeignKeyBuilder{
		Column:    column,
		Blueprint: b,
	}
}

// DropColumn adds a command to drop a column
func (b *Blueprint) DropColumn(name string) *Blueprint {
	b.Commands = append(b.Commands, "DROP COLUMN "+name)
	return b
}

// RenameColumn adds a command to rename a column
func (b *Blueprint) RenameColumn(from, to string) *Blueprint {
	b.Commands = append(b.Commands, "RENAME COLUMN "+from+" TO "+to)
	return b
}

// DropIndex adds a command to drop an index
func (b *Blueprint) DropIndex(name string) *Blueprint {
	b.Commands = append(b.Commands, "DROP INDEX "+name)
	return b
}
