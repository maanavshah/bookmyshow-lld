package models

import (
	"time"

	"github.com/google/uuid"
)

// PaymentMethod represents different payment methods
type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "CREDIT_CARD"
	PaymentMethodDebitCard  PaymentMethod = "DEBIT_CARD"
	PaymentMethodUPI        PaymentMethod = "UPI"
	PaymentMethodNetBanking PaymentMethod = "NET_BANKING"
	PaymentMethodWallet     PaymentMethod = "WALLET"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusSuccess   PaymentStatus = "SUCCESS"
	PaymentStatusFailed    PaymentStatus = "FAILED"
	PaymentStatusRefunded  PaymentStatus = "REFUNDED"
	PaymentStatusCancelled PaymentStatus = "CANCELLED"
)

// Payment represents a payment transaction
type Payment struct {
	ID              string        `json:"id"`
	BookingID       string        `json:"booking_id"`
	UserID          string        `json:"user_id"`
	Amount          float64       `json:"amount"`
	Method          PaymentMethod `json:"method"`
	Status          PaymentStatus `json:"status"`
	TransactionID   string        `json:"transaction_id,omitempty"`
	GatewayResponse string        `json:"gateway_response,omitempty"`
	FailureReason   string        `json:"failure_reason,omitempty"`
	RefundAmount    float64       `json:"refund_amount,omitempty"`
	RefundReason    string        `json:"refund_reason,omitempty"`
	ProcessedAt     *time.Time    `json:"processed_at,omitempty"`
	RefundedAt      *time.Time    `json:"refunded_at,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

// NewPayment creates a new payment
func NewPayment(bookingID, userID string, amount float64, method PaymentMethod) (*Payment, error) {
	if bookingID == "" || userID == "" || amount <= 0 {
		return nil, ErrInvalidPaymentData
	}

	return &Payment{
		ID:        uuid.New().String(),
		BookingID: bookingID,
		UserID:    userID,
		Amount:    amount,
		Method:    method,
		Status:    PaymentStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// MarkSuccess marks the payment as successful
func (p *Payment) MarkSuccess(transactionID, gatewayResponse string) {
	now := time.Now()
	p.Status = PaymentStatusSuccess
	p.TransactionID = transactionID
	p.GatewayResponse = gatewayResponse
	p.ProcessedAt = &now
	p.UpdatedAt = now
}

// MarkFailed marks the payment as failed
func (p *Payment) MarkFailed(failureReason string) {
	now := time.Now()
	p.Status = PaymentStatusFailed
	p.FailureReason = failureReason
	p.ProcessedAt = &now
	p.UpdatedAt = now
}

// MarkCancelled marks the payment as cancelled
func (p *Payment) MarkCancelled() {
	p.Status = PaymentStatusCancelled
	p.UpdatedAt = time.Now()
}

// ProcessRefund processes a refund for the payment
func (p *Payment) ProcessRefund(refundAmount float64, refundReason string) error {
	if p.Status != PaymentStatusSuccess {
		return ErrPaymentNotSuccessful
	}

	if refundAmount <= 0 || refundAmount > p.Amount {
		return ErrInvalidRefundAmount
	}

	now := time.Now()
	p.Status = PaymentStatusRefunded
	p.RefundAmount = refundAmount
	p.RefundReason = refundReason
	p.RefundedAt = &now
	p.UpdatedAt = now
	return nil
}

// IsSuccessful checks if payment was successful
func (p *Payment) IsSuccessful() bool {
	return p.Status == PaymentStatusSuccess
}

// IsPending checks if payment is pending
func (p *Payment) IsPending() bool {
	return p.Status == PaymentStatusPending
}

// IsFailed checks if payment failed
func (p *Payment) IsFailed() bool {
	return p.Status == PaymentStatusFailed
}

// IsRefunded checks if payment was refunded
func (p *Payment) IsRefunded() bool {
	return p.Status == PaymentStatusRefunded
}

// CanBeRefunded checks if payment can be refunded
func (p *Payment) CanBeRefunded() bool {
	return p.Status == PaymentStatusSuccess
}
