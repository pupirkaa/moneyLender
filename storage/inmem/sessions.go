package inmem

import (
	"time"

	ml "github.com/pupirkaa/moneyLender"
)

type SessionsStorage struct {
	sessions map[string]ml.SessionArgs
}

func NewSessionsStorage() (ss *SessionsStorage) {
	ss = &SessionsStorage{
		sessions: map[string]ml.SessionArgs{},
	}
	return ss
}

func (ss SessionsStorage) Close() error {
	return nil
}

func (ss SessionsStorage) SessionExist(session string) (error, bool) {
	if _, ok := ss.sessions[session]; !ok {
		return nil, false
	}
	return nil, true
}

func (ss *SessionsStorage) AddSession(session string, name string, creationDate time.Time) error {
	ss.sessions[session] = ml.SessionArgs{Name: name, CreationDate: creationDate}
	return nil
}

func (ss *SessionsStorage) DeleteSession(session string) error {
	delete(ss.sessions, session)
	return nil
}

func (ss *SessionsStorage) GetSessions() (map[string]ml.SessionArgs, error) {
	return ss.sessions, nil
}
