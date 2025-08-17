package models

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// BookingStatus represents the status of a booking
type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "PENDING"
	BookingStatusConfirmed BookingStatus = "CONFIRMED"
	BookingStatusCancelled BookingStatus = "CANCELLED"
	BookingStatusExpired   BookingStatus = "EXPIRED"
)

// Booking represents a ticket booking
type Booking struct {
	ID          string        `json:"id"`
	UserID      string        `json:"user_id"`
	ShowID      string        `json:"show_id"`
	SeatIDs     []string      `json:"seat_ids"`
	TotalAmount float64       `json:"total_amount"`
	Status      BookingStatus `json:"status"`
	BookingTime time.Time     `json:"booking_time"`
	ExpiryTime  time.Time     `json:"expiry_time"`
	PaymentID   string        `json:"payment_id,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	mutex       sync.RWMutex
}

// BookingTimeout represents the timeout for pending bookings
const BookingTimeout = 15 * time.Minute

// NewBooking creates a new booking
func NewBooking(userID, showID string, seatIDs []string, totalAmount float64) (*Booking, error) {
	if userID == "" || showID == "" || len(seatIDs) == 0 || totalAmount <= 0 {
		return nil, ErrInvalidBookingData
	}

	now := time.Now()
	return &Booking{
		ID:          uuid.New().String(),
		UserID:      userID,
		ShowID:      showID,
		SeatIDs:     seatIDs,
		TotalAmount: totalAmount,
		Status:      BookingStatusPending,
		BookingTime: now,
		ExpiryTime:  now.Add(BookingTimeout),
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// IsExpired checks if the booking has expired
func (b *Booking) IsExpired() bool {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return time.Now().After(b.ExpiryTime) && b.Status == BookingStatusPending
}

// Confirm confirms the booking after successful payment
func (b *Booking) Confirm(paymentID string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.Status != BookingStatusPending {
		return ErrBookingNotPending
	}

	if time.Now().After(b.ExpiryTime) {
		b.Status = BookingStatusExpired
		return ErrBookingExpired
	}

	b.Status = BookingStatusConfirmed
	b.PaymentID = paymentID
	b.UpdatedAt = time.Now()
	return nil
}

// Cancel cancels the booking
func (b *Booking) Cancel() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.Status == BookingStatusConfirmed {
		return ErrBookingAlreadyConfirmed
	}

	if b.Status == BookingStatusCancelled {
		return ErrBookingAlreadyCancelled
	}

	b.Status = BookingStatusCancelled
	b.UpdatedAt = time.Now()
	return nil
}

// Expire marks the booking as expired
func (b *Booking) Expire() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.Status != BookingStatusPending {
		return ErrBookingNotPending
	}

	b.Status = BookingStatusExpired
	b.UpdatedAt = time.Now()
	return nil
}

// GetStatus returns the current booking status (thread-safe)
func (b *Booking) GetStatus() BookingStatus {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.Status
}

// GetSeatCount returns number of seats booked
func (b *Booking) GetSeatCount() int {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return len(b.SeatIDs)
}

// TimeUntilExpiry returns time until booking expires
func (b *Booking) TimeUntilExpiry() time.Duration {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if b.Status != BookingStatusPending {
		return 0
	}

	remaining := time.Until(b.ExpiryTime)
	if remaining < 0 {
		return 0
	}
	return remaining
}

// CanBeCancelled checks if booking can be cancelled
func (b *Booking) CanBeCancelled() bool {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.Status == BookingStatusPending || b.Status == BookingStatusConfirmed
}
