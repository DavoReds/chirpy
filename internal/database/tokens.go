package database

func (db *DB) WasTokenRevoked(token string) (bool, error) {
	data, err := db.loadDB()
	if err != nil {
		return false, err
	}

	_, exists := data.RevokedTokens[token]

	return exists, nil
}
