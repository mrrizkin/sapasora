package kuproy

type KuproyMap map[string]any

type Field struct {
	Name       string
	Type       string
	Attributes []string // unique, required, etc.
}

type Generator struct {
	Module    string
	Entity    string
	TableName string
	Fields    []Field

	ModulePath string

	MigrationName string
}
