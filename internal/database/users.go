package database

import (
	"errors"

	"github.com/DavoReds/chirpy/internal/domain"
)

func (db *DB) CreateUser(email string) (domain.User, error) {
	data, err := db.loadDB()
	if err != nil {
		return domain.User{}, err
	}

	var lastID int
	if len(data.Users) == 0 {
		lastID = 0
	} else {
		lastID = data.Users[len(data.Users)-1].ID
	}
	newUser := domain.User{
		Email: email,
		ID:    lastID + 1,
	}

	data.Users = append(data.Users, newUser)

	if err = db.writeDB(data); err != nil {
		return domain.User{}, err
	}

	return newUser, nil
}

func (db *DB) GetUsers() ([]domain.User, error) {
	data, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	return data.Users, nil
}

func (db *DB) GetUser(id int) (domain.User, error) {
	data, err := db.loadDB()
	if err != nil {
		return domain.User{}, err
	}

	for _, user := range data.Users {
		if user.ID == id {
			return user, nil
		}
	}

	return domain.User{}, errors.New("Doesn't exist")
}
