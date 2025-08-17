package repositories

import (
	"bookmyshow-lld/internal/models"
	"time"
)

// UserRepository defines core user data access operations
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id string) (*models.User, error)
}

// MovieRepository defines core movie data access operations
type MovieRepository interface {
	Create(movie *models.Movie) error
	GetByID(id string) (*models.Movie, error)
	GetReleased() ([]*models.Movie, error) // For demo
}

// TheatreRepository defines core theatre data access operations
type TheatreRepository interface {
	Create(theatre *models.Theatre) error
	GetByID(id string) (*models.Theatre, error)
	Update(theatre *models.Theatre) error // Needed for adding screens
}

// ScreenRepository defines core screen data access operations
type ScreenRepository interface {
	Create(screen *models.Screen) error
	GetByID(id string) (*models.Screen, error)
	Update(screen *models.Screen) error // Needed for seat blocking/booking
}

// ShowRepository defines core show data access operations
type ShowRepository interface {
	Create(show *models.Show) error
	GetByID(id string) (*models.Show, error)
	GetByMovieID(movieID string) ([]*models.Show, error)                       // For demo
	CheckConflict(screenID string, startTime, endTime time.Time) (bool, error) // Business rule
}

// BookingRepository defines core booking data access operations
type BookingRepository interface {
	Create(booking *models.Booking) error
	GetByID(id string) (*models.Booking, error)
	Update(booking *models.Booking) error // Needed for confirming bookings
}

// PaymentRepository defines core payment data access operations
type PaymentRepository interface {
	Create(payment *models.Payment) error
	GetByID(id string) (*models.Payment, error)
	Update(payment *models.Payment) error // Needed for updating payment status
}
