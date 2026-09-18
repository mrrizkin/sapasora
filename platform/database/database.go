package database

import (
	"sapasora/platform/config"
	"sapasora/platform/logger"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type Database struct {
	*gorm.DB

	log    *logger.Logger
	config *Config
}

func NewDatabase(config *Config, cfg config.Config, log *logger.Logger) *Database {
	var db *gorm.DB
	var err error
	var level gormLogger.LogLevel

	if !isDevelopment(cfg) {
		level = gormLogger.Silent
	} else {
		switch log.GetLevel() {
		case "trace", "debug", "info":
			level = gormLogger.Info
		case "warn":
			level = gormLogger.Warn
		case "error", "fatal":
			level = gormLogger.Error
		case "none", "disable":
			level = gormLogger.Silent
		default:
			level = gormLogger.Silent
		}
	}

	gormConfig := &gorm.Config{
		Logger: log.GetGormLogger(level),
	}

	switch config.Protocol() {
	case "postgres":
		db, err = gorm.Open(postgres.Open(config.DSN()), gormConfig)
		if err != nil {
			log.Fatal("Failed to connect to database", "error", err)
		}
	case "sqlite3":
		db, err = gorm.Open(sqlite.Open(config.DSN()), gormConfig)
		if err != nil {
			log.Fatal("Failed to connect to database", "error", err)
		}
	case "mysql":
		db, err = gorm.Open(mysql.Open(config.DSN()), gormConfig)
		if err != nil {
			log.Fatal("Failed to connect to database", "error", err)
		}
	}

	return &Database{
		DB:     db,
		log:    log,
		config: config,
	}
}

func (d *Database) RunMigrations(models ...any) {
	if err := d.AutoMigrate(models...); err != nil {
		d.log.Fatal("Failed to run migrations", "error", err)
	}

	d.log.Info("Migrations completed successfully")
}

func isDevelopment(cfg config.Config) bool {
	env := cfg.GetString("app.env", "development")
	return env == "development" || env == "dev"
}
