package migrations

import "sapasora/platform/database/migrations"

type CreateMUsersTable struct{}

// NewCreateMUsersTable creates a new CreateMUsersTable instance
// @wired:provide(group=migrations)
func NewCreateMUsersTable() migrations.Migration {
	return &CreateMUsersTable{}
}

func (m *CreateMUsersTable) Name() string { return "2025_12_17_102646_create_m_users_table.go" }

func (m *CreateMUsersTable) Up(schema *migrations.Schema) error {
	return schema.Create("m_users", func(table *migrations.Blueprint) {
		table.ID()
		table.Timestamps()
		table.SoftDeletes()
		table.Text("name")
		table.Text("username").Unique()
		table.Text("hashed_password")

		table.BigInteger("role_id").Nullable().Unsigned()
		table.Foreign("role_id").
			On("m_roles").
			References("id").
			Build()

		table.Index("username")
		table.Index("role_id")
		// Add your columns here
	})
}

func (m *CreateMUsersTable) Down(schema *migrations.Schema) error {
	return schema.Drop("m_users")
}
