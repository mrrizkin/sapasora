package provider

import (
	"sapasora/platform/config"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/memory/v2"
)

type Memory struct {
	config config.Config
}

func NewMemory(config config.Config) *Memory {
	return &Memory{config: config}
}

func (m *Memory) Setup() (fiber.Storage, error) {
	driver := m.config.GetString("session.driver")
	switch driver {
	case "memory":
		return createMemoryStorage()
	case "redis", "valkey":
		return nil, fmt.Errorf("driver %s is not yet supported", driver)
	default:
		return nil, fmt.Errorf("unknown database driver: %s", driver)
	}
}

func createMemoryStorage() (fiber.Storage, error) {
	config := memory.Config{
		GCInterval: 10 * time.Second,
	}

	return memory.New(config), nil
}
