package services

import (
	"bookmyshow-lld/internal/models"
	"bookmyshow-lld/internal/repositories"
)

// PaymentServiceImpl implements PaymentService - demonstrates Strategy Pattern
type PaymentServiceImpl struct {
	paymentRepo     repositories.PaymentRepository
	bookingRepo     repositories.BookingRepository
	paymentGateway  PaymentGateway // Strategy Pattern - different payment methods
	notificationSvc NotificationService
}

// NewPaymentService creates a new payment service
func NewPaymentService(
	paymentRepo repositories.PaymentRepository,
	bookingRepo repositories.BookingRepository,
	paymentGateway PaymentGateway,
	notificationSvc NotificationService,
) PaymentService {
	return &PaymentServiceImpl{
		paymentRepo:     paymentRepo,
		bookingRepo:     bookingRepo,
		paymentGateway:  paymentGateway,
		notificationSvc: notificationSvc,
	}
}

// ProcessPayment processes a payment for a booking - demonstrates Strategy Pattern
func (ps *PaymentServiceImpl) ProcessPayment(bookingID string, paymentMethod models.PaymentMethod) (*models.Payment, error) {
	// Get booking
	booking, err := ps.bookingRepo.GetByID(bookingID)
	if err != nil {
		return nil, err
	}

	if booking.GetStatus() != models.BookingStatusPending {
		return nil, models.ErrBookingNotPending
	}

	if booking.IsExpired() {
		return nil, models.ErrBookingExpired
	}

	// Create payment record
	payment, err := models.NewPayment(bookingID, booking.UserID, booking.TotalAmount, paymentMethod)
	if err != nil {
		return nil, err
	}

	// Save payment
	if err := ps.paymentRepo.Create(payment); err != nil {
		return nil, err
	}

	// Process payment through gateway using Strategy Pattern
	metadata := ps.buildPaymentMetadata(paymentMethod, booking)
	result, err := ps.paymentGateway.ProcessPayment(booking.TotalAmount, paymentMethod, metadata)
	if err != nil {
		payment.MarkFailed(err.Error())
		ps.paymentRepo.Update(payment)
		return payment, err
	}

	if result.Success {
		payment.MarkSuccess(result.TransactionID, result.Response)
	} else {
		payment.MarkFailed(result.ErrorMessage)
	}

	// Update payment
	if err := ps.paymentRepo.Update(payment); err != nil {
		return payment, err
	}

	return payment, nil
}

// GetPayment retrieves a payment by ID
func (ps *PaymentServiceImpl) GetPayment(id string) (*models.Payment, error) {
	return ps.paymentRepo.GetByID(id)
}

// buildPaymentMetadata builds metadata for payment processing - demonstrates Strategy Pattern setup
func (ps *PaymentServiceImpl) buildPaymentMetadata(method models.PaymentMethod, booking *models.Booking) map[string]string {
	metadata := map[string]string{
		"booking_id": booking.ID,
		"user_id":    booking.UserID,
		"amount":     string(rune(booking.TotalAmount)),
	}

	// Add method-specific metadata - in real implementation, this would come from user input
	switch method {
	case models.PaymentMethodCreditCard:
		metadata["card_number"] = "1234-5678-9012-3456"
		metadata["cvv"] = "123"
		metadata["expiry"] = "12/25"
	case models.PaymentMethodDebitCard:
		metadata["card_number"] = "1234-5678-9012-3456"
		metadata["pin"] = "1234"
	case models.PaymentMethodUPI:
		metadata["upi_id"] = "user@paytm"
	case models.PaymentMethodNetBanking:
		metadata["bank_code"] = "HDFC"
		metadata["account_number"] = "1234567890"
	case models.PaymentMethodWallet:
		metadata["wallet_id"] = "wallet123"
	}

	return metadata
}
