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

	lastID := maxIntKey(data.Users)
	newID := lastID + 1

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return domain.User{}, err
	}

	newUser := domain.User{
		ID:       newID,
		Email:    email,
		Password: hashedPassword,
	}

	data.Users[newID] = newUser

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

	return getValues(data.Users), nil
}

func (db *DB) GetUserByID(id int) (domain.User, error) {
	data, err := db.loadDB()
	if err != nil {
		return domain.User{}, err
	}

	user, exists := data.Users[id]
	if !exists {
		return domain.User{}, errors.New("User doesn't exist")
	}

	return user, nil
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

func (db *DB) UpdateUser(id int, email string, password []byte) (domain.User, error) {
	data, err := db.loadDB()
	if err != nil {
		return domain.User{}, err
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return domain.User{}, err
	}

	newUser := domain.User{
		ID:       id,
		Email:    email,
		Password: hashedPassword,
	}
	data.Users[id] = newUser

	if err = db.writeDB(data); err != nil {
		return domain.User{}, err
	}

	return newUser, nil
}
