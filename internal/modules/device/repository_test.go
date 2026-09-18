package device

import (
	"context"
	"testing"
	"time"

	"sapasora/platform/database"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newDeviceTokenRepositoryTest(t *testing.T) (*DeviceRepositoryImpl, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	require.NoError(t, db.Exec(`
		CREATE TABLE m_devices (
			id INTEGER PRIMARY KEY,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME NULL,
			public_id TEXT,
			name TEXT,
			type TEXT,
			qr_code TEXT,
			webhook TEXT,
			jid TEXT,
			status TEXT,
			expired_at DATETIME NULL,
			events TEXT,
			permissions TEXT,
			user_id INTEGER
		)
	`).Error)
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

	return &DeviceRepositoryImpl{db: &database.Database{DB: db}}, db
}

func insertDeviceTokenFixture(
	t *testing.T,
	db *gorm.DB,
	deviceStatus, tokenStatus string,
	tokenExpiredAt, tokenDeletedAt any,
	deviceDeletedAt any,
	tokenUserID, deviceUserID uint,
) {
	t.Helper()

	require.NoError(t, db.Exec(`
		INSERT INTO m_devices (id, public_id, name, type, status, user_id, deleted_at)
		VALUES (1, 'device-1', 'test device', 'whatsapp', ?, ?, ?)
	`, deviceStatus, deviceUserID, deviceDeletedAt).Error)
	require.NoError(t, db.Exec(`
		INSERT INTO m_device_tokens (id, public_id, token, status, expired_at, deleted_at, device_id, user_id)
		VALUES (1, 'token-1', 'sk-dat-test', ?, ?, ?, 1, ?)
	`, tokenStatus, tokenExpiredAt, tokenDeletedAt, tokenUserID).Error)
}

func TestGetDeviceByTokenAcceptsOnlyActiveUnexpiredOwnedDevice(t *testing.T) {
	tests := []struct {
		name            string
		deviceStatus    string
		tokenStatus     string
		tokenExpiredAt  any
		tokenDeletedAt  any
		deviceDeletedAt any
		tokenUserID     uint
		deviceUserID    uint
		wantDevice      bool
	}{
		{
			name:         "active token with no expiry",
			deviceStatus: DeviceStatusActive.String(),
			tokenStatus:  "active",
			tokenUserID:  7,
			deviceUserID: 7,
			wantDevice:   true,
		},
		{
			name:           "active token before expiry",
			deviceStatus:   DeviceStatusActive.String(),
			tokenStatus:    "active",
			tokenExpiredAt: time.Now().Add(time.Minute),
			tokenUserID:    7,
			deviceUserID:   7,
			wantDevice:     true,
		},
		{
			name:         "revoked token",
			deviceStatus: DeviceStatusActive.String(),
			tokenStatus:  "inactive",
			tokenUserID:  7,
			deviceUserID: 7,
			wantDevice:   false,
		},
		{
			name:           "expired token",
			deviceStatus:   DeviceStatusActive.String(),
			tokenStatus:    "active",
			tokenExpiredAt: time.Now().Add(-time.Minute),
			tokenUserID:    7,
			deviceUserID:   7,
			wantDevice:     false,
		},
		{
			name:           "soft deleted token",
			deviceStatus:   DeviceStatusActive.String(),
			tokenStatus:    "active",
			tokenDeletedAt: time.Now(),
			tokenUserID:    7,
			deviceUserID:   7,
			wantDevice:     false,
		},
		{
			name:            "soft deleted device",
			deviceStatus:    DeviceStatusActive.String(),
			tokenStatus:     "active",
			deviceDeletedAt: time.Now(),
			tokenUserID:     7,
			deviceUserID:    7,
			wantDevice:      false,
		},
		{
			name:         "inactive device",
			deviceStatus: "inactive",
			tokenStatus:  "active",
			tokenUserID:  7,
			deviceUserID: 7,
			wantDevice:   false,
		},
		{
			name:         "token owner differs from device owner",
			deviceStatus: DeviceStatusActive.String(),
			tokenStatus:  "active",
			tokenUserID:  8,
			deviceUserID: 7,
			wantDevice:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository, db := newDeviceTokenRepositoryTest(t)
			insertDeviceTokenFixture(
				t,
				db,
				test.deviceStatus,
				test.tokenStatus,
				test.tokenExpiredAt,
				test.tokenDeletedAt,
				test.deviceDeletedAt,
				test.tokenUserID,
				test.deviceUserID,
			)

			result, err := repository.GetDeviceByToken(context.Background(), "sk-dat-test")
			if test.wantDevice {
				require.NoError(t, err)
				require.Equal(t, uint(1), result.ID)
				return
			}

			require.ErrorIs(t, err, gorm.ErrRecordNotFound)
		})
	}
}

func TestGetDeviceByTokenForUserEnforcesCallerOwnerScope(t *testing.T) {
	repository, db := newDeviceTokenRepositoryTest(t)
	insertDeviceTokenFixture(
		t,
		db,
		DeviceStatusActive.String(),
		"active",
		nil,
		nil,
		nil,
		7,
		7,
	)

	result, err := repository.GetDeviceByTokenForUser(context.Background(), "sk-dat-test", 7)
	require.NoError(t, err)
	require.Equal(t, uint(1), result.ID)

	_, err = repository.GetDeviceByTokenForUser(context.Background(), "sk-dat-test", 8)
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
