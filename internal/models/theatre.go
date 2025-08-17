package models

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// Theatre represents a theatre with multiple screens
type Theatre struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	Address   string             `json:"address"`
	City      string             `json:"city"`
	Screens   map[string]*Screen `json:"screens"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	mutex     sync.RWMutex
}

// NewTheatre creates a new theatre
func NewTheatre(name, address, city string) (*Theatre, error) {
	if name == "" || address == "" || city == "" {
		return nil, ErrInvalidTheatreData
	}

	return &Theatre{
		ID:        uuid.New().String(),
		Name:      name,
		Address:   address,
		City:      city,
		Screens:   make(map[string]*Screen),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// AddScreen adds a screen to the theatre
func (t *Theatre) AddScreen(screen *Screen) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	screen.TheatreID = t.ID
	t.Screens[screen.ID] = screen
	t.UpdatedAt = time.Now()
}

// GetScreen retrieves a screen by ID
func (t *Theatre) GetScreen(screenID string) (*Screen, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	screen, exists := t.Screens[screenID]
	if !exists {
		return nil, ErrScreenNotFound
	}
	return screen, nil
}

// GetAllScreens returns all screens
func (t *Theatre) GetAllScreens() []*Screen {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	screens := make([]*Screen, 0, len(t.Screens))
	for _, screen := range t.Screens {
		screens = append(screens, screen)
	}
	return screens
}

// RemoveScreen removes a screen from the theatre
func (t *Theatre) RemoveScreen(screenID string) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if _, exists := t.Screens[screenID]; !exists {
		return ErrScreenNotFound
	}

	delete(t.Screens, screenID)
	t.UpdatedAt = time.Now()
	return nil
}

// GetTotalCapacity returns total capacity of all screens
func (t *Theatre) GetTotalCapacity() int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	totalCapacity := 0
	for _, screen := range t.Screens {
		totalCapacity += screen.GetCapacity()
	}
	return totalCapacity
}

// UpdateTheatre updates theatre information
func (t *Theatre) UpdateTheatre(name, address, city string) error {
	if name == "" || address == "" || city == "" {
		return ErrInvalidTheatreData
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.Name = name
	t.Address = address
	t.City = city
	t.UpdatedAt = time.Now()
	return nil
}
