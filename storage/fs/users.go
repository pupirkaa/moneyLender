package fs

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type UserFileStorage struct {
	users      map[string]string
	path       string
	exitCh     chan bool
	syncDoneCh chan bool
}

func NewUserStorage(path string) (usf UserFileStorage) {
	usf.users = parseUsers(readFile(path))
	usf.path = path
	usf.exitCh = make(chan bool)
	usf.syncDoneCh = make(chan bool)
	go func() {
		t := time.NewTicker(1 * time.Minute)
		for {
			select {
			case <-usf.exitCh:
				usf.SaveUsersToFile()
				close(usf.syncDoneCh)
				return
			case <-t.C:
				usf.SaveUsersToFile()
			}
		}
	}()

	return usf
}

func (usf UserFileStorage) Close() error {
	close(usf.exitCh)
	<-usf.syncDoneCh
	return nil
}

func (usf UserFileStorage) UserExist(name string) (bool, error) {
	_, ok := usf.users[name]
	return ok, nil
}

func (usf UserFileStorage) UserAdd(name string, password string) error {
	usf.users[name] = password
	return nil
}

func (usf UserFileStorage) UserGet(name string) (string, error) {
	return usf.users[name], nil
}

func (usf UserFileStorage) SaveUsersToFile() {
	f, err := os.Create(usf.path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file: %v", err)
		os.Exit(1)
	}
	for k := range usf.users {
		_, err = f.WriteString(fmt.Sprintf("%v,%v\n", k, usf.users[k]))
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("Saved users")
	defer f.Close()
}

func readFile(s string) []string {
	var text []string
	f, err := os.Open(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file: %v", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	return text
}

func parseUsers(s []string) map[string]string {
	users := make(map[string]string)

	for i := range s {
		k := strings.Split(s[i], ",")
		users[k[0]] = k[1]
	}
	return users
}
