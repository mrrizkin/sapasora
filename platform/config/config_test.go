package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigManager(t *testing.T) {
	// Set up test environment variables
	os.Setenv("TEST_DB_HOST", "localhost")
	os.Setenv("TEST_DB_PORT", "5432")
	os.Setenv("TEST_DB_SSL", "true")
	os.Setenv("TEST_TIMEOUT", "30s")
	os.Setenv("TEST_RETRY_COUNT", "3")
	os.Setenv("TEST_RATE_LIMIT", "100.5")

	// Clean up environment variables after test
	defer func() {
		os.Unsetenv("TEST_DB_HOST")
		os.Unsetenv("TEST_DB_PORT")
		os.Unsetenv("TEST_DB_SSL")
		os.Unsetenv("TEST_TIMEOUT")
		os.Unsetenv("TEST_RETRY_COUNT")
		os.Unsetenv("TEST_RATE_LIMIT")
	}()

	config := NewConfigManager()

	t.Run("Set and Get", func(t *testing.T) {
		config.Set("app.name", "test-app")
		assert.Equal(t, "test-app", config.Get("app.name"))
	})

	t.Run("GetString", func(t *testing.T) {
		result := config.GetString("test.db.host")
		assert.Equal(t, "localhost", result)
		assert.Equal(t, "default", config.GetString("nonexistent", "default"))
	})

	t.Run("GetInt", func(t *testing.T) {
		result := config.GetInt("test.db.port")
		assert.Equal(t, 5432, result)
		assert.Equal(t, 999, config.GetInt("nonexistent", 999))
	})

	t.Run("GetBool", func(t *testing.T) {
		result := config.GetBool("test.db.ssl")
		assert.True(t, result)
		assert.False(t, config.GetBool("nonexistent", false))
	})

	t.Run("GetFloat64", func(t *testing.T) {
		result := config.GetFloat64("test.rate.limit")
		assert.Equal(t, 100.5, result)
		assert.Equal(t, 50.0, config.GetFloat64("nonexistent", 50.0))
	})

	t.Run("GetDuration", func(t *testing.T) {
		result := config.GetDuration("test.timeout")
		assert.Equal(t, 30*time.Second, result)
		assert.Equal(t, 10*time.Second, config.GetDuration("nonexistent", 10*time.Second))
	})

	t.Run("Has", func(t *testing.T) {
		assert.True(t, config.Has("test.db.host"))
		assert.False(t, config.Has("nonexistent"))
	})

	t.Run("All", func(t *testing.T) {
		all := config.All()
		assert.NotEmpty(t, all)
	})
}

func TestLoadStruct(t *testing.T) {
	os.Setenv("APP_NAME", "test-app")
	os.Setenv("APP_PORT", "8080")
	os.Setenv("APP_DEBUG", "true")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_SSL", "false")
	os.Setenv("CACHE_SIZE", "100")
	os.Setenv("API_TIMEOUT", "30s")

	defer func() {
		os.Unsetenv("APP_NAME")
		os.Unsetenv("APP_PORT")
		os.Unsetenv("APP_DEBUG")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_SSL")
		os.Unsetenv("CACHE_SIZE")
		os.Unsetenv("API_TIMEOUT")
	}()

	type DatabaseConfig struct {
		Host string `env:"DB_HOST"`
		Port int    `env:"DB_PORT"`
		SSL  bool   `env:"DB_SSL"`
	}

	type AppConfig struct {
		Name       string         `env:"APP_NAME"`
		Port       int            `env:"APP_PORT"`
		Debug      bool           `env:"APP_DEBUG"`
		Database   DatabaseConfig `env:"DATABASE"`
		CacheSize  uint           `env:"CACHE_SIZE"`
		APITimeout time.Duration  `env:"API_TIMEOUT"`
	}

	config := NewConfigManager()
	appConfig := AppConfig{}

	err := config.LoadStruct("app", &appConfig)
	require.NoError(t, err)

	assert.Equal(t, "test-app", appConfig.Name)
	assert.Equal(t, 8080, appConfig.Port)
	assert.True(t, appConfig.Debug)
	assert.Equal(t, "localhost", appConfig.Database.Host)
	assert.Equal(t, 5432, appConfig.Database.Port)
	assert.False(t, appConfig.Database.SSL)
	assert.Equal(t, uint(100), appConfig.CacheSize)
	assert.Equal(t, 30*time.Second, appConfig.APITimeout)
}

func TestLoadStructWithPointers(t *testing.T) {
	os.Setenv("PTR_NAME", "pointer-test")
	os.Setenv("PTR_COUNT", "42")

	defer func() {
		os.Unsetenv("PTR_NAME")
		os.Unsetenv("PTR_COUNT")
	}()

	type PointerConfig struct {
		Name  *string `env:"PTR_NAME"`
		Count *int    `env:"PTR_COUNT"`
	}

	config := NewConfigManager()
	ptrConfig := PointerConfig{}

	err := config.LoadStruct("ptr", &ptrConfig)
	require.NoError(t, err)

	require.NotNil(t, ptrConfig.Name)
	assert.Equal(t, "pointer-test", *ptrConfig.Name)

	require.NotNil(t, ptrConfig.Count)
	assert.Equal(t, 42, *ptrConfig.Count)
}

func TestLoadStructRequiredField(t *testing.T) {
	t.Setenv("REQUIRED_NAME", "")

	type RequiredConfig struct {
		Name string `env:"REQUIRED_NAME,required"`
	}

	config := NewConfigManager()
	reqConfig := RequiredConfig{}

	err := config.LoadStruct("req", &reqConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required config")
}

func TestLoadStructRequiredFieldRejectsEmptyValue(t *testing.T) {
	t.Setenv("REQUIRED_EMPTY_NAME", "   ")

	type RequiredConfig struct {
		Name string `env:"REQUIRED_EMPTY_NAME,required"`
	}

	config := NewConfigManager()
	err := config.LoadStruct("req", &RequiredConfig{})
	assert.ErrorContains(t, err, "required config")
}

func TestLoadStructWithDefault(t *testing.T) {
	type DefaultConfig struct {
		Name    string `env:"DEFAULT_NAME,default=myapp"`
		Count   int    `env:"DEFAULT_COUNT,default=10"`
		Enabled bool   `env:"DEFAULT_ENABLED,default=true"`
	}

	config := NewConfigManager()
	defConfig := DefaultConfig{}

	err := config.LoadStruct("default", &defConfig)
	require.NoError(t, err)

	assert.Equal(t, "myapp", defConfig.Name)
	assert.Equal(t, 10, defConfig.Count)
	assert.True(t, defConfig.Enabled)
}

func TestHelpers(t *testing.T) {
	os.Setenv("HELPER_NAME", "helper-test")

	defer func() {
		os.Unsetenv("HELPER_NAME")
	}()

	config := NewConfigManager()

	t.Run("MustLoad", func(t *testing.T) {
		type HelperConfig struct {
			Name string `env:"HELPER_NAME"`
		}

		helperConfig := HelperConfig{}
		result := MustLoad(config, "helper", &helperConfig)
		assert.Equal(t, "helper-test", result.Name)
	})

	t.Run("LoadOrDefault", func(t *testing.T) {
		type HelperConfig struct {
			Name string `env:"NONEXISTENT,default=default-value"`
		}

		helperConfig := HelperConfig{}
		defaultConfig := HelperConfig{Name: "default"}

		result := LoadOrDefault(config, "helper", &helperConfig, &defaultConfig)
		assert.Equal(t, "default-value", result.Name)
	})
}
