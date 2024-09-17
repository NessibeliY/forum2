package store

import (
	"database/sql"
	"fmt"
	"os"

	"forum/pkg/logger"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(l *logger.Logger, dbDriver, dbPath, migrationPath string) (*sql.DB, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open(dbDriver, dbPath)
	if err != nil {
		fmt.Println("OPEN ERROR")
		return nil, fmt.Errorf("%w %s", err, op)
	}

	if err = db.Ping(); err != nil {
		fmt.Println("PING ERROR")
		return nil, fmt.Errorf("%w %s", err, op)
	}

	stmt, err := os.ReadFile(migrationPath)
	if err != nil {
		fmt.Println("READ FILE ERROR")
		return nil, fmt.Errorf("%w %s", err, op)
	}

	_, err = db.Exec(string(stmt))
	if err != nil {
		return nil, fmt.Errorf("%w %s", err, op)
	}

	l.Info("connected to db")
	return db, nil
}
