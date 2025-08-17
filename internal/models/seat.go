package models

import (
	"sync"

	"github.com/google/uuid"
)

// SeatType represents different types of seats
type SeatType string

const (
	SeatTypeRegular  SeatType = "REGULAR"
	SeatTypePremium  SeatType = "PREMIUM"
	SeatTypeVIP      SeatType = "VIP"
	SeatTypeRecliner SeatType = "RECLINER"
)

// SeatStatus represents the status of a seat
type SeatStatus string

const (
	SeatStatusAvailable SeatStatus = "AVAILABLE"
	SeatStatusBooked    SeatStatus = "BOOKED"
	SeatStatusBlocked   SeatStatus = "BLOCKED"
)

// Seat represents a seat in a screen with thread-safe operations
type Seat struct {
	ID      string     `json:"id"`
	RowName string     `json:"row_name"`
	Number  int        `json:"number"`
	Type    SeatType   `json:"type"`
	Status  SeatStatus `json:"status"`
	Price   float64    `json:"price"`
	mutex   sync.RWMutex
}

// NewSeat creates a new seat
func NewSeat(rowName string, number int, seatType SeatType, price float64) *Seat {
	return &Seat{
		ID:      uuid.New().String(),
		RowName: rowName,
		Number:  number,
		Type:    seatType,
		Status:  SeatStatusAvailable,
		Price:   price,
	}
}

// IsAvailable checks if the seat is available for booking (thread-safe)
func (s *Seat) IsAvailable() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.Status == SeatStatusAvailable
}

// Block blocks the seat temporarily (thread-safe)
func (s *Seat) Block() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.Status != SeatStatusAvailable {
		return ErrSeatNotAvailable
	}

	s.Status = SeatStatusBlocked
	return nil
}

// Book books the seat (thread-safe)
func (s *Seat) Book() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.Status != SeatStatusBlocked {
		return ErrSeatNotBlocked
	}

	s.Status = SeatStatusBooked
	return nil
}

// Unblock unblocks the seat (thread-safe)
func (s *Seat) Unblock() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.Status != SeatStatusBlocked {
		return ErrSeatNotBlocked
	}

	s.Status = SeatStatusAvailable
	return nil
}

// GetStatus returns the current status (thread-safe)
func (s *Seat) GetStatus() SeatStatus {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.Status
}

// GetPrice returns the seat price
func (s *Seat) GetPrice() float64 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.Price
}

// GetSeatNumber returns formatted seat number
func (s *Seat) GetSeatNumber() string {
	return s.RowName + string(rune('0'+s.Number))
}
