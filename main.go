package main

import (
	"bookmyshow-lld/internal/controllers"
	"bookmyshow-lld/internal/factories"
	"bookmyshow-lld/internal/models"
	"bookmyshow-lld/internal/services"
	"fmt"
	"log"
	"time"
)

func main() {
	fmt.Println("🎬 BookMyShow Low Level Design Learning Prototype")
	fmt.Println("==================================================")
	fmt.Println("🎯 Focus: Core Design Patterns & SOLID Principles")

	// Get application controller - demonstrates Singleton + Dependency Injection
	appController := controllers.GetAppController()
	defer appController.Shutdown()

	// Get services through controller - demonstrates clean architecture
	userService := appController.GetUserService()
	movieService := appController.GetMovieService()
	theatreService := appController.GetTheatreService()
	showService := appController.GetShowService()
	bookingService := appController.GetBookingService()
	paymentService := appController.GetPaymentService()

	// Run focused demo showcasing design patterns
	runApi(userService, movieService, theatreService, showService, bookingService, paymentService)
}

func runApi(
	userService services.UserService,
	movieService services.MovieService,
	theatreService services.TheatreService,
	showService services.ShowService,
	bookingService services.BookingService,
	paymentService services.PaymentService,
) {
	fmt.Println("\n📚 1. Repository Pattern - Creating Core Entities")

	// Create users - demonstrates Repository Pattern
	user1, err := userService.CreateUser("John Doe", "john@example.com", "+1234567890")
	if err != nil {
		log.Fatal("Failed to create user:", err)
	}
	fmt.Printf("✅ Created user: %s (Repository Pattern)\n", user1.Name)

	// Create movie - demonstrates Repository Pattern
	movie1, err := movieService.CreateMovie(
		"Avengers: Endgame",
		"Epic superhero finale",
		3*time.Hour,
		models.GenreAction,
		models.LanguageEnglish,
		8.4,
		time.Now().AddDate(0, -1, 0),
	)
	if err != nil {
		log.Fatal("Failed to create movie:", err)
	}
	fmt.Printf("✅ Created movie: %s (Repository Pattern)\n", movie1.Title)

	// Create theatre - demonstrates Repository Pattern
	theatre1, err := theatreService.CreateTheatre("PVR Cinemas", "Phoenix Mall", "Mumbai")
	if err != nil {
		log.Fatal("Failed to create theatre:", err)
	}
	fmt.Printf("✅ Created theatre: %s (Repository Pattern)\n", theatre1.Name)

	fmt.Println("\n🏭 2. Factory Pattern - Creating Complex Objects")

	// Create screen with seats using Factory Pattern
	screen1 := models.NewScreen("Screen 1", theatre1.ID)

	// Use SeatFactory to create seats - demonstrates Factory Pattern
	seatFactory := factories.NewSeatFactory()
	seats := seatFactory.CreateDefaultScreenSeats(100.0) // Base price: $100
	fmt.Printf("🏭 Factory Pattern: Created %d seats with different types and pricing\n", len(seats))

	for _, seat := range seats {
		screen1.AddSeat(seat)
	}

	// Add screen to theatre
	err = theatreService.AddScreen(theatre1.ID, screen1)
	if err != nil {
		log.Fatal("Failed to add screen:", err)
	}
	fmt.Printf("✅ Added screen with %d seats (Factory Pattern)\n", screen1.GetCapacity())

	fmt.Println("\n🎯 3. Business Logic & Validation")

	// Create show - demonstrates business rules and validation
	showTime1 := time.Now().Add(2 * time.Hour)
	show1, err := showService.CreateShow(movie1.ID, theatre1.ID, screen1.ID, showTime1, 100.0)
	if err != nil {
		log.Fatal("Failed to create show:", err)
	}
	fmt.Printf("✅ Created show with business validation (Start: %s)\n", showTime1.Format("15:04"))

	fmt.Println("\n🔒 4. Concurrency Control - Thread-Safe Booking")

	// Get available seats for booking
	availableSeats := screen1.GetAvailableSeats()
	if len(availableSeats) < 3 {
		log.Fatal("Not enough seats available")
	}
	seatIDs := []string{availableSeats[0].ID, availableSeats[1].ID, availableSeats[2].ID}

	// Book seats - demonstrates concurrency control
	booking1, err := bookingService.CreateBooking(user1.ID, show1.ID, seatIDs)
	if err != nil {
		log.Fatal("Failed to create booking:", err)
	}
	fmt.Printf("🔒 Thread-safe booking created: $%.2f (Concurrency Control)\n", booking1.TotalAmount)

	fmt.Println("\n🔄 5. Strategy Pattern - Payment Processing")

	// Process payment using Strategy Pattern - different payment methods
	payment1, err := paymentService.ProcessPayment(booking1.ID, models.PaymentMethodUPI)
	if err != nil {
		log.Printf("❌ Payment failed: %v", err)
	} else {
		fmt.Printf("🔄 Strategy Pattern: %s payment processed ($%.2f)\n", payment1.Method, payment1.Amount)

		if payment1.IsSuccessful() {
			// Confirm booking
			err = bookingService.ConfirmBooking(booking1.ID, payment1.ID)
			if err != nil {
				log.Printf("❌ Failed to confirm booking: %v", err)
			} else {
				fmt.Printf("✅ Booking confirmed! (Transaction: %s)\n", payment1.TransactionID)
			}
		}
	}

	fmt.Println("\n📢 6. Observer Pattern - Notifications")
	fmt.Println("📧 Notification sent via Observer Pattern (check logs above)")

	fmt.Println("\n🎯 7. Demonstrating Seat Pricing (Factory Pattern)")

	// Show seat pricing using Factory pattern
	seatInfo := seatFactory.GetSeatTypeInfo()
	fmt.Println("💺 Factory Pattern - Different seat types and pricing:")
	for seatType, info := range seatInfo {
		fmt.Printf("   %s: %s (%.1fx base price)\n",
			seatType, info.Description, info.Multiplier)
	}

	fmt.Println("\n🏗️ 8. Getting Aggregate Data")

	// Get detailed booking information - demonstrates aggregate construction
	bookingDetails, err := bookingService.GetBookingDetails(booking1.ID)
	if err != nil {
		log.Printf("Failed to get booking details: %v", err)
	} else {
		fmt.Printf("📋 Aggregate Construction:\n")
		fmt.Printf("   Movie: %s (%s)\n", bookingDetails.Movie.Title, bookingDetails.Movie.Language)
		fmt.Printf("   Theatre: %s, %s\n", bookingDetails.Theatre.Name, bookingDetails.Theatre.City)
		fmt.Printf("   Seats: ")
		for i, seat := range bookingDetails.Seats {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%s%d (%s-$%.0f)", seat.RowName, seat.Number, seat.Type, seat.Price)
		}
		fmt.Printf("\n   Total: $%.2f | Status: %s\n", bookingDetails.Booking.TotalAmount, bookingDetails.Booking.GetStatus())
	}

	fmt.Println("\n✨ Learning Demo Completed Successfully!")
	fmt.Println("\n🎓 Key Design Patterns Demonstrated:")
	fmt.Println("   🏭 Factory Pattern: SeatFactory creates different seat types with pricing")
	fmt.Println("   🔄 Strategy Pattern: Multiple payment methods (UPI, Credit Card, etc.)")
	fmt.Println("   🔒 Singleton Pattern: AppController manages application lifecycle")
	fmt.Println("   📦 Repository Pattern: Clean data access abstraction")
	fmt.Println("   📢 Observer Pattern: Notification system for booking events")
	fmt.Println("   🔐 Concurrency Control: Thread-safe seat booking with atomic operations")

	fmt.Println("\n⚡ SOLID Principles Applied:")
	fmt.Println("   S - Single Responsibility: Each class has one clear purpose")
	fmt.Println("   O - Open/Closed: Easy to add new payment methods or seat types")
	fmt.Println("   L - Liskov Substitution: All strategies implement PaymentStrategy")
	fmt.Println("   I - Interface Segregation: Small, focused interfaces")
	fmt.Println("   D - Dependency Inversion: Services depend on abstractions, not concrete classes")

	fmt.Println("\n🎯 Learning Objectives Achieved:")
	fmt.Println("   ✓ Clean Architecture with proper layering")
	fmt.Println("   ✓ Design Patterns implementation")
	fmt.Println("   ✓ SOLID principles application")
	fmt.Println("   ✓ Thread-safe operations")
	fmt.Println("   ✓ Business logic separation")
	fmt.Println("   ✓ Dependency injection")
}
