package database

import (
	"sapasora/platform/config"
	"fmt"
)

type Config struct {
	Driver   string
	Host     string
	Port     int
	Name     string
	Username string
	Password string
	SSLMode  string
}

func NewConfig(cfg config.Config) *Config {
	return &Config{
		Driver:   cfg.GetString("database.driver"),
		Host:     cfg.GetString("database.host"),
		Port:     cfg.GetInt("database.port"),
		Name:     cfg.GetString("database.name"),
		Username: cfg.GetString("database.username"),
		Password: cfg.GetString("database.password"),
		SSLMode:  cfg.GetString("database.ssl_mode"),
	}
}

func (c *Config) DSN() string {
	switch c.Protocol() {
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
			c.Host, c.Port, c.Username, c.Name, c.SSLMode, c.Password,
		)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.Username, c.Password, c.Host, c.Port, c.Name,
		)
	case "sqlite3":
		return fmt.Sprintf(
			"file:%s?cache=shared&mode=rwc&_foreign_keys=true&?_pragma=foreign_keys(1)",
			c.Host,
		)
	default:
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.Username, c.Password, c.Host, c.Port, c.Name,
		)
	}
}

func (c *Config) Protocol() string {
	switch c.Driver {
	case "postgres", "postgresql", "psql", "pgx", "pgsql":
		return "postgres"
	case "mysql", "mysqlx", "mssql", "mariadb":
		return "mysql"
	case "sqlite3", "sqlite", "file":
		return "sqlite3"
	default:
		return c.Driver
	}
}
