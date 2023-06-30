package inmem

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
	return ok
}

func (usf UserStorage) UserAdd(name string, password string) {
	usf.users[name] = password
}

func (usf UserStorage) UserGet(name string) string {
	return usf.users[name]
}
