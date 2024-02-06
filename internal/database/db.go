package database

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"sync"

	"github.com/DavoReds/chirpy/internal/domain"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type Data struct {
	Chirps []domain.Chirp `json:"chirps"`
}

func NewDB(path string) *DB {
	database := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}

	return database
}

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

func (db *DB) loadDB() (Data, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	fileData, err := os.ReadFile(db.path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Data{}, nil
		}

		return Data{}, err
	}

	data := Data{}
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		return Data{}, err
	}

	return data, nil
}

func (db *DB) writeDB(data Data) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, jsonData, 0666)
	if err != nil {
		return err
	}

	return nil
}
