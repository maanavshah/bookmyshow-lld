package models

import "errors"

// User errors
var (
	ErrInvalidUserData = errors.New("invalid user data provided")
	ErrUserNotFound    = errors.New("user not found")
)

// Movie errors
var (
	ErrInvalidMovieData = errors.New("invalid movie data provided")
	ErrMovieNotFound    = errors.New("movie not found")
)

// Theatre errors
var (
	ErrInvalidTheatreData = errors.New("invalid theatre data provided")
	ErrTheatreNotFound    = errors.New("theatre not found")
)

// Screen errors
var (
	ErrScreenNotFound = errors.New("screen not found")
)

// Seat errors
var (
	ErrSeatNotFound      = errors.New("seat not found")
	ErrSeatNotAvailable  = errors.New("seat is not available")
	ErrSeatNotBlocked    = errors.New("seat is not blocked")
	ErrSeatAlreadyBooked = errors.New("seat is already booked")
)

// Show errors
var (
	ErrInvalidShowData = errors.New("invalid show data provided")
	ErrInvalidShowTime = errors.New("invalid show time")
	ErrShowNotFound    = errors.New("show not found")
	ErrShowNotBookable = errors.New("show is not available for booking")
)

// Booking errors
var (
	ErrInvalidBookingData      = errors.New("invalid booking data provided")
	ErrBookingNotFound         = errors.New("booking not found")
	ErrBookingNotPending       = errors.New("booking is not in pending status")
	ErrBookingExpired          = errors.New("booking has expired")
	ErrBookingAlreadyConfirmed = errors.New("booking is already confirmed")
	ErrBookingAlreadyCancelled = errors.New("booking is already cancelled")
	ErrInsufficientSeats       = errors.New("insufficient available seats")
)

// Payment errors
var (
	ErrInvalidPaymentData    = errors.New("invalid payment data provided")
	ErrPaymentNotFound       = errors.New("payment not found")
	ErrPaymentNotSuccessful  = errors.New("payment was not successful")
	ErrInvalidRefundAmount   = errors.New("invalid refund amount")
	ErrPaymentGatewayError   = errors.New("payment gateway error")
	ErrPaymentProcessingFail = errors.New("payment processing failed")
)

// Service errors
var (
	ErrServiceUnavailable = errors.New("service temporarily unavailable")
	ErrInternalError      = errors.New("internal server error")
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrConcurrencyIssue   = errors.New("concurrency conflict occurred")
)
