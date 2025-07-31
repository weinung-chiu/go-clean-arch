package adapter

import (
	"context"
	"fmt"
	"go-clean-arch/internal/entity"
	"go-clean-arch/internal/usecase"
	"sync"
)

// MemoryBlogRepository provides an in-memory implementation of BlogRepository
// for demonstration and testing purposes. It stores articles in memory using
// a map with thread-safe operations via mutex.
type MemoryBlogRepository struct {
	articles map[string]*entity.Article
	mu       sync.RWMutex
}

// NewMemoryBlogRepository creates a new in-memory blog repository
func NewMemoryBlogRepository() *MemoryBlogRepository {
	return &MemoryBlogRepository{
		articles: make(map[string]*entity.Article),
	}
}

// Create stores a new article in memory
func (r *MemoryBlogRepository) Create(ctx context.Context, article *entity.Article) error {
	if article == nil {
		return fmt.Errorf("article cannot be nil")
	}
	if article.ID == "" {
		return fmt.Errorf("article ID cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if article already exists
	if _, exists := r.articles[article.ID]; exists {
		return fmt.Errorf("article with ID %s already exists", article.ID)
	}

	// Create a copy to avoid external modifications
	articleCopy := *article
	r.articles[article.ID] = &articleCopy

	return nil
}

// GetByID retrieves an article by its ID
func (r *MemoryBlogRepository) GetByID(ctx context.Context, id string) (*entity.Article, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	article, exists := r.articles[id]
	if !exists {
		return nil, fmt.Errorf("article with ID %s not found", id)
	}

	// Return a copy to avoid external modifications
	articleCopy := *article
	return &articleCopy, nil
}

// Update modifies an existing article
func (r *MemoryBlogRepository) Update(ctx context.Context, article *entity.Article) error {
	if article == nil {
		return fmt.Errorf("article cannot be nil")
	}
	if article.ID == "" {
		return fmt.Errorf("article ID cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if article exists
	if _, exists := r.articles[article.ID]; !exists {
		return fmt.Errorf("article with ID %s not found", article.ID)
	}

	// Update with a copy
	articleCopy := *article
	r.articles[article.ID] = &articleCopy

	return nil
}

// Delete removes an article from memory
func (r *MemoryBlogRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if article exists
	if _, exists := r.articles[id]; !exists {
		return fmt.Errorf("article with ID %s not found", id)
	}

	delete(r.articles, id)
	return nil
}

// List retrieves articles based on filter criteria
func (r *MemoryBlogRepository) List(ctx context.Context, filter usecase.BlogFilter) ([]*entity.Article, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entity.Article

	for _, article := range r.articles {
		// Apply filters
		if filter.ID != nil && *filter.ID != article.ID {
			continue
		}

		if filter.AuthorID != nil && *filter.AuthorID != article.AuthorID {
			continue
		}

		// Publication status filters
		isPublished := article.IsPublished()
		isScheduled := article.IsScheduled()
		isDraft := article.IsDraft()

		if filter.PublishedOnly != nil && *filter.PublishedOnly && !isPublished {
			continue
		}

		if filter.DraftsOnly != nil && *filter.DraftsOnly && !isDraft {
			continue
		}

		if filter.ScheduledOnly != nil && *filter.ScheduledOnly && !isScheduled {
			continue
		}

		// Create a copy for the result
		articleCopy := *article
		result = append(result, &articleCopy)
	}

	return result, nil
}
