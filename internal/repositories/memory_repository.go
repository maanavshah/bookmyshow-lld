package repositories

import (
	"bookmyshow-lld/internal/models"
	"sync"
)

// MemoryUserRepository implements UserRepository - demonstrates Repository Pattern
type MemoryUserRepository struct {
	users map[string]*models.User
	mutex sync.RWMutex
}

func NewMemoryUserRepository() UserRepository {
	return &MemoryUserRepository{
		users: make(map[string]*models.User),
	}
}

func (r *MemoryUserRepository) Create(user *models.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Simple validation - prevent duplicate emails
	for _, existingUser := range r.users {
		if existingUser.Email == user.Email {
			return models.ErrInvalidUserData
		}
	}

	r.users[user.ID] = user
	return nil
}

func (r *MemoryUserRepository) GetByID(id string) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, models.ErrUserNotFound
	}
	return user, nil
}

// MemoryMovieRepository implements MovieRepository - demonstrates Repository Pattern
type MemoryMovieRepository struct {
	movies map[string]*models.Movie
	mutex  sync.RWMutex
}

func NewMemoryMovieRepository() MovieRepository {
	return &MemoryMovieRepository{
		movies: make(map[string]*models.Movie),
	}
}

func (r *MemoryMovieRepository) Create(movie *models.Movie) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.movies[movie.ID] = movie
	return nil
}

func (r *MemoryMovieRepository) GetByID(id string) (*models.Movie, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	movie, exists := r.movies[id]
	if !exists {
		return nil, models.ErrMovieNotFound
	}
	return movie, nil
}

func (r *MemoryMovieRepository) GetReleased() ([]*models.Movie, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var movies []*models.Movie
	for _, movie := range r.movies {
		if movie.IsReleased() {
			movies = append(movies, movie)
		}
	}
	return movies, nil
}

// MemoryTheatreRepository implements TheatreRepository - demonstrates Repository Pattern
type MemoryTheatreRepository struct {
	theatres map[string]*models.Theatre
	mutex    sync.RWMutex
}

func NewMemoryTheatreRepository() TheatreRepository {
	return &MemoryTheatreRepository{
		theatres: make(map[string]*models.Theatre),
	}
}

func (r *MemoryTheatreRepository) Create(theatre *models.Theatre) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.theatres[theatre.ID] = theatre
	return nil
}

func (r *MemoryTheatreRepository) GetByID(id string) (*models.Theatre, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	theatre, exists := r.theatres[id]
	if !exists {
		return nil, models.ErrTheatreNotFound
	}
	return theatre, nil
}

func (r *MemoryTheatreRepository) Update(theatre *models.Theatre) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.theatres[theatre.ID]; !exists {
		return models.ErrTheatreNotFound
	}

	r.theatres[theatre.ID] = theatre
	return nil
}

// MemoryScreenRepository implements ScreenRepository - demonstrates Repository Pattern
type MemoryScreenRepository struct {
	screens map[string]*models.Screen
	mutex   sync.RWMutex
}

func NewMemoryScreenRepository() ScreenRepository {
	return &MemoryScreenRepository{
		screens: make(map[string]*models.Screen),
	}
}

func (r *MemoryScreenRepository) Create(screen *models.Screen) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.screens[screen.ID] = screen
	return nil
}

func (r *MemoryScreenRepository) GetByID(id string) (*models.Screen, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	screen, exists := r.screens[id]
	if !exists {
		return nil, models.ErrScreenNotFound
	}
	return screen, nil
}

func (r *MemoryScreenRepository) Update(screen *models.Screen) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.screens[screen.ID]; !exists {
		return models.ErrScreenNotFound
	}

	r.screens[screen.ID] = screen
	return nil
}
