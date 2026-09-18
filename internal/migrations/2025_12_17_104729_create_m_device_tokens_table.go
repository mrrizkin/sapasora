package migrations

import "sapasora/platform/database/migrations"

type CreateMDeviceTokensTable struct{}

// NewCreateMDeviceTokensTable creates a new CreateMDeviceTokensTable instance
// @wired:provide(group=migrations)
func NewCreateMDeviceTokensTable() migrations.Migration {
	return &CreateMDeviceTokensTable{}
}

func (m *CreateMDeviceTokensTable) Name() string {
	return "2025_12_17_104729_create_m_device_tokens_table.go"
}

func (m *CreateMDeviceTokensTable) Up(schema *migrations.Schema) error {
	return schema.Create("m_device_tokens", func(table *migrations.Blueprint) {
		table.ID()
		table.Timestamps()
		table.SoftDeletes()
		table.Text("token")
		table.Text("status")
		table.Timestamp("expired_at").Nullable()

		table.BigInteger("device_id").Nullable().Unsigned()
		table.Foreign("device_id").
			On("m_devices").
			References("id").
			Build()

		table.BigInteger("user_id").Nullable().Unsigned()
		table.Foreign("user_id").
			On("m_users").
			References("id").
			Build()

		table.Index("token")
		table.Index("status")
		table.Index("device_id")
		table.Index("user_id")
	})
}

func (m *CreateMDeviceTokensTable) Down(schema *migrations.Schema) error {
	return schema.Drop("m_device_tokens")
}
