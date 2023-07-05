package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(path string) (usf UserStorage) {
	var err error

	dsn := "file:" + path
	usf.db, err = sql.Open("sqlite", dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "opening db: %v\n", err)
		os.Exit(1)
	}

	_, err = usf.db.Exec("CREATE TABLE IF NOT EXISTS users(name string UNIQUE,password string);", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "creating table: %v\n", err)
		os.Exit(1)
	}
	return usf
}

func (usf UserStorage) Close() error {
	if err := usf.db.Close(); err != nil {
		return fmt.Errorf("closing db: %v", err)
	}
	return nil
}

func (usf UserStorage) UserExist(name string) (bool, error) {
	row := usf.db.QueryRow("SELECT name FROM users WHERE name=?;", name)

	var s string
	err := row.Scan(&s)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, sql.ErrNoRows):
		return false, nil
	default:
		return false, fmt.Errorf("querying user: %v", err)
	}
}

func (usf UserStorage) UserAdd(name string, password string) error {
	_, err := usf.db.Exec("INSERT INTO users (name, password) VALUES (?, ?);", name, password)
	if err != nil {
		return fmt.Errorf("checking is user exist: %v", err)
	}

	return nil
}

func (usf UserStorage) UserGet(name string) (string, error) {
	row := usf.db.QueryRow("SELECT password FROM users WHERE name=?;", name)

	var s string
	if err := row.Scan(&s); err != nil {
		return "", fmt.Errorf("getting user: %v", err)
	}

	return s, nil
}
