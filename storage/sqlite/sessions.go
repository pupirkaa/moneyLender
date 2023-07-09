package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
)

type SessionsStorage struct {
	db *sql.DB
}

func NewSessionsStorage(path string) (ss SessionsStorage) {
	dsn := "file:" + path
	d, err := sql.Open("sqlite", dsn)
	ss = SessionsStorage{
		db: d,
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "opening db: %v\n", err)
		os.Exit(1)
	}

	_, err = ss.db.Exec("CREATE TABLE IF NOT EXISTS sessions(session string);", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "creating transactions table: %v\n", err)
		os.Exit(1)
	}

	return ss
}

func (ss SessionsStorage) Close() error {
	if err := ss.db.Close(); err != nil {
		return fmt.Errorf("closing db: %v", err)
	}
	return nil
}

func (ss SessionsStorage) SessionExist(session string) (error, bool) {
	row := ss.db.QueryRow("SELECT session FROM sessions WHERE session=?;", session)

	var s string
	err := row.Scan(&s)
	switch {
	case err == nil:
		return nil, true
	case errors.Is(err, sql.ErrNoRows):
		return nil, false
	default:
		return fmt.Errorf("querying session: %v", err), false
	}
}

func (ss *SessionsStorage) AddSession(session string) error {
	_, err := ss.db.Exec("INSERT INTO sessions (session) VALUES (?);", session)
	if err != nil {
		return fmt.Errorf("checking is session exist: %v", err)
	}

	return nil
}

func (ss *SessionsStorage) DeleteSession(session string) error {
	_, err := ss.db.Exec("DELETE FROM sessions WHERE session=?;", session)
	if err != nil {
		return fmt.Errorf("deleting session: %v", err)
	}

	return nil
}
