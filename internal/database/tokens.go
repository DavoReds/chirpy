package database

import "time"

func (db *DB) WasTokenRevoked(token string) (bool, error) {
	data, err := db.loadDB()
	if err != nil {
		return false, err
	}

	_, exists := data.RevokedTokens[token]

	return exists, nil
}

func (db *DB) RevokeToken(token string) error {
	data, err := db.loadDB()
	if err != nil {
		return err
	}

	data.RevokedTokens[token] = time.Now()

	err = db.writeDB(data)
	if err != nil {
		return err
	}

	return nil
}
