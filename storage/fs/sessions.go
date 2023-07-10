package fs

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	ml "github.com/pupirkaa/moneyLender"
)

type SessionsStorage struct {
	sessions   map[string]ml.SessionArgs
	path       string
	exitCh     chan bool
	syncDoneCh chan bool
}

func NewSessionsStorage(path string) (ss *SessionsStorage) {
	ss = &SessionsStorage{
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

func (ss *SessionsStorage) saveSessionsToFile() {
	f, err := os.Create(ss.path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file: %v", err)
		os.Exit(1)
	}
	defer f.Close()
	for t := range ss.sessions {
		_, err = f.WriteString(fmt.Sprintf("%s,%s,%d\n", t, ss.sessions[t].Name, ss.sessions[t].CreationDate.Unix()))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func parseSesions(s []string) map[string]ml.SessionArgs {
	sessions := make(map[string]ml.SessionArgs)
	for i := range s {
		k := strings.Split(s[i], ",")
		i, err := strconv.ParseInt(k[2], 10, 64)
		if err != nil {
			fmt.Println(err)
		}
		sessions[k[0]] = ml.SessionArgs{Name: k[1], CreationDate: time.Unix(i, 0)}
	}
	return sessions
}
