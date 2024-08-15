package infra

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type ID uuid.UUID

func (id ID) String() string {
	return uuid.UUID(id).String()
}

type User struct {
	FirstName string
	LastName  string
	Biography string
}

type UserOut struct {
	ID ID
	User
}

type AppDB struct {
	data map[ID]User
}

func NewAppDB() *AppDB {
	return &AppDB{
		data: make(map[ID]User),
	}
}

func (db *AppDB) CreateUser(m *sync.Mutex, user User) ID {
	m.Lock()
	defer m.Unlock()
	id := ID(uuid.New())
	db.data[id] = user
	return id
}

func (db *AppDB) GetUser(m *sync.Mutex, id ID) (UserOut, error) {
	m.Lock()
	defer m.Unlock()
	user, ok := db.data[id]
	if !ok {
		return UserOut{}, errors.New("could not find user")
	}

	return UserOut{
		ID:   id,
		User: user,
	}, nil
}

func (db *AppDB) ListUsers(m *sync.Mutex) []UserOut {
	m.Lock()
	defer m.Unlock()
	users := make([]UserOut, 0, len(db.data))
	for id := range db.data {
		user := db.data[id]
		users = append(users, UserOut{
			ID:   id,
			User: user,
		})
	}
	return users
}

func (db *AppDB) UpdateUser(m *sync.Mutex, id ID, user User) error {
	m.Lock()
	defer m.Unlock()
	if _, ok := db.data[id]; !ok {
		return errors.New("could not find user")
	}
	db.data[id] = user

	return nil
}

func (db *AppDB) DeleteUser(m *sync.Mutex, id ID) error {
	m.Lock()
	defer m.Unlock()
	if _, ok := db.data[id]; !ok {
		return errors.New("could not find user")
	}
	delete(db.data, id)
	return nil
}
