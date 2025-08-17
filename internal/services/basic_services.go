package services

import (
	"bookmyshow-lld/internal/models"
	"bookmyshow-lld/internal/repositories"
	"time"
)

// UserServiceImpl implements UserService - demonstrates Repository Pattern
type UserServiceImpl struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

func (us *UserServiceImpl) CreateUser(name, email, phoneNumber string) (*models.User, error) {
	user, err := models.NewUser(name, email, phoneNumber)
	if err != nil {
		return nil, err
	}

	if err := us.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserServiceImpl) GetUser(id string) (*models.User, error) {
	return us.userRepo.GetByID(id)
}

// MovieServiceImpl implements MovieService - demonstrates Repository Pattern
type MovieServiceImpl struct {
	movieRepo repositories.MovieRepository
}

func NewMovieService(movieRepo repositories.MovieRepository) MovieService {
	return &MovieServiceImpl{
		movieRepo: movieRepo,
	}
}

func (ms *MovieServiceImpl) CreateMovie(title, description string, duration time.Duration, genre models.Genre, language models.Language, rating float32, releaseDate time.Time) (*models.Movie, error) {
	movie, err := models.NewMovie(title, description, duration, genre, language, rating, releaseDate)
	if err != nil {
		return nil, err
	}

	if err := ms.movieRepo.Create(movie); err != nil {
		return nil, err
	}

	return movie, nil
}

func (ms *MovieServiceImpl) GetMovie(id string) (*models.Movie, error) {
	return ms.movieRepo.GetByID(id)
}

func (ms *MovieServiceImpl) GetReleasedMovies() ([]*models.Movie, error) {
	return ms.movieRepo.GetReleased()
}

// TheatreServiceImpl implements TheatreService - demonstrates Repository Pattern + Business Logic
type TheatreServiceImpl struct {
	theatreRepo repositories.TheatreRepository
	screenRepo  repositories.ScreenRepository
}

func NewTheatreService(theatreRepo repositories.TheatreRepository, screenRepo repositories.ScreenRepository) TheatreService {
	return &TheatreServiceImpl{
		theatreRepo: theatreRepo,
		screenRepo:  screenRepo,
	}
}

func (ts *TheatreServiceImpl) CreateTheatre(name, address, city string) (*models.Theatre, error) {
	theatre, err := models.NewTheatre(name, address, city)
	if err != nil {
		return nil, err
	}

	if err := ts.theatreRepo.Create(theatre); err != nil {
		return nil, err
	}

	return theatre, nil
}

func (ts *TheatreServiceImpl) GetTheatre(id string) (*models.Theatre, error) {
	return ts.theatreRepo.GetByID(id)
}

func (ts *TheatreServiceImpl) AddScreen(theatreID string, screen *models.Screen) error {
	theatre, err := ts.theatreRepo.GetByID(theatreID)
	if err != nil {
		return err
	}

	theatre.AddScreen(screen)

	if err := ts.screenRepo.Create(screen); err != nil {
		return err
	}

	return ts.theatreRepo.Update(theatre)
}

// ShowServiceImpl implements ShowService - demonstrates business rules and validation
type ShowServiceImpl struct {
	showRepo    repositories.ShowRepository
	movieRepo   repositories.MovieRepository
	theatreRepo repositories.TheatreRepository
	screenRepo  repositories.ScreenRepository
}

func NewShowService(showRepo repositories.ShowRepository, movieRepo repositories.MovieRepository, theatreRepo repositories.TheatreRepository, screenRepo repositories.ScreenRepository) ShowService {
	return &ShowServiceImpl{
		showRepo:    showRepo,
		movieRepo:   movieRepo,
		theatreRepo: theatreRepo,
		screenRepo:  screenRepo,
	}
}

func (ss *ShowServiceImpl) CreateShow(movieID, theatreID, screenID string, startTime time.Time, basePrice float64) (*models.Show, error) {
	// Validate movie exists
	movie, err := ss.movieRepo.GetByID(movieID)
	if err != nil {
		return nil, err
	}

	// Validate theatre exists
	if _, err := ss.theatreRepo.GetByID(theatreID); err != nil {
		return nil, err
	}

	// Validate screen exists and belongs to theatre
	screen, err := ss.screenRepo.GetByID(screenID)
	if err != nil {
		return nil, err
	}

	if screen.TheatreID != theatreID {
		return nil, models.ErrInvalidShowData
	}

	// Check for scheduling conflicts - demonstrates business rules
	endTime := startTime.Add(movie.Duration)
	hasConflict, err := ss.showRepo.CheckConflict(screenID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	if hasConflict {
		return nil, models.ErrInvalidShowTime
	}

	// Create show
	show, err := models.NewShow(movieID, theatreID, screenID, startTime, basePrice, movie.Duration)
	if err != nil {
		return nil, err
	}

	if err := ss.showRepo.Create(show); err != nil {
		return nil, err
	}

	return show, nil
}

func (ss *ShowServiceImpl) GetShow(id string) (*models.Show, error) {
	return ss.showRepo.GetByID(id)
}

func (ss *ShowServiceImpl) GetShowsByMovie(movieID string) ([]*models.Show, error) {
	return ss.showRepo.GetByMovieID(movieID)
}
