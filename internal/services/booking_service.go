package services

import (
	"bookmyshow-lld/internal/models"
	"bookmyshow-lld/internal/repositories"
	"fmt"
	"sync"
)

// BookingServiceImpl implements BookingService - demonstrates Concurrency Control and Business Logic
type BookingServiceImpl struct {
	bookingRepo     repositories.BookingRepository
	showRepo        repositories.ShowRepository
	screenRepo      repositories.ScreenRepository
	theatreRepo     repositories.TheatreRepository
	movieRepo       repositories.MovieRepository
	paymentRepo     repositories.PaymentRepository
	notificationSvc NotificationService
	mutex           sync.RWMutex // Demonstrates thread-safe operations
}

// NewBookingService creates a new booking service
func NewBookingService(
	bookingRepo repositories.BookingRepository,
	showRepo repositories.ShowRepository,
	screenRepo repositories.ScreenRepository,
	theatreRepo repositories.TheatreRepository,
	movieRepo repositories.MovieRepository,
	paymentRepo repositories.PaymentRepository,
	notificationSvc NotificationService,
) BookingService {
	return &BookingServiceImpl{
		bookingRepo:     bookingRepo,
		showRepo:        showRepo,
		screenRepo:      screenRepo,
		theatreRepo:     theatreRepo,
		movieRepo:       movieRepo,
		paymentRepo:     paymentRepo,
		notificationSvc: notificationSvc,
	}
}

// CreateBooking creates a new booking with atomic seat blocking - demonstrates Concurrency Control
func (bs *BookingServiceImpl) CreateBooking(userID, showID string, seatIDs []string) (*models.Booking, error) {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	// Validate show
	show, err := bs.showRepo.GetByID(showID)
	if err != nil {
		return nil, err
	}

	if !show.CanBeBooked() {
		return nil, models.ErrShowNotBookable
	}

	// Get screen and validate seats
	screen, err := bs.screenRepo.GetByID(show.ScreenID)
	if err != nil {
		return nil, err
	}

	// Calculate total amount using Factory Pattern pricing
	totalAmount := 0.0
	for _, seatID := range seatIDs {
		seat, err := screen.GetSeat(seatID)
		if err != nil {
			return nil, err
		}
		if !seat.IsAvailable() {
			return nil, models.ErrSeatNotAvailable
		}
		totalAmount += seat.GetPrice()
	}

	// Block seats atomically - demonstrates atomic operations
	if err := screen.BlockSeats(seatIDs); err != nil {
		return nil, err
	}

	// Create booking
	booking, err := models.NewBooking(userID, showID, seatIDs, totalAmount)
	if err != nil {
		// Rollback seat blocking on failure
		bs.rollbackSeatBlocking(screen, seatIDs)
		return nil, err
	}

	// Save booking
	if err := bs.bookingRepo.Create(booking); err != nil {
		// Rollback seat blocking on failure
		bs.rollbackSeatBlocking(screen, seatIDs)
		return nil, err
	}

	// Update screen in repository
	if err := bs.screenRepo.Update(screen); err != nil {
		// Log error but don't fail the booking
		fmt.Printf("Warning: Failed to update screen after booking creation: %v\n", err)
	}

	return booking, nil
}

// GetBooking retrieves a booking by ID
func (bs *BookingServiceImpl) GetBooking(id string) (*models.Booking, error) {
	return bs.bookingRepo.GetByID(id)
}

// ConfirmBooking confirms a booking after successful payment - demonstrates Observer Pattern
func (bs *BookingServiceImpl) ConfirmBooking(bookingID, paymentID string) error {
	booking, err := bs.bookingRepo.GetByID(bookingID)
	if err != nil {
		return err
	}

	if err := booking.Confirm(paymentID); err != nil {
		return err
	}

	// Update booking in repository
	if err := bs.bookingRepo.Update(booking); err != nil {
		return err
	}

	// Book the actual seats
	show, err := bs.showRepo.GetByID(booking.ShowID)
	if err != nil {
		return err
	}

	screen, err := bs.screenRepo.GetByID(show.ScreenID)
	if err != nil {
		return err
	}

	// Book all seats
	for _, seatID := range booking.SeatIDs {
		seat, err := screen.GetSeat(seatID)
		if err != nil {
			continue // Log error but don't fail
		}
		if err := seat.Book(); err != nil {
			// Log error but continue
			fmt.Printf("Warning: Failed to book seat %s: %v\n", seatID, err)
		}
	}

	// Update screen
	if err := bs.screenRepo.Update(screen); err != nil {
		fmt.Printf("Warning: Failed to update screen after booking confirmation: %v\n", err)
	}

	// Send notification - demonstrates Observer Pattern
	if bs.notificationSvc != nil {
		bs.notificationSvc.SendBookingConfirmation(booking.UserID, booking.ID)
	}

	return nil
}

// GetBookingDetails retrieves detailed booking information - demonstrates Aggregate Construction
func (bs *BookingServiceImpl) GetBookingDetails(bookingID string) (*BookingDetails, error) {
	booking, err := bs.bookingRepo.GetByID(bookingID)
	if err != nil {
		return nil, err
	}

	show, err := bs.showRepo.GetByID(booking.ShowID)
	if err != nil {
		return nil, err
	}

	movie, err := bs.movieRepo.GetByID(show.MovieID)
	if err != nil {
		return nil, err
	}

	theatre, err := bs.theatreRepo.GetByID(show.TheatreID)
	if err != nil {
		return nil, err
	}

	screen, err := bs.screenRepo.GetByID(show.ScreenID)
	if err != nil {
		return nil, err
	}

	// Get seats
	var seats []*models.Seat
	for _, seatID := range booking.SeatIDs {
		seat, err := screen.GetSeat(seatID)
		if err == nil {
			seats = append(seats, seat)
		}
	}

	// Get payment if exists
	var payment *models.Payment
	if booking.PaymentID != "" {
		payment, _ = bs.paymentRepo.GetByID(booking.PaymentID)
	}

	return &BookingDetails{
		Booking: booking,
		Show:    show,
		Movie:   movie,
		Theatre: theatre,
		Screen:  screen,
		Seats:   seats,
		Payment: payment,
	}, nil
}

// Helper method to rollback seat blocking - demonstrates Error Handling
func (bs *BookingServiceImpl) rollbackSeatBlocking(screen *models.Screen, seatIDs []string) {
	for _, seatID := range seatIDs {
		seat, err := screen.GetSeat(seatID)
		if err == nil {
			seat.Unblock()
		}
	}
}
