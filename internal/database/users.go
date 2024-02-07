package database

import (
	"errors"

	"github.com/DavoReds/chirpy/internal/domain"
)

func (db *DB) CreateUser(email string, password []byte) (domain.User, error) {
	data, err := db.loadDB()
	if err != nil {
		return domain.User{}, err
	}

	_, err = db.GetUserByEmail(email)
	if err == nil {
		return domain.User{}, errors.New("User already exists")
	}

	var lastID int
	if len(data.Users) == 0 {
		lastID = 0
	} else {
		lastID = data.Users[len(data.Users)-1].ID
	}
	newUser := domain.User{
		ID:       lastID + 1,
		Email:    email,
		Password: password,
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

func (db *DB) GetUserByID(id int) (domain.User, error) {
	data, err := db.loadDB()
	if err != nil {
		return domain.User{}, err
	}

	for _, user := range data.Users {
		if user.ID == id {
			return user, nil
		}
	}

	return domain.User{}, errors.New("User doesn't exist")
}

func (db *DB) GetUserByEmail(email string) (domain.User, error) {
	data, err := db.loadDB()
	if err != nil {
		return domain.User{}, err
	}

	for _, user := range data.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return domain.User{}, errors.New("User doesn't exist")
}
