package controllers

import (
	"bookmyshow-lld/internal/repositories"
	"bookmyshow-lld/internal/services"
	"bookmyshow-lld/internal/strategies"
	"sync"
)

// AppController manages application lifecycle and dependency injection
// This is the proper place for orchestration logic
type AppController struct {
	// Business Services
	userService    services.UserService
	movieService   services.MovieService
	theatreService services.TheatreService
	showService    services.ShowService
	bookingService services.BookingService
	paymentService services.PaymentService

	// Repository Layer - explicit dependencies for type safety
	userRepo    repositories.UserRepository
	movieRepo   repositories.MovieRepository
	theatreRepo repositories.TheatreRepository
	screenRepo  repositories.ScreenRepository
	showRepo    repositories.ShowRepository
	bookingRepo repositories.BookingRepository
	paymentRepo repositories.PaymentRepository

	// External Services Layer
	paymentGateway  services.PaymentGateway
	notificationSvc services.NotificationService
}

var (
	instance *AppController
	once     sync.Once
)

// GetAppController returns singleton instance using dependency injection
func GetAppController() *AppController {
	once.Do(func() {
		instance = &AppController{}
		instance.initializeApp()
	})
	return instance
}

// initializeApp sets up the entire application with proper dependency injection
func (ac *AppController) initializeApp() {
	// Step 1: Initialize Infrastructure Layer (Repositories)
	ac.initializeRepositories()

	// Step 2: Initialize External Services
	ac.initializeExternalServices()

	// Step 3: Initialize Business Services with Dependencies
	ac.initializeBusinessServices()
}

// initializeRepositories creates all repository instances - explicit and type-safe
func (ac *AppController) initializeRepositories() {
	ac.userRepo = repositories.NewMemoryUserRepository()
	ac.movieRepo = repositories.NewMemoryMovieRepository()
	ac.theatreRepo = repositories.NewMemoryTheatreRepository()
	ac.screenRepo = repositories.NewMemoryScreenRepository()
	ac.showRepo = repositories.NewMemoryShowRepository()
	ac.bookingRepo = repositories.NewMemoryBookingRepository()
	ac.paymentRepo = repositories.NewMemoryPaymentRepository()
}

// initializeExternalServices creates external service connections - explicit and type-safe
func (ac *AppController) initializeExternalServices() {
	ac.paymentGateway = strategies.NewPaymentGateway()
	ac.notificationSvc = services.NewNotificationService()
}

// initializeBusinessServices creates business services with proper dependencies
func (ac *AppController) initializeBusinessServices() {
	// Create business services with explicit dependencies - no type assertions needed
	ac.userService = services.NewUserService(ac.userRepo)
	ac.movieService = services.NewMovieService(ac.movieRepo)
	ac.theatreService = services.NewTheatreService(ac.theatreRepo, ac.screenRepo)
	ac.showService = services.NewShowService(ac.showRepo, ac.movieRepo, ac.theatreRepo, ac.screenRepo)
	ac.bookingService = services.NewBookingService(
		ac.bookingRepo,
		ac.showRepo,
		ac.screenRepo,
		ac.theatreRepo,
		ac.movieRepo,
		ac.paymentRepo,
		ac.notificationSvc,
	)
	ac.paymentService = services.NewPaymentService(
		ac.paymentRepo,
		ac.bookingRepo,
		ac.paymentGateway,
		ac.notificationSvc,
	)
}

// Business Service Getters - Clean interface for accessing services
func (ac *AppController) GetUserService() services.UserService {
	return ac.userService
}

func (ac *AppController) GetMovieService() services.MovieService {
	return ac.movieService
}

func (ac *AppController) GetTheatreService() services.TheatreService {
	return ac.theatreService
}

func (ac *AppController) GetShowService() services.ShowService {
	return ac.showService
}

func (ac *AppController) GetBookingService() services.BookingService {
	return ac.bookingService
}

func (ac *AppController) GetPaymentService() services.PaymentService {
	return ac.paymentService
}

// Application lifecycle management
func (ac *AppController) Shutdown() {
	// Cleanup operations:
	// - Close database connections
	// - Stop background workers
	// - Release resources
	// - Graceful shutdown of services
}

// Health check for monitoring
func (ac *AppController) HealthCheck() map[string]string {
	return map[string]string{
		"status":       "healthy",
		"services":     "6 services running",
		"repositories": "7 repositories connected",
	}
}
