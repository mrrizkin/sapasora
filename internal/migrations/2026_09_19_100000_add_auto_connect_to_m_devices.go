// Package migrations provides the migrations for the database
package migrations

import "sapasora/platform/database/migrations"

type AddAutoConnectToMDevices struct{}

// NewAddAutoConnectToMDevices creates the explicit provider-startup opt-in migration.
// @wired:provide(group=migrations)
func NewAddAutoConnectToMDevices() migrations.Migration {
	return &AddAutoConnectToMDevices{}
}

func (m *AddAutoConnectToMDevices) Name() string {
	return "2026_09_19_100000_add_auto_connect_to_m_devices.go"
}

func (m *AddAutoConnectToMDevices) Up(schema *migrations.Schema) error {
	// FALSE is deliberate: upgrading must not silently reconnect existing devices.
	return schema.Exec(`
		ALTER TABLE m_devices
		ADD COLUMN auto_connect BOOLEAN NOT NULL DEFAULT FALSE
	`)
}

func (m *AddAutoConnectToMDevices) Down(schema *migrations.Schema) error {
	return schema.Exec(`ALTER TABLE m_devices DROP COLUMN auto_connect`)
}
