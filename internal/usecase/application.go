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

// CreateArticle creates a new blog article in draft state
func (a *Application) CreateArticle(ctx context.Context, title, content, authorID string) (*entity.Article, error) {
	a.logger.DebugContext(ctx, "Creating new article", "title", title, "authorID", authorID)

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
		return nil, fmt.Errorf("failed to create article: %w", err)
	}

	return article, nil
}

// GetArticle retrieves a blog article by ID
func (a *Application) GetArticle(ctx context.Context, id string) (*entity.Article, error) {
	a.logger.DebugContext(ctx, "Getting article", "id", id)

	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	article, err := a.blogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get article: %w", err)
	}

	return article, nil
}

// ListArticles retrieves blog articles based on filter criteria
func (a *Application) ListArticles(ctx context.Context, filter BlogFilter) ([]*entity.Article, error) {
	a.logger.DebugContext(ctx, "Listing articles", "filter", filter)

	articles, err := a.blogRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list articles: %w", err)
	}

	return articles, nil
}

// UpdateArticle updates an existing blog article's content
func (a *Application) UpdateArticle(ctx context.Context, id, title, content string) (*entity.Article, error) {
	a.logger.DebugContext(ctx, "Updating article", "id", id, "title", title)

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
		return nil, fmt.Errorf("failed to get article for update: %w", err)
	}

	article.Title = title
	article.Content = content
	article.UpdatedAt = time.Now()

	if err := a.blogRepo.Update(ctx, article); err != nil {
		return nil, fmt.Errorf("failed to update article: %w", err)
	}

	return article, nil
}

// PublishArticle publishes a blog article immediately
func (a *Application) PublishArticle(ctx context.Context, id string) (*entity.Article, error) {
	a.logger.DebugContext(ctx, "Publishing article", "id", id)

	return a.scheduleArticle(ctx, id, time.Now())
}

// ScheduleArticle schedules a blog article for future publication
func (a *Application) ScheduleArticle(ctx context.Context, id string, publishAt time.Time) (*entity.Article, error) {
	a.logger.DebugContext(ctx, "Scheduling article", "id", id, "publishAt", publishAt)

	return a.scheduleArticle(ctx, id, publishAt)
}

// scheduleArticle internal helper for publishing/scheduling
func (a *Application) scheduleArticle(ctx context.Context, id string, publishAt time.Time) (*entity.Article, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	article, err := a.blogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get article for publishing: %w", err)
	}

	article.PublishedAt = &publishAt
	article.UpdatedAt = time.Now()

	if err := a.blogRepo.Update(ctx, article); err != nil {
		return nil, fmt.Errorf("failed to publish article: %w", err)
	}

	return article, nil
}

// UnpublishArticle unpublishes a blog article (makes it draft)
func (a *Application) UnpublishArticle(ctx context.Context, id string) (*entity.Article, error) {
	a.logger.DebugContext(ctx, "Unpublishing article", "id", id)

	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	article, err := a.blogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get article for unpublishing: %w", err)
	}

	article.PublishedAt = nil
	article.UpdatedAt = time.Now()

	if err := a.blogRepo.Update(ctx, article); err != nil {
		return nil, fmt.Errorf("failed to unpublish article: %w", err)
	}

	return article, nil
}

// DeleteArticle deletes a blog article
func (a *Application) DeleteArticle(ctx context.Context, id string) error {
	a.logger.DebugContext(ctx, "Deleting article", "id", id)

	if id == "" {
		return fmt.Errorf("id is required")
	}

	if err := a.blogRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete article: %w", err)
	}

	return nil
}
