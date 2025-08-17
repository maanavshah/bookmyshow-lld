package services

import (
	"fmt"
	"log"
)

// NotificationServiceImpl implements NotificationService - demonstrates Observer Pattern
type NotificationServiceImpl struct {
	// In a real implementation, this would have email/SMS service clients
}

// NewNotificationService creates a new notification service
func NewNotificationService() NotificationService {
	return &NotificationServiceImpl{}
}

// SendBookingConfirmation sends booking confirmation notification - demonstrates Observer Pattern
func (ns *NotificationServiceImpl) SendBookingConfirmation(userID, bookingID string) error {
	message := fmt.Sprintf("Booking confirmed! Booking ID: %s for User: %s", bookingID, userID)
	log.Printf("📧 NOTIFICATION: %s", message)

	// In real implementation:
	// - Send email confirmation
	// - Send SMS notification
	// - Push notification to mobile app
	// - Update user's notification preferences

	return nil
}
