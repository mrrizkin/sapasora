package migrations

import "sapasora/platform/database/migrations"

type CreateMDevicesTable struct{}

// NewCreateMDevicesTable creates a new CreateMDevicesTable instance
// @wired:provide(group=migrations)
func NewCreateMDevicesTable() migrations.Migration {
	return &CreateMDevicesTable{}
}

func (m *CreateMDevicesTable) Name() string { return "2025_12_17_102702_create_m_devices_table.go" }

func (m *CreateMDevicesTable) Up(schema *migrations.Schema) error {
	return schema.Create("m_devices", func(table *migrations.Blueprint) {
		table.ID()
		table.Timestamps()
		table.SoftDeletes()
		table.Text("name")
		table.Text("type")
		table.Text("qr_code").Nullable()
		table.Text("webhook").Nullable()
		table.Text("jid").Nullable()
		table.Text("status")
		table.Timestamp("expired_at").Nullable()
		table.Text("events").Nullable()
		table.Text("permissions")

		table.BigInteger("user_id").Nullable().Unsigned()
		table.Foreign("user_id").
			On("m_users").
			References("id").
			Build()

		table.Index("name")
		table.Index("type")
		table.Index("jid")
		table.Index("status")
		table.Index("expired_at")
		table.Index("user_id")
	})
}

func (m *CreateMDevicesTable) Down(schema *migrations.Schema) error {
	return schema.Drop("m_devices")
}
