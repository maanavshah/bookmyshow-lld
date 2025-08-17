package models

import (
	"sync"

	"github.com/google/uuid"
)

// Screen represents a screen in a theatre
type Screen struct {
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	TheatreID  string           `json:"theatre_id"`
	Capacity   int              `json:"capacity"`
	Seats      map[string]*Seat `json:"seats"`
	seatsMutex sync.RWMutex
}

// NewScreen creates a new screen
func NewScreen(name, theatreID string) *Screen {
	return &Screen{
		ID:        uuid.New().String(),
		Name:      name,
		TheatreID: theatreID,
		Seats:     make(map[string]*Seat),
	}
}

// AddSeat adds a seat to the screen
func (s *Screen) AddSeat(seat *Seat) {
	s.seatsMutex.Lock()
	defer s.seatsMutex.Unlock()

	s.Seats[seat.ID] = seat
	s.Capacity++
}

// GetSeat retrieves a seat by ID (thread-safe)
func (s *Screen) GetSeat(seatID string) (*Seat, error) {
	s.seatsMutex.RLock()
	defer s.seatsMutex.RUnlock()

	seat, exists := s.Seats[seatID]
	if !exists {
		return nil, ErrSeatNotFound
	}
	return seat, nil
}

// GetAvailableSeats returns all available seats (thread-safe)
func (s *Screen) GetAvailableSeats() []*Seat {
	s.seatsMutex.RLock()
	defer s.seatsMutex.RUnlock()

	var availableSeats []*Seat
	for _, seat := range s.Seats {
		if seat.IsAvailable() {
			availableSeats = append(availableSeats, seat)
		}
	}
	return availableSeats
}

// GetSeatsByType returns seats of a specific type
func (s *Screen) GetSeatsByType(seatType SeatType) []*Seat {
	s.seatsMutex.RLock()
	defer s.seatsMutex.RUnlock()

	var seats []*Seat
	for _, seat := range s.Seats {
		if seat.Type == seatType {
			seats = append(seats, seat)
		}
	}
	return seats
}

// BlockSeats blocks multiple seats atomically
func (s *Screen) BlockSeats(seatIDs []string) error {
	s.seatsMutex.Lock()
	defer s.seatsMutex.Unlock()

	// First check if all seats are available
	for _, seatID := range seatIDs {
		seat, exists := s.Seats[seatID]
		if !exists {
			return ErrSeatNotFound
		}
		if !seat.IsAvailable() {
			return ErrSeatNotAvailable
		}
	}

	// Block all seats
	for _, seatID := range seatIDs {
		if err := s.Seats[seatID].Block(); err != nil {
			// Rollback previous blocks
			for i := 0; i < len(seatIDs); i++ {
				if seatIDs[i] == seatID {
					break
				}
				s.Seats[seatIDs[i]].Unblock()
			}
			return err
		}
	}

	return nil
}

// GetCapacity returns screen capacity
func (s *Screen) GetCapacity() int {
	s.seatsMutex.RLock()
	defer s.seatsMutex.RUnlock()
	return s.Capacity
}
