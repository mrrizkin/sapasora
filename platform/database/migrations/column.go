package migrations

// Nullable sets the column to allow NULL values
func (c *Column) Nullable() *Column {
	c.IsNullable = true
	return c
}

// Unique adds a unique constraint to the column
func (c *Column) Unique() *Column {
	c.IsUnique = true
	return c
}

// Default sets the default value for the column
func (c *Column) Default(value any) *Column {
	c.DefaultVal = value
	return c
}

// Unsigned makes the column unsigned (MySQL only)
func (c *Column) Unsigned() *Column {
	c.IsUnsigned = true
	return c
}

// Comment adds a comment to the column
func (c *Column) Comment(comment string) *Column {
	c.ColComment = comment
	return c
}

// After specifies the column position (after another column - MySQL only)
func (c *Column) After(column string) *Column {
	c.PosAfter = column
	return c
}
