package migrations

import (
	"testing"

	platformdatabase "sapasora/platform/database"
	dbschema "sapasora/platform/database/migrations"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAddAutoConnectDefaultsToOptOutAndRollsBack(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.Exec(`CREATE TABLE m_devices (id INTEGER PRIMARY KEY, name TEXT)`).Error)

	schema := dbschema.NewSchema(&platformdatabase.Database{DB: db})
	migration := &AddAutoConnectToMDevices{}
	require.NoError(t, migration.Up(schema))
	require.True(t, schema.HasColumn("m_devices", "auto_connect"))

	require.NoError(t, db.Exec(`INSERT INTO m_devices (id, name) VALUES (1, 'existing')`).Error)
	var autoConnect bool
	require.NoError(t, db.Raw(`SELECT auto_connect FROM m_devices WHERE id = 1`).Scan(&autoConnect).Error)
	require.False(t, autoConnect)

	require.NoError(t, migration.Down(schema))
	require.False(t, schema.HasColumn("m_devices", "auto_connect"))
}
