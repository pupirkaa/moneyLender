package inmem

type SessionsStorage struct {
	sessions map[string]bool
}

func NewSessionsStorage() (ss SessionsStorage) {
	ss.sessions = map[string]bool{}
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

func (ss *SessionsStorage) AddSession(session string) error {
	ss.sessions[session] = true
	return nil
}

func (ss *SessionsStorage) DeleteSession(session string) error {
	delete(ss.sessions, session)
	return nil
}
