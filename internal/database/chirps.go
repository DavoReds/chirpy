package database

import (
	"errors"

	"github.com/DavoReds/chirpy/internal/domain"
)

func (db *DB) CreateChirp(body string) (domain.Chirp, error) {
	data, err := db.loadDB()
	if err != nil {
		return domain.Chirp{}, err
	}

	lastID := maxIntKey(data.Chirps)
	newID := lastID + 1
	newChirp := domain.Chirp{
		Body: body,
		ID:   newID,
	}

	data.Chirps[newID] = newChirp

	if err = db.writeDB(data); err != nil {
		return domain.Chirp{}, err
	}

	return newChirp, nil
}

func (db *DB) GetChirps() ([]domain.Chirp, error) {
	data, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	return getValues(data.Chirps), nil
}

func (db *DB) GetChirpByID(id int) (domain.Chirp, error) {
	data, err := db.loadDB()
	if err != nil {
		return domain.Chirp{}, err
	}

	for _, chirp := range data.Chirps {
		if chirp.ID == id {
			return chirp, nil
		}
	}

	return domain.Chirp{}, errors.New("Doesn't exist")
}
