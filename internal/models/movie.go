package models

import (
	"time"

	"github.com/google/uuid"
)

// Genre represents movie genres
type Genre string

const (
	GenreAction   Genre = "ACTION"
	GenreComedy   Genre = "COMEDY"
	GenreDrama    Genre = "DRAMA"
	GenreHorror   Genre = "HORROR"
	GenreRomance  Genre = "ROMANCE"
	GenreSciFi    Genre = "SCI_FI"
	GenreThriller Genre = "THRILLER"
)

// Language represents movie languages
type Language string

const (
	LanguageEnglish Language = "ENGLISH"
	LanguageHindi   Language = "HINDI"
	LanguageTamil   Language = "TAMIL"
	LanguageTelugu  Language = "TELUGU"
)

// Movie represents a movie in the system
type Movie struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Duration    time.Duration `json:"duration"`
	Genre       Genre         `json:"genre"`
	Language    Language      `json:"language"`
	Rating      float32       `json:"rating"`
	ReleaseDate time.Time     `json:"release_date"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// NewMovie creates a new movie with validation
func NewMovie(title, description string, duration time.Duration, genre Genre, language Language, rating float32, releaseDate time.Time) (*Movie, error) {
	if title == "" || duration <= 0 || rating < 0 || rating > 10 {
		return nil, ErrInvalidMovieData
	}

	return &Movie{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Duration:    duration,
		Genre:       genre,
		Language:    language,
		Rating:      rating,
		ReleaseDate: releaseDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// UpdateMovie updates movie information
func (m *Movie) UpdateMovie(title, description string, rating float32) error {
	if title == "" || rating < 0 || rating > 10 {
		return ErrInvalidMovieData
	}

	m.Title = title
	m.Description = description
	m.Rating = rating
	m.UpdatedAt = time.Now()
	return nil
}

// IsReleased checks if the movie has been released
func (m *Movie) IsReleased() bool {
	return time.Now().After(m.ReleaseDate)
}
