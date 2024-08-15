package infra

import (
	"errors"
	"fmt"

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

func (db *AppDB) CreateUser(user User) ID {
	id := ID(uuid.New())
	fmt.Println(id.String())
	db.data[id] = user
	return id
}

func (db *AppDB) GetUser(id ID) (UserOut, error) {
	user, ok := db.data[id]
	if !ok {
		return UserOut{}, errors.New("could not find user")
	}

	return UserOut{
		ID:   id,
		User: user,
	}, nil
}

func (db *AppDB) ListUsers() []UserOut {
	users := make([]UserOut, 0, len(db.data))
	for id := range db.data {
		fmt.Println(id.String())
		user := db.data[id]
		users = append(users, UserOut{
			ID:   id,
			User: user,
		})
	}
	return users
}

func (db *AppDB) UpdateUser(id ID, user User) error {
	if _, ok := db.data[id]; !ok {
		return errors.New("could not find user")
	}
	db.data[id] = user

	return nil
}

func (db *AppDB) DeleteUser(id ID) error {
	if _, ok := db.data[id]; !ok {
		return errors.New("could not find user")
	}
	delete(db.data, id)
	return nil
}
