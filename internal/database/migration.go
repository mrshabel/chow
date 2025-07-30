package database

import (
	"database/sql"
	"io"
	"os"
)

const (
	migrationFilePath = "init.sql"
)

func migrate(db *sql.DB) error {
	f, err := os.Open(migrationFilePath)
	if err != nil {
		return err
	}

	contents, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	// run migration in a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(string(contents))
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
