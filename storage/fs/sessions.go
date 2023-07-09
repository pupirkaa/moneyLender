package fs

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type SessionsStorage struct {
	sessions   map[string]bool
	path       string
	exitCh     chan bool
	syncDoneCh chan bool
}

func NewSessionsStorage(path string) (ss SessionsStorage) {
	ss = SessionsStorage{
		sessions:   parseSesions(readFile(path)),
		path:       path,
		exitCh:     make(chan bool),
		syncDoneCh: make(chan bool),
	}
	go func() {
		t := time.NewTicker(1 * time.Minute)
		for {
			select {
			case <-ss.exitCh:
				ss.saveSessionsToFile()
				close(ss.syncDoneCh)
				return
			case <-t.C:
				ss.saveSessionsToFile()
			}
		}
	}()
	return ss
}

func (ss SessionsStorage) Close() error {
	close(ss.exitCh)
	<-ss.syncDoneCh
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

func (ss *SessionsStorage) saveSessionsToFile() {
	f, err := os.Create(ss.path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file: %v", err)
		os.Exit(1)
	}
	defer f.Close()
	for t := range ss.sessions {
		_, err = f.WriteString(fmt.Sprintf("%s\n", t))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func parseSesions(s []string) map[string]bool {
	sessions := make(map[string]bool)

	for i := range s {
		k := strings.Split(s[i], ",")
		sessions[k[0]] = true
	}
	return sessions
}
