package usecase

import (
	"context"
	"go-clean-arch/internal/entity"
	"time"
)

// Repository interfaces define contracts for data access in the use case layer.
// These interfaces follow Clean Architecture principles by depending only on
// the entity layer and providing abstractions for outer layers to implement.

// BlogFilter defines filter criteria for blog queries
type BlogFilter struct {
	ID            *string
	AuthorID      *string
	PublishedOnly *bool
	DraftsOnly    *bool
	ScheduledOnly *bool
}

// BlogRepository defines the contract for blog data access operations
type BlogRepository interface {
	Create(ctx context.Context, article *entity.Article) error
	GetByID(ctx context.Context, id string) (*entity.Article, error)
	Update(ctx context.Context, article *entity.Article) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter BlogFilter) ([]*entity.Article, error)
}

// AuthResult represents the result of successful authentication
type AuthResult struct {
	Token     string
	User      *entity.User
	ExpiresAt time.Time
}

// UserRepository defines the contract for user data access operations
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	GetByID(ctx context.Context, id string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
}

// AuthService defines the contract for authentication operations
type AuthService interface {
	LoginWithPassword(username, password string) (*AuthResult, error)
	RegisterWithPassword(username, password string) (*AuthResult, error)
	ValidateToken(token string) (*entity.User, error)
}
