package services

import (
	"bookmyshow-lld/internal/models"
	"time"
)

// UserService defines core user operations for LLD learning
type UserService interface {
	CreateUser(name, email, phoneNumber string) (*models.User, error)
	GetUser(id string) (*models.User, error)
}

// MovieService defines core movie operations for LLD learning
type MovieService interface {
	CreateMovie(title, description string, duration time.Duration, genre models.Genre, language models.Language, rating float32, releaseDate time.Time) (*models.Movie, error)
	GetMovie(id string) (*models.Movie, error)
	GetReleasedMovies() ([]*models.Movie, error) // Needed for demo
}

// TheatreService defines core theatre operations for LLD learning
type TheatreService interface {
	CreateTheatre(name, address, city string) (*models.Theatre, error)
	GetTheatre(id string) (*models.Theatre, error)
	AddScreen(theatreID string, screen *models.Screen) error // Core to booking flow
}

// ShowService defines core show operations for LLD learning
type ShowService interface {
	CreateShow(movieID, theatreID, screenID string, startTime time.Time, basePrice float64) (*models.Show, error)
	GetShow(id string) (*models.Show, error)
	GetShowsByMovie(movieID string) ([]*models.Show, error) // Needed for demo
}

// BookingService defines core booking operations for LLD learning
type BookingService interface {
	CreateBooking(userID, showID string, seatIDs []string) (*models.Booking, error)
	GetBooking(id string) (*models.Booking, error)
	ConfirmBooking(bookingID, paymentID string) error
	GetBookingDetails(bookingID string) (*BookingDetails, error)
}

// PaymentService defines core payment operations for LLD learning (Strategy Pattern)
type PaymentService interface {
	ProcessPayment(bookingID string, paymentMethod models.PaymentMethod) (*models.Payment, error)
	GetPayment(id string) (*models.Payment, error)
}

// PaymentGateway defines payment gateway operations (Strategy Pattern)
type PaymentGateway interface {
	ProcessPayment(amount float64, method models.PaymentMethod, metadata map[string]string) (*PaymentResult, error)
}

// NotificationService defines notification operations (Observer Pattern)
type NotificationService interface {
	SendBookingConfirmation(userID, bookingID string) error
}

// BookingDetails represents detailed booking information
type BookingDetails struct {
	Booking *models.Booking `json:"booking"`
	Show    *models.Show    `json:"show"`
	Movie   *models.Movie   `json:"movie"`
	Theatre *models.Theatre `json:"theatre"`
	Screen  *models.Screen  `json:"screen"`
	Seats   []*models.Seat  `json:"seats"`
	Payment *models.Payment `json:"payment,omitempty"`
}

// PaymentResult represents payment processing result (Strategy Pattern)
type PaymentResult struct {
	Success       bool   `json:"success"`
	TransactionID string `json:"transaction_id"`
	Response      string `json:"response"`
	ErrorMessage  string `json:"error_message,omitempty"`
}
