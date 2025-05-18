package sqlite

import (
	"database/sql"
	"errors"
)

type migrator struct{}

func NewMigrator() *migrator {
	return nil
}

func (m *migrator) Setup(txn *sql.Tx, expectedMigrations []string) error {
	_, err := txn.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id TEXT PRIMARY KEY
		)
	`)
	return err
}

func (m *migrator) Migrate(txn *sql.Tx, schema string, content string) error {
	if exists, err := m.checkExists(txn, schema); err != nil {
		return err
	} else if exists {
		return nil
	}

	return m.migrate(txn, schema, content)
}

func (m *migrator) checkExists(txn *sql.Tx, schema string) (found bool, err error) {
	row := txn.QueryRow(`
		SELECT * FROM schema_migrations WHERE id=? LIMIT 1
	`, schema)

	err = row.Scan()
	found = errors.Is(err, sql.ErrNoRows)
	return
}

func (m *migrator) migrate(txn *sql.Tx, schema string, content string) error {
	if _, err := txn.Exec(content); err != nil {
		return err
	}

	if _, err := txn.Exec(`INSERT INTO schema_migrations(id) VALUES (?)`, schema); err != nil {
		return err
	}

	return nil
}
