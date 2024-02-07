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

	var lastID int
	if len(data.Chirps) == 0 {
		lastID = 0
	} else {
		lastID = data.Chirps[len(data.Chirps)-1].ID
	}
	newChirp := domain.Chirp{
		Body: body,
		ID:   lastID + 1,
	}

	data.Chirps = append(data.Chirps, newChirp)

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

	return data.Chirps, nil
}

func (db *DB) GetChirp(id int) (domain.Chirp, error) {
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
