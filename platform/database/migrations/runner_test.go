package migrations

import (
	"errors"
	"testing"

	"sapasora/platform/database"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type testMigration struct {
	name string
	up   func(*Schema) error
	down func(*Schema) error
}

func (m testMigration) Name() string            { return m.name }
func (m testMigration) Up(schema *Schema) error { return m.up(schema) }
func (m testMigration) Down(schema *Schema) error {
	if m.down != nil {
		return m.down(schema)
	}
	return nil
}

func TestRunUsesTransactionPerMigration(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	runner := NewMigrationRunner(MigrationRunnerIn{
		DB: &database.Database{DB: db},
		Migrations: []Migration{
			testMigration{name: "001_success", up: func(schema *Schema) error {
				return schema.Exec("CREATE TABLE committed_migration (id INTEGER PRIMARY KEY)")
			}},
			testMigration{name: "002_failure", up: func(schema *Schema) error {
				if err := schema.Exec("CREATE TABLE rolled_back_migration (id INTEGER PRIMARY KEY)"); err != nil {
					return err
				}
				return errors.New("intentional migration failure")
			}},
		},
	})

	err = runner.Run()
	require.Error(t, err)
	require.True(t, db.Migrator().HasTable("committed_migration"))
	require.False(t, db.Migrator().HasTable("rolled_back_migration"))

	var records []MigrationRecord
	require.NoError(t, db.Find(&records).Error)
	require.Len(t, records, 1)
	require.Equal(t, "001_success", records[0].Migration)
}
