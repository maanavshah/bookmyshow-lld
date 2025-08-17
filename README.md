# BookMyShow Low Level Design (LLD) 🎬

A comprehensive Go implementation of BookMyShow system following SOLID principles, design patterns, and concurrent programming best practices.

## 🏗️ Architecture Overview

This project demonstrates a complete low-level design of a movie ticket booking system with proper separation of concerns, thread safety, and scalable architecture.

### Core Components

- **Models**: Domain entities with proper encapsulation
- **Interfaces**: Abstractions for dependency inversion
- **Services**: Business logic layer
- **Repositories**: Data access layer
- **Factories**: Object creation patterns
- **Strategies**: Different algorithm implementations

## 🎯 SOLID Principles Implementation

### 1. Single Responsibility Principle (SRP)
- Each class has one reason to change
- `User` handles user data, `Booking` handles booking logic
- Services are focused on specific business domains

### 2. Open/Closed Principle (OCP)
- Payment strategies can be extended without modifying existing code
- New seat types can be added via factory pattern
- Repository implementations can be swapped

### 3. Liskov Substitution Principle (LSP)
- All payment strategies are interchangeable
- Repository implementations can substitute each other
- Service interfaces ensure consistent behavior

### 4. Interface Segregation Principle (ISP)
- Focused interfaces for specific responsibilities
- `PaymentGateway`, `NotificationService` are separate concerns
- Clients depend only on interfaces they use

### 5. Dependency Inversion Principle (DIP)
- High-level services depend on abstractions
- Repositories are injected via interfaces
- Service manager orchestrates dependencies

## 🎨 Design Patterns

### 1. Factory Pattern
```go
// SeatFactory creates different types of seats with appropriate pricing
seatFactory := factories.NewSeatFactory()
vipSeat := seatFactory.CreateSeat("A", 1, models.SeatTypeVIP, 100.0)
```

### 2. Strategy Pattern
```go
// Different payment methods using strategy pattern
paymentGateway.ProcessPayment(amount, models.PaymentMethodUPI, metadata)
paymentGateway.ProcessPayment(amount, models.PaymentMethodCreditCard, metadata)
```

### 3. Singleton Pattern
```go
// Service manager ensures single instance
serviceManager := services.GetServiceManager()
```

### 4. Repository Pattern
```go
// Data access abstraction
type UserRepository interface {
    Create(user *models.User) error
    GetByID(id string) (*models.User, error)
    // ... other methods
}
```

### 5. Observer Pattern (via Notifications)
```go
// Notification service observes booking events
notificationSvc.SendBookingConfirmation(userID, bookingID)
```

## 🧵 Concurrency Handling

### Thread-Safe Operations
- **Seat Booking**: Atomic seat blocking with mutex
- **Repository Access**: RWMutex for concurrent read/write
- **Booking Expiry**: Safe status transitions

### Example Concurrency Control
```go
// Thread-safe seat blocking
func (s *Seat) Block() error {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    if s.Status != SeatStatusAvailable {
        return ErrSeatNotAvailable
    }
    
    s.Status = SeatStatusBlocked
    return nil
}
```

## 📊 Core Entities

### User Management
- User creation and profile management
- Email uniqueness validation
- Thread-safe operations

### Movie Management
- Multi-genre and multi-language support
- Release date validation
- Search functionality

### Theatre & Screen Management
- Multi-screen theatres
- Configurable seating arrangements
- Capacity management

### Show Management
- Schedule conflict detection
- Time-based availability
- Dynamic pricing support

### Booking System
- Atomic seat reservation
- Automatic expiry handling
- Concurrent booking prevention

### Payment Processing
- Multiple payment methods
- Transaction tracking
- Refund processing

## 🚀 Usage Example

```go
// Initialize system
serviceManager := services.GetServiceManager()

// Create user
user, _ := serviceManager.GetUserService().CreateUser(
    "John Doe", "john@email.com", "+1234567890"
)

// Create movie
movie, _ := serviceManager.GetMovieService().CreateMovie(
    "Avengers", "Action movie", 3*time.Hour,
    models.GenreAction, models.LanguageEnglish, 8.5,
    time.Now().AddDate(0, -1, 0)
)

// Book tickets
booking, _ := serviceManager.GetBookingService().CreateBooking(
    user.ID, show.ID, []string{seat1.ID, seat2.ID}
)

// Process payment
payment, _ := serviceManager.GetPaymentService().ProcessPayment(
    booking.ID, models.PaymentMethodUPI
)
```

## 🏃‍♂️ Running the Application

```bash
# Clone the repository
git clone <repository-url>
cd bookmyshow-lld

# Install dependencies
go mod tidy

# Run the application
go run main.go
```

## 📁 Project Structure

```
bookmyshow-lld/
├── main.go                  # Application entry point
├── internal/
│   ├── models/              # Domain entities
│   │   ├── user.go
│   │   ├── movie.go
│   │   ├── theatre.go
│   │   ├── screen.go
│   │   ├── seat.go
│   │   ├── show.go
│   │   ├── booking.go
│   │   ├── payment.go
│   │   └── errors.go
│   ├── interfaces/          # Abstractions
│   │   ├── repositories.go
│   │   └── services.go
│   ├── repositories/        # Data access layer
│   │   ├── memory_repository.go
│   │   └── show_booking_repositories.go
│   ├── services/           # Business logic
│   │   ├── basic_services.go
│   │   ├── booking_service.go
│   │   ├── payment_service.go
│   │   ├── notification_service.go
│   │   └── manager.go
│   ├── factories/          # Object creation
│   │   └── seat_factory.go
│   └── strategies/         # Algorithm implementations
│       └── payment_strategy.go
├── go.mod
└── README.md
```

## 🎯 Key Features

### ✅ Functional Features
- User registration and management
- Movie catalog with search
- Theatre and screen management
- Show scheduling with conflict detection
- Seat booking with different types
- Payment processing with multiple methods
- Booking confirmation and management
- Automatic booking expiry

### ✅ Non-Functional Features
- **Concurrency**: Thread-safe operations
- **Scalability**: Repository pattern for data layer
- **Maintainability**: Clean architecture with SOLID principles
- **Extensibility**: Factory and Strategy patterns
- **Reliability**: Proper error handling and validation
- **Performance**: Efficient data structures and algorithms

## 🔧 Technical Highlights

### Error Handling
- Custom error types for different domains
- Comprehensive error propagation
- Graceful failure handling

### Validation
- Input validation at service layer
- Business rule enforcement
- Data consistency checks

### Memory Management
- Efficient data structures
- Proper resource cleanup
- Memory-safe operations

### Testing Considerations
- Interface-based design enables easy mocking
- Dependency injection supports unit testing
- Clear separation of concerns

## 🚀 Future Enhancements

- Database integration with SQL/NoSQL
- REST API endpoints
- Real-time notifications
- Caching layer for performance
- Distributed system support
- Monitoring and logging
- Load balancing strategies

## 📝 Learning Outcomes

This project demonstrates:
- **Object-Oriented Programming**: Encapsulation, inheritance, polymorphism
- **Design Patterns**: Factory, Strategy, Singleton, Repository, Observer
- **SOLID Principles**: All five principles with practical examples
- **Concurrency**: Thread-safe programming in Go
- **Clean Architecture**: Separation of concerns and dependency management
- **Error Handling**: Robust error management strategies
- **Testing**: Design for testability

---

*Built with ❤️ using Go and best practices in software engineering* 
