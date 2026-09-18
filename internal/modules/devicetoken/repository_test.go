package devicetoken

import (
	"context"
	"testing"

	"sapasora/platform/database"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestListDeviceTokensByDeviceIDScopesToDeviceAndSkipsDeleted(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.Exec(`
		CREATE TABLE m_device_tokens (
			id INTEGER PRIMARY KEY,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME NULL,
			public_id TEXT,
			token TEXT,
			status TEXT,
			expired_at DATETIME NULL,
			device_id INTEGER,
			user_id INTEGER
		)
	`).Error)
	require.NoError(t, db.Exec(`
		INSERT INTO m_device_tokens (id, public_id, token, status, device_id)
		VALUES (1, 'token-1', 'secret-1', 'active', 7),
		       (2, 'token-2', 'secret-2', 'active', 8),
		       (3, 'token-3', 'secret-3', 'active', 7)
	`).Error)
	require.NoError(t, db.Exec("UPDATE m_device_tokens SET deleted_at = CURRENT_TIMESTAMP WHERE id = 3").Error)

	repository := &DeviceTokenRepositoryImpl{db: &database.Database{DB: db}}
	result, err := repository.ListDeviceTokensByDeviceID(context.Background(), 7)

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, uint(1), result[0].ID)
}
