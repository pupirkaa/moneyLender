package inmem

import "fmt"

type UserStorage struct {
	users map[string]string
}

func NewUserStorage() (usf UserStorage) {
	usf.users = make(map[string]string)
	return usf
}

func (usf UserStorage) Close() error {
	return nil
}

func (usf UserStorage) UserExist(name string) bool {
	_, ok := usf.users[name]
	fmt.Println("check user", name, " ", ok)
	return ok
}

func (usf UserStorage) UserAdd(name string, password string) {
	usf.users[name] = password
	fmt.Println(usf.users)
}

func (usf UserStorage) UserGet(name string) string {
	return usf.users[name]
}
