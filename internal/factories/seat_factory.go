package factories

import (
	"bookmyshow-lld/internal/models"
	"fmt"
)

// SeatFactory creates different types of seats
type SeatFactory struct{}

// NewSeatFactory creates a new seat factory
func NewSeatFactory() *SeatFactory {
	return &SeatFactory{}
}

// CreateSeat creates a seat based on type with appropriate pricing
func (sf *SeatFactory) CreateSeat(rowName string, number int, seatType models.SeatType, basePrice float64) *models.Seat {
	price := sf.calculatePrice(seatType, basePrice)
	return models.NewSeat(rowName, number, seatType, price)
}

// CreateSeatsForScreen creates seats for an entire screen
func (sf *SeatFactory) CreateSeatsForScreen(screenID string, config ScreenConfig, basePrice float64) []*models.Seat {
	var seats []*models.Seat

	for _, rowConfig := range config.Rows {
		for i := 1; i <= rowConfig.Count; i++ {
			seat := sf.CreateSeat(rowConfig.Name, i, rowConfig.Type, basePrice)
			seats = append(seats, seat)
		}
	}

	return seats
}

// CreateDefaultScreenSeats creates a default seat configuration
func (sf *SeatFactory) CreateDefaultScreenSeats(basePrice float64) []*models.Seat {
	config := ScreenConfig{
		Rows: []RowConfig{
			{Name: "A", Count: 10, Type: models.SeatTypeVIP},
			{Name: "B", Count: 12, Type: models.SeatTypeVIP},
			{Name: "C", Count: 14, Type: models.SeatTypePremium},
			{Name: "D", Count: 14, Type: models.SeatTypePremium},
			{Name: "E", Count: 16, Type: models.SeatTypeRegular},
			{Name: "F", Count: 16, Type: models.SeatTypeRegular},
			{Name: "G", Count: 18, Type: models.SeatTypeRegular},
			{Name: "H", Count: 18, Type: models.SeatTypeRegular},
		},
	}

	return sf.CreateSeatsForScreen("", config, basePrice)
}

// calculatePrice calculates price based on seat type
func (sf *SeatFactory) calculatePrice(seatType models.SeatType, basePrice float64) float64 {
	multiplier := sf.getPriceMultiplier(seatType)
	return basePrice * multiplier
}

// getPriceMultiplier returns price multiplier for different seat types
func (sf *SeatFactory) getPriceMultiplier(seatType models.SeatType) float64 {
	switch seatType {
	case models.SeatTypeVIP:
		return 2.0
	case models.SeatTypePremium:
		return 1.5
	case models.SeatTypeRecliner:
		return 2.5
	case models.SeatTypeRegular:
		return 1.0
	default:
		return 1.0
	}
}

// ValidateSeatType validates if seat type is supported
func (sf *SeatFactory) ValidateSeatType(seatType models.SeatType) error {
	switch seatType {
	case models.SeatTypeRegular, models.SeatTypePremium, models.SeatTypeVIP, models.SeatTypeRecliner:
		return nil
	default:
		return fmt.Errorf("unsupported seat type: %s", seatType)
	}
}

// ScreenConfig represents screen seat configuration
type ScreenConfig struct {
	Rows []RowConfig `json:"rows"`
}

// RowConfig represents row configuration
type RowConfig struct {
	Name  string          `json:"name"`
	Count int             `json:"count"`
	Type  models.SeatType `json:"type"`
}

// GetSeatTypeInfo returns information about seat types
func (sf *SeatFactory) GetSeatTypeInfo() map[models.SeatType]SeatTypeInfo {
	return map[models.SeatType]SeatTypeInfo{
		models.SeatTypeRegular: {
			Name:        "Regular",
			Description: "Standard seating with basic comfort",
			Multiplier:  1.0,
		},
		models.SeatTypePremium: {
			Name:        "Premium",
			Description: "Enhanced comfort with extra legroom",
			Multiplier:  1.5,
		},
		models.SeatTypeVIP: {
			Name:        "VIP",
			Description: "Luxury seating with premium amenities",
			Multiplier:  2.0,
		},
		models.SeatTypeRecliner: {
			Name:        "Recliner",
			Description: "Fully reclining seats with maximum comfort",
			Multiplier:  2.5,
		},
	}
}

// SeatTypeInfo contains information about seat types
type SeatTypeInfo struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Multiplier  float64 `json:"multiplier"`
}
