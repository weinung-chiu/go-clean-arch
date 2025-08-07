package usecase

import (
	"context"
	"errors"
	"go-clean-arch/internal/entity"
	"log/slog"
	"os"
	"testing"
	"time"
)

// MockBlogRepository is a test implementation of BlogRepository
type MockBlogRepository struct {
	articles map[string]*entity.Article
	createFn func(ctx context.Context, article *entity.Article) error
	getFn    func(ctx context.Context, id string) (*entity.Article, error)
	updateFn func(ctx context.Context, article *entity.Article) error
	deleteFn func(ctx context.Context, id string) error
	listFn   func(ctx context.Context, filter BlogFilter) ([]*entity.Article, error)
}

func NewMockBlogRepository() *MockBlogRepository {
	return &MockBlogRepository{
		articles: make(map[string]*entity.Article),
	}
}

func (m *MockBlogRepository) Create(ctx context.Context, article *entity.Article) error {
	if m.createFn != nil {
		return m.createFn(ctx, article)
	}
	m.articles[article.ID] = article
	return nil
}

func (m *MockBlogRepository) GetByID(ctx context.Context, id string) (*entity.Article, error) {
	if m.getFn != nil {
		return m.getFn(ctx, id)
	}
	article, exists := m.articles[id]
	if !exists {
		return nil, errors.New("article not found")
	}
	return article, nil
}

func (m *MockBlogRepository) Update(ctx context.Context, article *entity.Article) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, article)
	}
	if _, exists := m.articles[article.ID]; !exists {
		return errors.New("article not found")
	}
	m.articles[article.ID] = article
	return nil
}

func (m *MockBlogRepository) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	if _, exists := m.articles[id]; !exists {
		return errors.New("article not found")
	}
	delete(m.articles, id)
	return nil
}

func (m *MockBlogRepository) List(ctx context.Context, filter BlogFilter) ([]*entity.Article, error) {
	if m.listFn != nil {
		return m.listFn(ctx, filter)
	}

	var articles []*entity.Article
	for _, article := range m.articles {
		if m.matchesFilter(article, filter) {
			articles = append(articles, article)
		}
	}
	return articles, nil
}

func (m *MockBlogRepository) matchesFilter(article *entity.Article, filter BlogFilter) bool {
	if filter.ID != nil && article.ID != *filter.ID {
		return false
	}
	if filter.AuthorID != nil && article.AuthorID != *filter.AuthorID {
		return false
	}
	if filter.PublishedOnly != nil && *filter.PublishedOnly && !article.IsPublished() {
		return false
	}
	if filter.DraftsOnly != nil && *filter.DraftsOnly && !article.IsDraft() {
		return false
	}
	if filter.ScheduledOnly != nil && *filter.ScheduledOnly && !article.IsScheduled() {
		return false
	}
	return true
}

func setupTestApplication() (*Application, *MockBlogRepository) {
	mockRepo := NewMockBlogRepository()
	mockAuthService := &MockAuthService{}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	app, _ := NewApplication(NewApplicationParams{
		Logger:      logger,
		BlogRepo:    mockRepo,
		AuthService: mockAuthService,
	})

	return app, mockRepo
}

// Mock implementations for testing
type MockUserRepository struct{}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error { return nil }
func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	return nil, nil
}
func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	return nil, nil
}
func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error { return nil }

type MockAuthService struct{}

func (m *MockAuthService) LoginWithPassword(username, password string) (*AuthResult, error) {
	return nil, nil
}
func (m *MockAuthService) RegisterWithPassword(username, password string) (*AuthResult, error) {
	return nil, nil
}
func (m *MockAuthService) ValidateToken(token string) (*entity.User, error) { return nil, nil }

func TestNewApplication(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	mockRepo := NewMockBlogRepository()

	t.Run("successful creation", func(t *testing.T) {
		app, err := NewApplication(NewApplicationParams{
			Logger:      logger,
			BlogRepo:    mockRepo,
			AuthService: &MockAuthService{},
		})

		if err != nil {
			t.Errorf("NewApplication() error = %v, want nil", err)
		}
		if app == nil {
			t.Error("NewApplication() returned nil application")
		}
	})

	t.Run("missing blog repository", func(t *testing.T) {
		_, err := NewApplication(NewApplicationParams{
			Logger:      logger,
			BlogRepo:    nil,
			AuthService: &MockAuthService{},
		})

		if err == nil {
			t.Error("NewApplication() error = nil, want error for missing BlogRepo")
		}
	})

	t.Run("missing auth service", func(t *testing.T) {
		_, err := NewApplication(NewApplicationParams{
			Logger:      logger,
			BlogRepo:    mockRepo,
			AuthService: nil,
		})

		if err == nil {
			t.Error("NewApplication() error = nil, want error for missing AuthService")
		}
	})
}

func TestApplication_CreateArticle(t *testing.T) {
	app, _ := setupTestApplication()
	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		article, err := app.CreateArticle(ctx, "Test Title", "Test Content", "author-123")

		if err != nil {
			t.Errorf("CreateArticle() error = %v, want nil", err)
		}
		if article == nil {
			t.Fatal("CreateArticle() returned nil article")
		}
		if article.Title != "Test Title" {
			t.Errorf("CreateArticle() title = %v, want 'Test Title'", article.Title)
		}
		if article.Content != "Test Content" {
			t.Errorf("CreateArticle() content = %v, want 'Test Content'", article.Content)
		}
		if article.AuthorID != "author-123" {
			t.Errorf("CreateArticle() authorID = %v, want 'author-123'", article.AuthorID)
		}
		if !article.IsDraft() {
			t.Error("CreateArticle() should create article in draft state")
		}
	})

	t.Run("missing title", func(t *testing.T) {
		_, err := app.CreateArticle(ctx, "", "Test Content", "author-123")
		if err == nil {
			t.Error("CreateArticle() error = nil, want error for missing title")
		}
	})

	t.Run("missing content", func(t *testing.T) {
		_, err := app.CreateArticle(ctx, "Test Title", "", "author-123")
		if err == nil {
			t.Error("CreateArticle() error = nil, want error for missing content")
		}
	})

	t.Run("missing author ID", func(t *testing.T) {
		_, err := app.CreateArticle(ctx, "Test Title", "Test Content", "")
		if err == nil {
			t.Error("CreateArticle() error = nil, want error for missing author ID")
		}
	})
}

func TestApplication_GetArticle(t *testing.T) {
	app, mockRepo := setupTestApplication()
	ctx := context.Background()

	// Setup test data
	testArticle := &entity.Article{
		ID:        "test-id",
		Title:     "Test Article",
		Content:   "Test content",
		AuthorID:  "test-author",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.articles["test-id"] = testArticle

	t.Run("successful retrieval", func(t *testing.T) {
		article, err := app.GetArticle(ctx, "test-id")

		if err != nil {
			t.Errorf("GetArticle() error = %v, want nil", err)
		}
		if article == nil {
			t.Fatal("GetArticle() returned nil article")
		}
		if article.ID != "test-id" {
			t.Errorf("GetArticle() ID = %v, want 'test-id'", article.ID)
		}
	})

	t.Run("missing ID", func(t *testing.T) {
		_, err := app.GetArticle(ctx, "")
		if err == nil {
			t.Error("GetArticle() error = nil, want error for missing ID")
		}
	})

	t.Run("article not found", func(t *testing.T) {
		_, err := app.GetArticle(ctx, "nonexistent-id")
		if err == nil {
			t.Error("GetArticle() error = nil, want error for nonexistent article")
		}
	})
}

func TestApplication_UpdateArticle(t *testing.T) {
	app, mockRepo := setupTestApplication()
	ctx := context.Background()

	// Setup test data
	testArticle := &entity.Article{
		ID:        "test-id",
		Title:     "Original Title",
		Content:   "Original content",
		AuthorID:  "test-author",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.articles["test-id"] = testArticle

	t.Run("successful update", func(t *testing.T) {
		updatedArticle, err := app.UpdateArticle(ctx, "test-id", "Updated Title", "Updated content")

		if err != nil {
			t.Errorf("UpdateArticle() error = %v, want nil", err)
		}
		if updatedArticle == nil {
			t.Fatal("UpdateArticle() returned nil article")
		}
		if updatedArticle.Title != "Updated Title" {
			t.Errorf("UpdateArticle() title = %v, want 'Updated Title'", updatedArticle.Title)
		}
		if updatedArticle.Content != "Updated content" {
			t.Errorf("UpdateArticle() content = %v, want 'Updated content'", updatedArticle.Content)
		}
	})

	t.Run("missing parameters", func(t *testing.T) {
		tests := []struct {
			name    string
			id      string
			title   string
			content string
		}{
			{"missing ID", "", "Title", "Content"},
			{"missing title", "test-id", "", "Content"},
			{"missing content", "test-id", "Title", ""},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := app.UpdateArticle(ctx, tt.id, tt.title, tt.content)
				if err == nil {
					t.Errorf("UpdateArticle() error = nil, want error for %s", tt.name)
				}
			})
		}
	})
}

func TestApplication_PublishArticle(t *testing.T) {
	app, mockRepo := setupTestApplication()
	ctx := context.Background()

	// Setup test data
	testArticle := &entity.Article{
		ID:        "test-id",
		Title:     "Test Article",
		Content:   "Test content",
		AuthorID:  "test-author",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.articles["test-id"] = testArticle

	t.Run("successful publish", func(t *testing.T) {
		publishedArticle, err := app.PublishArticle(ctx, "test-id")

		if err != nil {
			t.Errorf("PublishArticle() error = %v, want nil", err)
		}
		if publishedArticle == nil {
			t.Fatal("PublishArticle() returned nil article")
		}
		if !publishedArticle.IsPublished() {
			t.Error("PublishArticle() should mark article as published")
		}
	})

	t.Run("missing ID", func(t *testing.T) {
		_, err := app.PublishArticle(ctx, "")
		if err == nil {
			t.Error("PublishArticle() error = nil, want error for missing ID")
		}
	})
}

func TestApplication_UnpublishArticle(t *testing.T) {
	app, mockRepo := setupTestApplication()
	ctx := context.Background()

	// Setup published test data
	now := time.Now()
	testArticle := &entity.Article{
		ID:          "test-id",
		Title:       "Test Article",
		Content:     "Test content",
		AuthorID:    "test-author",
		CreatedAt:   now,
		UpdatedAt:   now,
		PublishedAt: &now,
	}
	mockRepo.articles["test-id"] = testArticle

	t.Run("successful unpublish", func(t *testing.T) {
		unpublishedArticle, err := app.UnpublishArticle(ctx, "test-id")

		if err != nil {
			t.Errorf("UnpublishArticle() error = %v, want nil", err)
		}
		if unpublishedArticle == nil {
			t.Fatal("UnpublishArticle() returned nil article")
		}
		if !unpublishedArticle.IsDraft() {
			t.Error("UnpublishArticle() should mark article as draft")
		}
	})

	t.Run("missing ID", func(t *testing.T) {
		_, err := app.UnpublishArticle(ctx, "")
		if err == nil {
			t.Error("UnpublishArticle() error = nil, want error for missing ID")
		}
	})
}

func TestApplication_DeleteArticle(t *testing.T) {
	app, mockRepo := setupTestApplication()
	ctx := context.Background()

	// Setup test data
	testArticle := &entity.Article{
		ID:       "test-id",
		Title:    "Test Article",
		Content:  "Test content",
		AuthorID: "test-author",
	}
	mockRepo.articles["test-id"] = testArticle

	t.Run("successful delete", func(t *testing.T) {
		err := app.DeleteArticle(ctx, "test-id")

		if err != nil {
			t.Errorf("DeleteArticle() error = %v, want nil", err)
		}

		// Verify article is deleted
		_, exists := mockRepo.articles["test-id"]
		if exists {
			t.Error("DeleteArticle() should remove article from repository")
		}
	})

	t.Run("missing ID", func(t *testing.T) {
		err := app.DeleteArticle(ctx, "")
		if err == nil {
			t.Error("DeleteArticle() error = nil, want error for missing ID")
		}
	})
}

func TestApplication_ListArticles(t *testing.T) {
	app, mockRepo := setupTestApplication()
	ctx := context.Background()

	// Setup test data
	now := time.Now()
	draftArticle := &entity.Article{ID: "draft", Title: "Draft", AuthorID: "author1", CreatedAt: now, UpdatedAt: now}
	publishedArticle := &entity.Article{ID: "published", Title: "Published", AuthorID: "author1", CreatedAt: now, UpdatedAt: now, PublishedAt: &now}
	scheduledArticle := &entity.Article{ID: "scheduled", Title: "Scheduled", AuthorID: "author2", CreatedAt: now, UpdatedAt: now, PublishedAt: timePtr(now.Add(1 * time.Hour))}

	mockRepo.articles["draft"] = draftArticle
	mockRepo.articles["published"] = publishedArticle
	mockRepo.articles["scheduled"] = scheduledArticle

	t.Run("list all articles", func(t *testing.T) {
		articles, err := app.ListArticles(ctx, BlogFilter{})

		if err != nil {
			t.Errorf("ListArticles() error = %v, want nil", err)
		}
		if len(articles) != 3 {
			t.Errorf("ListArticles() returned %d articles, want 3", len(articles))
		}
	})

	t.Run("filter by author", func(t *testing.T) {
		filter := BlogFilter{AuthorID: stringPtr("author1")}
		articles, err := app.ListArticles(ctx, filter)

		if err != nil {
			t.Errorf("ListArticles() error = %v, want nil", err)
		}
		if len(articles) != 2 {
			t.Errorf("ListArticles() returned %d articles, want 2", len(articles))
		}
	})

	t.Run("filter published only", func(t *testing.T) {
		filter := BlogFilter{PublishedOnly: boolPtr(true)}
		articles, err := app.ListArticles(ctx, filter)

		if err != nil {
			t.Errorf("ListArticles() error = %v, want nil", err)
		}
		if len(articles) != 1 {
			t.Errorf("ListArticles() returned %d articles, want 1", len(articles))
		}
		if articles[0].ID != "published" {
			t.Errorf("ListArticles() returned wrong article, want published article")
		}
	})
}

// Helper functions
func stringPtr(s string) *string     { return &s }
func boolPtr(b bool) *bool           { return &b }
func timePtr(t time.Time) *time.Time { return &t }
