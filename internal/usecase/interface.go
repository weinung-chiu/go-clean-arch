package usecase

import (
	"context"
	"go-clean-arch/internal/entity"
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
