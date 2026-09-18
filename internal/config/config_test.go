package config

import (
	"testing"
	"time"

	platformconfig "sapasora/platform/config"

	"github.com/stretchr/testify/require"
)

func setRequiredEnvironment(t *testing.T) {
	t.Helper()
	t.Setenv("APP_NAME", "test")
	t.Setenv("APP_ENV", "production")
	t.Setenv("APP_URL", "https://example.test")
	t.Setenv("APP_PORT", "3000")
	t.Setenv("DB_DRIVER", "sqlite")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_NAME", "test")
	t.Setenv("DB_USERNAME", "test")
	t.Setenv("TELEGRAM_API_ID", "123")
	t.Setenv("TELEGRAM_API_HASH", "hash")
}

func TestNewAppConfigPropagatesLoadErrors(t *testing.T) {
	setRequiredEnvironment(t)
	t.Setenv("DB_PORT", "not-a-port")

	_, err := NewAppConfig(platformconfig.NewConfigManager())
	require.Error(t, err)
	require.ErrorContains(t, err, "load database config")
	require.ErrorContains(t, err, "invalid int value")
}

func TestNewAppConfigRejectsDurationWithoutUnit(t *testing.T) {
	setRequiredEnvironment(t)
	t.Setenv("SERVER_READ_TIMEOUT", "30")

	_, err := NewAppConfig(platformconfig.NewConfigManager())
	require.Error(t, err)
	require.ErrorContains(t, err, "load server config")
	require.ErrorContains(t, err, "invalid duration value")
}

func TestNewAppConfigLoadsServerConfig(t *testing.T) {
	setRequiredEnvironment(t)
	t.Setenv("SERVER_READ_TIMEOUT", "2s")
	t.Setenv("SERVER_WRITE_TIMEOUT", "3s")
	t.Setenv("SERVER_IDLE_TIMEOUT", "4s")
	t.Setenv("SERVER_BODY_LIMIT", "1024")

	cfg, err := NewAppConfig(platformconfig.NewConfigManager())
	require.NoError(t, err)
	require.Equal(t, 2*time.Second, cfg.GetDuration("server.read_timeout"))
	require.Equal(t, 3*time.Second, cfg.GetDuration("server.write_timeout"))
	require.Equal(t, 4*time.Second, cfg.GetDuration("server.idle_timeout"))
	require.Equal(t, 1024, cfg.GetInt("server.body_limit"))
}
