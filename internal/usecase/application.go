package usecase

import (
	"context"
	"fmt"
	"go-clean-arch/internal/entity"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Application represents the main application service that orchestrates
// business logic and coordinates between different use cases.
// It follows Clean Architecture principles by depending only on
// interfaces and entities, not on external frameworks or infrastructure.
type Application struct {
	logger   *slog.Logger
	blogRepo BlogRepository
}

// NewApplicationParams holds the dependencies needed to create a new Application.
type NewApplicationParams struct {
	Logger   *slog.Logger
	BlogRepo BlogRepository
}

// NewApplication creates a new Application instance with the provided dependencies.
// This follows dependency injection patterns to maintain loose coupling.
func NewApplication(params NewApplicationParams) (*Application, error) {
	if params.BlogRepo == nil {
		return nil, fmt.Errorf("BlogRepo is required")
	}

	return &Application{
		logger:   params.Logger.With("component", "application"),
		blogRepo: params.BlogRepo,
	}, nil
}

// CreatePost creates a new blog post in draft state
func (a *Application) CreatePost(ctx context.Context, title, content, authorID string) (*entity.Article, error) {
	a.logger.DebugContext(ctx, "Creating new post", "title", title, "authorID", authorID)

	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if content == "" {
		return nil, fmt.Errorf("content is required")
	}
	if authorID == "" {
		return nil, fmt.Errorf("authorID is required")
	}

	now := time.Now()
	article := &entity.Article{
		ID:          uuid.New().String(),
		Title:       title,
		Content:     content,
		AuthorID:    authorID,
		CreatedAt:   now,
		UpdatedAt:   now,
		PublishedAt: nil, // Draft state
	}

	if err := a.blogRepo.Create(ctx, article); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return article, nil
}

// GetPost retrieves a blog post by ID
func (a *Application) GetPost(ctx context.Context, id string) (*entity.Article, error) {
	a.logger.DebugContext(ctx, "Getting post", "id", id)

	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	article, err := a.blogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return article, nil
}

// ListPosts retrieves blog posts based on filter criteria
func (a *Application) ListPosts(ctx context.Context, filter BlogFilter) ([]*entity.Article, error) {
	a.logger.DebugContext(ctx, "Listing posts", "filter", filter)

	articles, err := a.blogRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}

	return articles, nil
}

// UpdatePost updates an existing blog post's content
func (a *Application) UpdatePost(ctx context.Context, id, title, content string) (*entity.Article, error) {
	a.logger.DebugContext(ctx, "Updating post", "id", id, "title", title)

	if id == "" {
		return nil, fmt.Errorf("id is required")
	}
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if content == "" {
		return nil, fmt.Errorf("content is required")
	}

	article, err := a.blogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post for update: %w", err)
	}

	article.Title = title
	article.Content = content
	article.UpdatedAt = time.Now()

	if err := a.blogRepo.Update(ctx, article); err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return article, nil
}

// PublishPost publishes a blog post immediately
func (a *Application) PublishPost(ctx context.Context, id string) (*entity.Article, error) {
	a.logger.DebugContext(ctx, "Publishing post", "id", id)

	return a.schedulePost(ctx, id, time.Now())
}

// SchedulePost schedules a blog post for future publication
func (a *Application) SchedulePost(ctx context.Context, id string, publishAt time.Time) (*entity.Article, error) {
	a.logger.DebugContext(ctx, "Scheduling post", "id", id, "publishAt", publishAt)

	return a.schedulePost(ctx, id, publishAt)
}

// schedulePost internal helper for publishing/scheduling
func (a *Application) schedulePost(ctx context.Context, id string, publishAt time.Time) (*entity.Article, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	article, err := a.blogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post for publishing: %w", err)
	}

	article.PublishedAt = &publishAt
	article.UpdatedAt = time.Now()

	if err := a.blogRepo.Update(ctx, article); err != nil {
		return nil, fmt.Errorf("failed to publish post: %w", err)
	}

	return article, nil
}

// UnpublishPost unpublishes a blog post (makes it draft)
func (a *Application) UnpublishPost(ctx context.Context, id string) (*entity.Article, error) {
	a.logger.DebugContext(ctx, "Unpublishing post", "id", id)

	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	article, err := a.blogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post for unpublishing: %w", err)
	}

	article.PublishedAt = nil
	article.UpdatedAt = time.Now()

	if err := a.blogRepo.Update(ctx, article); err != nil {
		return nil, fmt.Errorf("failed to unpublish post: %w", err)
	}

	return article, nil
}

// DeletePost deletes a blog post
func (a *Application) DeletePost(ctx context.Context, id string) error {
	a.logger.DebugContext(ctx, "Deleting post", "id", id)

	if id == "" {
		return fmt.Errorf("id is required")
	}

	if err := a.blogRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}
