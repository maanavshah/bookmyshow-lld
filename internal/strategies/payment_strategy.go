package strategies

import (
	"bookmyshow-lld/internal/models"
	"bookmyshow-lld/internal/services"
	"fmt"
	"math/rand"
	"time"
)

// PaymentStrategy defines the strategy interface for payment processing - demonstrates Strategy Pattern
type PaymentStrategy interface {
	ProcessPayment(amount float64, metadata map[string]string) (*services.PaymentResult, error)
	ValidatePayment(metadata map[string]string) error
	GetPaymentMethod() models.PaymentMethod
}

// PaymentGatewayImpl implements the PaymentGateway interface using strategies
type PaymentGatewayImpl struct {
	strategies map[models.PaymentMethod]PaymentStrategy
}

// NewPaymentGateway creates a new payment gateway with all strategies - demonstrates Strategy Pattern
func NewPaymentGateway() *PaymentGatewayImpl {
	gateway := &PaymentGatewayImpl{
		strategies: make(map[models.PaymentMethod]PaymentStrategy),
	}

	// Register all payment strategies - demonstrates Strategy Pattern
	gateway.RegisterStrategy(&CreditCardStrategy{})
	gateway.RegisterStrategy(&DebitCardStrategy{})
	gateway.RegisterStrategy(&UPIStrategy{})
	gateway.RegisterStrategy(&NetBankingStrategy{})
	gateway.RegisterStrategy(&WalletStrategy{})

	return gateway
}

// RegisterStrategy registers a payment strategy
func (pg *PaymentGatewayImpl) RegisterStrategy(strategy PaymentStrategy) {
	pg.strategies[strategy.GetPaymentMethod()] = strategy
}

// ProcessPayment processes payment using the appropriate strategy - demonstrates Strategy Pattern
func (pg *PaymentGatewayImpl) ProcessPayment(amount float64, method models.PaymentMethod, metadata map[string]string) (*services.PaymentResult, error) {
	strategy, exists := pg.strategies[method]
	if !exists {
		return nil, fmt.Errorf("payment method %s not supported", method)
	}

	return strategy.ProcessPayment(amount, metadata)
}

// CreditCardStrategy implements payment processing for credit cards - demonstrates Concrete Strategy
type CreditCardStrategy struct{}

func (ccs *CreditCardStrategy) ProcessPayment(amount float64, metadata map[string]string) (*services.PaymentResult, error) {
	if err := ccs.ValidatePayment(metadata); err != nil {
		return &services.PaymentResult{
			Success:      false,
			ErrorMessage: err.Error(),
		}, err
	}

	// Mock payment processing - 90% success rate
	success := rand.Float32() > 0.1

	if success {
		return &services.PaymentResult{
			Success:       true,
			TransactionID: fmt.Sprintf("CC_%d", time.Now().Unix()),
			Response:      "Payment processed successfully via Credit Card",
		}, nil
	}

	return &services.PaymentResult{
		Success:      false,
		ErrorMessage: "Credit card payment failed",
	}, models.ErrPaymentProcessingFail
}

func (ccs *CreditCardStrategy) ValidatePayment(metadata map[string]string) error {
	if metadata["card_number"] == "" || metadata["cvv"] == "" || metadata["expiry"] == "" {
		return fmt.Errorf("missing required credit card details")
	}
	return nil
}

func (ccs *CreditCardStrategy) GetPaymentMethod() models.PaymentMethod {
	return models.PaymentMethodCreditCard
}

// DebitCardStrategy implements payment processing for debit cards - demonstrates Concrete Strategy
type DebitCardStrategy struct{}

func (dcs *DebitCardStrategy) ProcessPayment(amount float64, metadata map[string]string) (*services.PaymentResult, error) {
	if err := dcs.ValidatePayment(metadata); err != nil {
		return &services.PaymentResult{
			Success:      false,
			ErrorMessage: err.Error(),
		}, err
	}

	// Mock payment processing - 85% success rate
	success := rand.Float32() > 0.15

	if success {
		return &services.PaymentResult{
			Success:       true,
			TransactionID: fmt.Sprintf("DC_%d", time.Now().Unix()),
			Response:      "Payment processed successfully via Debit Card",
		}, nil
	}

	return &services.PaymentResult{
		Success:      false,
		ErrorMessage: "Debit card payment failed",
	}, models.ErrPaymentProcessingFail
}

func (dcs *DebitCardStrategy) ValidatePayment(metadata map[string]string) error {
	if metadata["card_number"] == "" || metadata["pin"] == "" {
		return fmt.Errorf("missing required debit card details")
	}
	return nil
}

func (dcs *DebitCardStrategy) GetPaymentMethod() models.PaymentMethod {
	return models.PaymentMethodDebitCard
}

// UPIStrategy implements payment processing for UPI - demonstrates Concrete Strategy
type UPIStrategy struct{}

func (upi *UPIStrategy) ProcessPayment(amount float64, metadata map[string]string) (*services.PaymentResult, error) {
	if err := upi.ValidatePayment(metadata); err != nil {
		return &services.PaymentResult{
			Success:      false,
			ErrorMessage: err.Error(),
		}, err
	}

	// Mock payment processing - 95% success rate (UPI is most reliable)
	success := rand.Float32() > 0.05

	if success {
		return &services.PaymentResult{
			Success:       true,
			TransactionID: fmt.Sprintf("UPI_%d", time.Now().Unix()),
			Response:      "Payment processed successfully via UPI",
		}, nil
	}

	return &services.PaymentResult{
		Success:      false,
		ErrorMessage: "UPI payment failed",
	}, models.ErrPaymentProcessingFail
}

func (upi *UPIStrategy) ValidatePayment(metadata map[string]string) error {
	if metadata["upi_id"] == "" {
		return fmt.Errorf("missing UPI ID")
	}
	return nil
}

func (upi *UPIStrategy) GetPaymentMethod() models.PaymentMethod {
	return models.PaymentMethodUPI
}

// NetBankingStrategy implements payment processing for net banking - demonstrates Concrete Strategy
type NetBankingStrategy struct{}

func (nb *NetBankingStrategy) ProcessPayment(amount float64, metadata map[string]string) (*services.PaymentResult, error) {
	if err := nb.ValidatePayment(metadata); err != nil {
		return &services.PaymentResult{
			Success:      false,
			ErrorMessage: err.Error(),
		}, err
	}

	// Mock payment processing - 92% success rate
	success := rand.Float32() > 0.08

	if success {
		return &services.PaymentResult{
			Success:       true,
			TransactionID: fmt.Sprintf("NB_%d", time.Now().Unix()),
			Response:      "Payment processed successfully via Net Banking",
		}, nil
	}

	return &services.PaymentResult{
		Success:      false,
		ErrorMessage: "Net banking payment failed",
	}, models.ErrPaymentProcessingFail
}

func (nb *NetBankingStrategy) ValidatePayment(metadata map[string]string) error {
	if metadata["bank_code"] == "" || metadata["account_number"] == "" {
		return fmt.Errorf("missing net banking details")
	}
	return nil
}

func (nb *NetBankingStrategy) GetPaymentMethod() models.PaymentMethod {
	return models.PaymentMethodNetBanking
}

// WalletStrategy implements payment processing for digital wallets - demonstrates Concrete Strategy
type WalletStrategy struct{}

func (ws *WalletStrategy) ProcessPayment(amount float64, metadata map[string]string) (*services.PaymentResult, error) {
	if err := ws.ValidatePayment(metadata); err != nil {
		return &services.PaymentResult{
			Success:      false,
			ErrorMessage: err.Error(),
		}, err
	}

	// Mock payment processing - 97% success rate (wallets are very reliable)
	success := rand.Float32() > 0.03

	if success {
		return &services.PaymentResult{
			Success:       true,
			TransactionID: fmt.Sprintf("WALLET_%d", time.Now().Unix()),
			Response:      "Payment processed successfully via Wallet",
		}, nil
	}

	return &services.PaymentResult{
		Success:      false,
		ErrorMessage: "Wallet payment failed",
	}, models.ErrPaymentProcessingFail
}

func (ws *WalletStrategy) ValidatePayment(metadata map[string]string) error {
	if metadata["wallet_id"] == "" {
		return fmt.Errorf("missing wallet ID")
	}
	return nil
}

func (ws *WalletStrategy) GetPaymentMethod() models.PaymentMethod {
	return models.PaymentMethodWallet
}
