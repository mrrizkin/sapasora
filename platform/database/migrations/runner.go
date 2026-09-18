package migrations

import (
	"sapasora/platform/database"
	"sapasora/platform/support/arr"
	"sapasora/platform/support/console"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"time"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

type MigrationRecord struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Migration string    `gorm:"unique;not null"`
	Batch     int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}

func (m *MigrationRecord) TableName() string {
	return "migrations"
}

type MigrationEntry struct {
	Name      string
	Migration Migration
}

type MigrationRunner struct {
	db         *database.Database
	migrations []MigrationEntry

	forbiddenTables []string // we can't modify this tables
}

type MigrationRunnerIn struct {
	fx.In

	DB         *database.Database
	Migrations []Migration `group:"migrations"`
}

func NewMigrationRunner(in MigrationRunnerIn) *MigrationRunner {
	return &MigrationRunner{
		db: in.DB,
		migrations: arr.Map(in.Migrations, func(m Migration) MigrationEntry {
			return MigrationEntry{
				Name:      m.Name(),
				Migration: m,
			}
		}),

		forbiddenTables: []string{
			"migrations",
			"sqlite_sequence",
		},
	}
}

func (mr *MigrationRunner) Register(name string, migration Migration) {
	mr.migrations = append(mr.migrations, MigrationEntry{
		Name:      name,
		Migration: migration,
	})
}

func (mr *MigrationRunner) Run() error {
	if err := mr.ensureMigrationsTable(); err != nil {
		return fmt.Errorf("failed to ensure migrations table: %w", err)
	}

	sort.Slice(mr.migrations, func(i, j int) bool {
		return mr.migrations[i].Name < mr.migrations[j].Name
	})

	lastBatch, err := mr.getLastBatchNumber()
	if err != nil {
		return fmt.Errorf("failed to get last batch number: %w", err)
	}

	nextBatch := lastBatch + 1

	executedMigrations, err := mr.getExecutedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get executed migrations: %w", err)
	}

	executedSet := make(map[string]bool)
	for _, name := range executedMigrations {
		executedSet[name] = true
	}

	for _, migrationEntry := range mr.migrations {
		if _, exists := executedSet[migrationEntry.Name]; exists {
			console.Info("Migration %s already executed, skipping...", migrationEntry.Name)
			continue
		}

		console.Info("Running migration: %s", migrationEntry.Name)

		schema := NewSchema(mr.db)

		if err := migrationEntry.Migration.Up(schema); err != nil {
			return fmt.Errorf("migration %s failed: %w", migrationEntry.Name, err)
		}

		record := MigrationRecord{
			Migration: migrationEntry.Name,
			Batch:     nextBatch,
			CreatedAt: time.Now(),
		}

		if err := mr.db.Create(&record).Error; err != nil {
			return fmt.Errorf("failed to record migration %s: %w", migrationEntry.Name, err)
		}

		console.Info("Migration %s completed successfully", migrationEntry.Name)
	}

	console.Info("All migrations completed. Batch: %d", nextBatch)
	return nil
}

func (mr *MigrationRunner) Rollback() error {
	lastBatch, err := mr.getLastBatchNumber()
	if err != nil {
		return fmt.Errorf("failed to get last batch number: %w", err)
	}

	if lastBatch == 0 {
		console.Info("No migrations to rollback")
		return nil
	}

	lastBatchMigrations, err := mr.getMigrationsByBatch(lastBatch)
	if err != nil {
		return fmt.Errorf("failed to get migrations for batch %d: %w", lastBatch, err)
	}

	sort.Slice(lastBatchMigrations, func(i, j int) bool {
		return lastBatchMigrations[i].Migration > lastBatchMigrations[j].Migration
	})

	migrationMap := make(map[string]Migration)
	for _, entry := range mr.migrations {
		migrationMap[entry.Name] = entry.Migration
	}

	for _, record := range lastBatchMigrations {
		migration, exists := migrationMap[record.Migration]
		if !exists {
			return fmt.Errorf("migration %s not found in registered migrations", record.Migration)
		}

		console.Info("Rolling back migration: %s", record.Migration)

		schema := NewSchema(mr.db)

		if err := migration.Down(schema); err != nil {
			return fmt.Errorf("rollback for migration %s failed: %w", record.Migration, err)
		}

		if err := mr.db.Where("migration = ? AND batch = ?", record.Migration, lastBatch).Delete(&MigrationRecord{}).Error; err != nil {
			return fmt.Errorf("failed to remove migration record %s: %w", record.Migration, err)
		}

		console.Info("Migration %s rolled back successfully", record.Migration)
	}

	console.Info("Batch %d rolled back successfully", lastBatch)
	return nil
}

func (mr *MigrationRunner) Reset() error {
	allBatches, err := mr.getAllBatches()
	if err != nil {
		return err
	}

	sort.Sort(sort.Reverse(sort.IntSlice(allBatches)))

	for _, batch := range allBatches {
		if err := mr.rollbackBatch(batch); err != nil {
			return err
		}
	}

	console.Info("All migrations have been reset")
	return nil
}

func (mr *MigrationRunner) rollbackBatch(batch int) error {
	lastBatchMigrations, err := mr.getMigrationsByBatch(batch)
	if err != nil {
		return fmt.Errorf("failed to get migrations for batch %d: %w", batch, err)
	}

	sort.Slice(lastBatchMigrations, func(i, j int) bool {
		return lastBatchMigrations[i].Migration > lastBatchMigrations[j].Migration
	})

	migrationMap := make(map[string]Migration)
	for _, entry := range mr.migrations {
		migrationMap[entry.Name] = entry.Migration
	}

	for _, record := range lastBatchMigrations {
		migration, exists := migrationMap[record.Migration]
		if !exists {
			return fmt.Errorf("migration %s not found in registered migrations", record.Migration)
		}

		console.Info("Rolling back migration: %s (batch %d)", record.Migration, batch)

		schema := NewSchema(mr.db)

		if err := migration.Down(schema); err != nil {
			return fmt.Errorf("rollback for migration %s failed: %w", record.Migration, err)
		}

		if err := mr.db.Where("migration = ? AND batch = ?", record.Migration, batch).Delete(&MigrationRecord{}).Error; err != nil {
			return fmt.Errorf("failed to remove migration record %s: %w", record.Migration, err)
		}

		console.Info("Migration %s (batch %d) rolled back successfully", record.Migration, batch)
	}

	console.Info("Batch %d rolled back successfully", batch)
	return nil
}

func (mr *MigrationRunner) Fresh() error {
	console.Info("Dropping all tables...")

	tables, err := mr.db.Migrator().GetTables()
	tables = arr.Filter(tables, func(s string) bool {
		return !slices.Contains(mr.forbiddenTables, s)
	})
	if err != nil {
		return fmt.Errorf("failed to get table names: %w", err)
	}

	for _, table := range tables {
		if err := mr.db.Migrator().DropTable(table); err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}

		console.Info("Dropped table: %s", table)
	}

	if err := mr.db.Where("1 = 1").Delete(&MigrationRecord{}).Error; err != nil {
		return fmt.Errorf("failed to clear migration records: %w", err)
	}

	console.Info("Cleared migration records")

	console.Info("Running all migrations again...")
	return mr.Run()
}

func (mr *MigrationRunner) Status() error {
	if err := mr.ensureMigrationsTable(); err != nil {
		return fmt.Errorf("failed to ensure migrations table: %w", err)
	}

	executedMigrations, err := mr.getExecutedMigrationsWithBatch()
	if err != nil {
		return fmt.Errorf("failed to get executed migrations: %w", err)
	}

	executedSet := make(map[string]int) // name -> batch
	for _, record := range executedMigrations {
		executedSet[record.Migration] = record.Batch
	}

	console.Line("Migration Status")

	headers := []string{"Migration", "Batch", "Status"}
	data := make([][]string, 0)

	if len(executedMigrations) > 0 {
		for _, record := range executedMigrations {
			data = append(data, []string{record.Migration, strconv.Itoa(record.Batch), "Executed"})
		}
	}

	pendingCount := 0
	for _, migrationEntry := range mr.migrations {
		if _, exists := executedSet[migrationEntry.Name]; !exists {
			data = append(data, []string{migrationEntry.Name, "", "Pending"})
		}
	}

	if len(data) > 0 {
		console.Table(headers, data)
	} else {
		console.Line("No migrations")
	}

	console.Info("Total: %d executed, %d pending", len(executedMigrations), pendingCount)
	return nil
}

func (mr *MigrationRunner) ensureMigrationsTable() error {
	if !mr.db.Migrator().HasTable(&MigrationRecord{}) {
		if err := mr.db.Migrator().CreateTable(&MigrationRecord{}); err != nil {
			return fmt.Errorf("failed to create migrations table: %w", err)
		}
	}
	return nil
}

func (mr *MigrationRunner) getLastBatchNumber() (int, error) {
	var lastRecord MigrationRecord
	if err := mr.db.Order("batch DESC").First(&lastRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil // No migrations executed yet
		}
		return 0, err
	}

	return int(lastRecord.Batch), nil
}

func (mr *MigrationRunner) getExecutedMigrations() ([]string, error) {
	var records []MigrationRecord
	if err := mr.db.Find(&records).Error; err != nil {
		return nil, err
	}

	var migrationNames []string
	for _, record := range records {
		migrationNames = append(migrationNames, record.Migration)
	}

	return migrationNames, nil
}

func (mr *MigrationRunner) getExecutedMigrationsWithBatch() ([]MigrationRecord, error) {
	var records []MigrationRecord
	if err := mr.db.Order("batch ASC, migration ASC").Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

func (mr *MigrationRunner) getMigrationsByBatch(batch int) ([]MigrationRecord, error) {
	var records []MigrationRecord
	if err := mr.db.Where("batch = ?", batch).Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

func (mr *MigrationRunner) getAllBatches() ([]int, error) {
	var batches []int
	err := mr.db.Raw("SELECT DISTINCT batch FROM migrations ORDER BY batch ASC").
		Scan(&batches).
		Error
	if err != nil {
		return nil, err
	}
	return batches, nil
}
