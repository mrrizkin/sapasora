package provider

import (
	"sapasora/platform/database"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/mysql/v2"
	"github.com/gofiber/storage/postgres/v3"
	"github.com/gofiber/storage/sqlite3/v2"
)

type Database struct {
	config *database.Config
}

func NewDatabase(config *database.Config) *Database {
	return &Database{config: config}
}

func (d *Database) Setup() (fiber.Storage, error) {
	switch d.config.Protocol() {
	case "postgres":
		return createPostgresStorage(d.config)
	case "mysql":
		return createMysqlStorage(d.config)
	case "sqlite3":
		return createSQLiteStorage(d.config)
	default:
		return nil, fmt.Errorf("unknown database driver: %s", d.config.Driver)
	}
}

func createPostgresStorage(cfg *database.Config) (fiber.Storage, error) {
	config := postgres.Config{
		Host:       cfg.Host,
		Port:       cfg.Port,
		Database:   cfg.Name,
		Username:   cfg.Username,
		Password:   cfg.Password,
		Table:      "sessions",
		SSLMode:    cfg.SSLMode,
		Reset:      false,
		GCInterval: 10 * time.Second,
	}

	return postgres.New(config), nil
}

func createMysqlStorage(cfg *database.Config) (fiber.Storage, error) {
	config := mysql.Config{
		Host:       cfg.Host,
		Port:       cfg.Port,
		Database:   cfg.Name,
		Username:   cfg.Username,
		Password:   cfg.Password,
		Table:      "sessions",
		Reset:      false,
		GCInterval: 10 * time.Second,
	}

	return mysql.New(config), nil
}

func createSQLiteStorage(cfg *database.Config) (fiber.Storage, error) {
	config := sqlite3.Config{
		Database:   cfg.Host,
		Table:      "sessions",
		Reset:      false,
		GCInterval: 10 * time.Second,
	}

	return sqlite3.New(config), nil
}
