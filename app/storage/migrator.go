package storage

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"io"
	"io/fs"

	"github.com/fiffu/mailtl/app/storage/sqlite"
	"github.com/fiffu/mailtl/lib"
)

//go:embed migrations/.gitkeep
//go:embed migrations/*
var migrationsFS embed.FS
var migrator Migrator = sqlite.NewMigrator()
var errNoMigrations = errors.New("no migrations defined")

type Migrator interface {
	Setup(txn *sql.Tx, expectedMigrations []string) error
	Migrate(txn *sql.Tx, schema string, content string) error
}

type migration struct {
	filepath string
	content  string
}

// String implements interface fmt.Stringer
func (m migration) String() string { return m.filepath }

func Migrate(ctx context.Context, storage Storage) error {
	migrations, err := findMigrations()
	if err != nil {
		return err
	}

	return WithTransaction(ctx, storage, func(txn *sql.Tx) error {
		return runMigrations(txn, migrations)
	})
}

func findMigrations() ([]migration, error) {
	migrationFiles, err := fs.Glob(migrationsFS, "**/*.sql")
	if err != nil {
		return nil, err
	}
	if len(migrationFiles) == 0 {
		return nil, errNoMigrations
	}

	migrationContent := make([]migration, len(migrationFiles))
	for i, filepath := range migrationFiles {
		file, err := migrationsFS.Open(filepath)
		if err != nil {
			return nil, err
		}

		content, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}

		migrationContent[i].filepath = filepath
		migrationContent[i].content = string(content)
	}

	return migrationContent, nil
}

func runMigrations(txn *sql.Tx, migrations []migration) error {
	if err := migrator.Setup(txn, lib.StringsOf(migrations)); err != nil {
		return err
	}

	for _, migration := range migrations {
		if err := migrator.Migrate(txn, migration.filepath, migration.content); err != nil {
			return err
		}
	}
	return nil
}
