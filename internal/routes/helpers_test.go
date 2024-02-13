package routes

import (
	"reflect"
	"testing"

	"github.com/DavoReds/chirpy/internal/domain"
)

func TestSortChirpsAsc(t *testing.T) {
	chirps := []domain.Chirp{
		{ID: 4, AuthorID: 4, Body: "Chirp 4"},
		{ID: 2, AuthorID: 2, Body: "Chirp 2"},
		{ID: 6, AuthorID: 6, Body: "Chirp 6"},
		{ID: 5, AuthorID: 5, Body: "Chirp 5"},
		{ID: 3, AuthorID: 3, Body: "Chirp 3"},
		{ID: 1, AuthorID: 1, Body: "Chirp 1"},
	}

	expected := []domain.Chirp{
		{ID: 1, AuthorID: 1, Body: "Chirp 1"},
		{ID: 2, AuthorID: 2, Body: "Chirp 2"},
		{ID: 3, AuthorID: 3, Body: "Chirp 3"},
		{ID: 4, AuthorID: 4, Body: "Chirp 4"},
		{ID: 5, AuthorID: 5, Body: "Chirp 5"},
		{ID: 6, AuthorID: 6, Body: "Chirp 6"},
	}

	sortChirps(chirps, "asc")

	if !reflect.DeepEqual(chirps, expected) {
		t.Errorf("Expected %v, got %v", expected, chirps)
	}
}

func TestSortChirpsDesc(t *testing.T) {
	chirps := []domain.Chirp{
		{ID: 4, AuthorID: 4, Body: "Chirp 4"},
		{ID: 2, AuthorID: 2, Body: "Chirp 2"},
		{ID: 6, AuthorID: 6, Body: "Chirp 6"},
		{ID: 5, AuthorID: 5, Body: "Chirp 5"},
		{ID: 3, AuthorID: 3, Body: "Chirp 3"},
		{ID: 1, AuthorID: 1, Body: "Chirp 1"},
	}

	expected := []domain.Chirp{
		{ID: 6, AuthorID: 6, Body: "Chirp 6"},
		{ID: 5, AuthorID: 5, Body: "Chirp 5"},
		{ID: 4, AuthorID: 4, Body: "Chirp 4"},
		{ID: 3, AuthorID: 3, Body: "Chirp 3"},
		{ID: 2, AuthorID: 2, Body: "Chirp 2"},
		{ID: 1, AuthorID: 1, Body: "Chirp 1"},
	}

	sortChirps(chirps, "desc")

	if !reflect.DeepEqual(chirps, expected) {
		t.Errorf("Expected %v, got %v", expected, chirps)
	}
}
