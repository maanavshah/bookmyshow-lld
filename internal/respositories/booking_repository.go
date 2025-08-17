package repositories

import (
	"bookmyshow-lld/internal/models"
	"sync"
	"time"
)

// MemoryShowRepository implements ShowRepository - demonstrates Repository Pattern
type MemoryShowRepository struct {
	shows map[string]*models.Show
	mutex sync.RWMutex
}

func NewMemoryShowRepository() ShowRepository {
	return &MemoryShowRepository{
		shows: make(map[string]*models.Show),
	}
}

func (r *MemoryShowRepository) Create(show *models.Show) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.shows[show.ID] = show
	return nil
}

func (r *MemoryShowRepository) GetByID(id string) (*models.Show, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	show, exists := r.shows[id]
	if !exists {
		return nil, models.ErrShowNotFound
	}
	return show, nil
}

func (r *MemoryShowRepository) GetByMovieID(movieID string) ([]*models.Show, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var shows []*models.Show
	for _, show := range r.shows {
		if show.MovieID == movieID {
			shows = append(shows, show)
		}
	}
	return shows, nil
}

func (r *MemoryShowRepository) CheckConflict(screenID string, startTime, endTime time.Time) (bool, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, show := range r.shows {
		if show.ScreenID == screenID {
			// Check for time overlap - demonstrates business rules
			if startTime.Before(show.EndTime) && endTime.After(show.StartTime) {
				return true, nil
			}
		}
	}
	return false, nil
}

// MemoryBookingRepository implements BookingRepository - demonstrates Repository Pattern
type MemoryBookingRepository struct {
	bookings map[string]*models.Booking
	mutex    sync.RWMutex
}

func NewMemoryBookingRepository() BookingRepository {
	return &MemoryBookingRepository{
		bookings: make(map[string]*models.Booking),
	}
}

func (r *MemoryBookingRepository) Create(booking *models.Booking) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.bookings[booking.ID] = booking
	return nil
}

func (r *MemoryBookingRepository) GetByID(id string) (*models.Booking, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	booking, exists := r.bookings[id]
	if !exists {
		return nil, models.ErrBookingNotFound
	}
	return booking, nil
}

func (r *MemoryBookingRepository) Update(booking *models.Booking) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.bookings[booking.ID]; !exists {
		return models.ErrBookingNotFound
	}

	r.bookings[booking.ID] = booking
	return nil
}

// MemoryPaymentRepository implements PaymentRepository - demonstrates Repository Pattern
type MemoryPaymentRepository struct {
	payments map[string]*models.Payment
	mutex    sync.RWMutex
}

func NewMemoryPaymentRepository() PaymentRepository {
	return &MemoryPaymentRepository{
		payments: make(map[string]*models.Payment),
	}
}

func (r *MemoryPaymentRepository) Create(payment *models.Payment) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.payments[payment.ID] = payment
	return nil
}

func (r *MemoryPaymentRepository) GetByID(id string) (*models.Payment, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	payment, exists := r.payments[id]
	if !exists {
		return nil, models.ErrPaymentNotFound
	}
	return payment, nil
}

func (r *MemoryPaymentRepository) Update(payment *models.Payment) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.payments[payment.ID]; !exists {
		return models.ErrPaymentNotFound
	}

	r.payments[payment.ID] = payment
	return nil
}
