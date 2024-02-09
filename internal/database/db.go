package database

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"sync"
	"time"

	"github.com/DavoReds/chirpy/internal/domain"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type Data struct {
	Chirps        []domain.Chirp       `json:"chirps"`
	Users         map[int]domain.User  `json:"users"`
	RevokedTokens map[string]time.Time `json:"revoked_tokens"`
}

func NewDB(path string) *DB {
	database := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}

	return database
}

func emptyData() Data {
	return Data{
		Chirps:        []domain.Chirp{},
		Users:         make(map[int]domain.User),
		RevokedTokens: make(map[string]time.Time),
	}
}

func (db *DB) loadDB() (Data, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	fileData, err := os.ReadFile(db.path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return emptyData(), nil
		}

		return emptyData(), err
	}

	data := Data{}
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		return emptyData(), err
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
