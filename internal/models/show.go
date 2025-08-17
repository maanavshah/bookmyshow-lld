package models

import (
	"time"

	"github.com/google/uuid"
)

// Show represents a movie show at a specific theatre and time
type Show struct {
	ID        string    `json:"id"`
	MovieID   string    `json:"movie_id"`
	TheatreID string    `json:"theatre_id"`
	ScreenID  string    `json:"screen_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	BasePrice float64   `json:"base_price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewShow creates a new show with validation
func NewShow(movieID, theatreID, screenID string, startTime time.Time, basePrice float64, movieDuration time.Duration) (*Show, error) {
	if movieID == "" || theatreID == "" || screenID == "" || basePrice <= 0 {
		return nil, ErrInvalidShowData
	}

	if startTime.Before(time.Now()) {
		return nil, ErrInvalidShowTime
	}

	endTime := startTime.Add(movieDuration)

	return &Show{
		ID:        uuid.New().String(),
		MovieID:   movieID,
		TheatreID: theatreID,
		ScreenID:  screenID,
		StartTime: startTime,
		EndTime:   endTime,
		BasePrice: basePrice,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// IsActive checks if the show is currently active
func (s *Show) IsActive() bool {
	now := time.Now()
	return now.After(s.StartTime) && now.Before(s.EndTime)
}

// IsUpcoming checks if the show is scheduled for the future
func (s *Show) IsUpcoming() bool {
	return time.Now().Before(s.StartTime)
}

// IsCompleted checks if the show has ended
func (s *Show) IsCompleted() bool {
	return time.Now().After(s.EndTime)
}

// CanBeBooked checks if the show can still be booked
func (s *Show) CanBeBooked() bool {
	// Allow booking until 30 minutes after start time
	bookingCutoff := s.StartTime.Add(30 * time.Minute)
	return time.Now().Before(bookingCutoff)
}

// UpdateShow updates show information
func (s *Show) UpdateShow(startTime time.Time, basePrice float64, movieDuration time.Duration) error {
	if startTime.Before(time.Now()) || basePrice <= 0 {
		return ErrInvalidShowData
	}

	s.StartTime = startTime
	s.EndTime = startTime.Add(movieDuration)
	s.BasePrice = basePrice
	s.UpdatedAt = time.Now()
	return nil
}

// GetDuration returns the show duration
func (s *Show) GetDuration() time.Duration {
	return s.EndTime.Sub(s.StartTime)
}

// TimeUntilStart returns duration until show starts
func (s *Show) TimeUntilStart() time.Duration {
	if s.IsUpcoming() {
		return s.StartTime.Sub(time.Now())
	}
	return 0
}
