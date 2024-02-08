package database

import "golang.org/x/crypto/bcrypt"

func maxIntKey[T any](m map[int]T) int {
	maxKey := 0
	for key := range m {
		if key > maxKey {
			maxKey = key
		}
	}

	return maxKey
}

func getValues[K comparable, V any](m map[K]V) []V {
	values := []V{}
	for _, value := range m {
		values = append(values, value)
	}

	return values
}

func hashPassword(password []byte) ([]byte, error) {
	encrypted, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, err
	}

	return encrypted, nil
}
