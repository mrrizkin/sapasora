// Package migrations provides the migrations for the database
package migrations

import "sapasora/platform/database/migrations"

type CreateMRolesTable struct{}

// NewCreateMRolesTable creates a new CreateMRolesTable instance
// @wired:provide(group=migrations)
func NewCreateMRolesTable() migrations.Migration {
	return &CreateMRolesTable{}
}

func (m *CreateMRolesTable) Name() string { return "2025_12_17_102624_create_m_roles_table.go" }

func (m *CreateMRolesTable) Up(schema *migrations.Schema) error {
	return schema.Create("m_roles", func(table *migrations.Blueprint) {
		table.ID()
		table.Timestamps()
		table.SoftDeletes()
		table.Text("name")
		table.Text("description")
		table.Text("permissions")

		table.Index("name")
		// Add your columns here
	})
}

func (m *CreateMRolesTable) Down(schema *migrations.Schema) error {
	return schema.Drop("m_roles")
}
