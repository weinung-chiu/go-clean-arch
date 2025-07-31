package adapter

// This package contains adapter implementations that connect the application
// to external systems and frameworks. Adapters implement interfaces defined
// in the use case layer and handle the technical details of external integrations.
//
// Common adapter types include:
// - Repository implementations (database, file system, cache)
// - External service clients (HTTP APIs, message queues)
// - Authentication and authorization services
// - Notification services (email, SMS, push)
//
// Adapter implementations should:
// - Implement interfaces defined in the use case layer
// - Handle technical concerns (serialization, networking, persistence)
// - Convert between domain entities and external formats
// - Manage connections and configuration for external systems
//
// Example structure:
//
// type DatabaseRepository struct {
//     db *sql.DB
//     logger *slog.Logger
// }
//
// func (r *DatabaseRepository) Create(ctx context.Context, entity *entity.Entity) error {
//     // Implementation details for database operations
// }

// Remove this file when adding real adapter implementations
