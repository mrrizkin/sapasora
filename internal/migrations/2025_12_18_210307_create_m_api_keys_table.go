// Package migrations provides the migrations for the database
package migrations

import "sapasora/platform/database/migrations"

type CreateMApiKeysTable struct{}

// NewCreateMApiKeysTable creates a new CreateMApiKeysTable instance
// @wired:provide(group=migrations)
func NewCreateMApiKeysTable() migrations.Migration {
	return &CreateMApiKeysTable{}
}

func (m *CreateMApiKeysTable) Name() string { return "2025_12_18_210307_create_m_api_keys_table.go" }

func (m *CreateMApiKeysTable) Up(schema *migrations.Schema) error {
	return schema.Create("m_api_keys", func(table *migrations.Blueprint) {
		table.ID()
		table.Timestamps()
		table.SoftDeletes()
		table.Text("name")
		table.Text("key").Unique()
		table.Text("status")
		table.Text("permissions")

		table.BigInteger("user_id").Nullable().Unsigned()
		table.Foreign("user_id").
			On("m_users").
			References("id").
			Build()

		table.Index("key")
		table.Index("status")
		table.Index("user_id")
	})
}

func (m *CreateMApiKeysTable) Down(schema *migrations.Schema) error {
	return schema.Drop("m_api_keys")
}
