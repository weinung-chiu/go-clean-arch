package adapter

import (
	"context"
	"fmt"
	"go-clean-arch/internal/entity"
	"go-clean-arch/internal/usecase"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ArticleModel represents the database model for articles
type ArticleModel struct {
	ID          string     `gorm:"primaryKey;type:varchar(255)"`
	Title       string     `gorm:"not null"`
	Content     string     `gorm:"not null"`
	AuthorID    string     `gorm:"not null;type:varchar(255);index"`
	CreatedAt   time.Time  `gorm:"not null;index"`
	UpdatedAt   time.Time  `gorm:"not null"`
	PublishedAt *time.Time `gorm:"index"`
}

// TableName specifies the table name for GORM
func (ArticleModel) TableName() string {
	return "articles"
}

// toEntity converts ArticleModel to domain entity
func (m *ArticleModel) toEntity() *entity.Article {
	return &entity.Article{
		ID:          m.ID,
		Title:       m.Title,
		Content:     m.Content,
		AuthorID:    m.AuthorID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		PublishedAt: m.PublishedAt,
	}
}

// fromEntity converts domain entity to ArticleModel
func (m *ArticleModel) fromEntity(article *entity.Article) {
	m.ID = article.ID
	m.Title = article.Title
	m.Content = article.Content
	m.AuthorID = article.AuthorID
	m.CreatedAt = article.CreatedAt
	m.UpdatedAt = article.UpdatedAt
	m.PublishedAt = article.PublishedAt
}

// PostgresBlogRepository provides a PostgreSQL implementation of BlogRepository using GORM
type PostgresBlogRepository struct {
	db *gorm.DB
}

// NewPostgresBlogRepository creates a new PostgreSQL blog repository with GORM
func NewPostgresBlogRepository(dsn string) (*PostgresBlogRepository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&ArticleModel{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &PostgresBlogRepository{db: db}, nil
}

// Create stores a new article in the database
func (r *PostgresBlogRepository) Create(ctx context.Context, article *entity.Article) error {
	if article == nil {
		return fmt.Errorf("article cannot be nil")
	}
	if article.ID == "" {
		return fmt.Errorf("article ID cannot be empty")
	}

	model := &ArticleModel{}
	model.fromEntity(article)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create article: %w", err)
	}

	return nil
}

// GetByID retrieves an article by its ID
func (r *PostgresBlogRepository) GetByID(ctx context.Context, id string) (*entity.Article, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	var model ArticleModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("article with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to get article: %w", err)
	}

	return model.toEntity(), nil
}

// Update modifies an existing article
func (r *PostgresBlogRepository) Update(ctx context.Context, article *entity.Article) error {
	if article == nil {
		return fmt.Errorf("article cannot be nil")
	}
	if article.ID == "" {
		return fmt.Errorf("article ID cannot be empty")
	}

	model := &ArticleModel{}
	model.fromEntity(article)

	result := r.db.WithContext(ctx).Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update article: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("article with ID %s not found", article.ID)
	}

	return nil
}

// Delete removes an article from the database
func (r *PostgresBlogRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	result := r.db.WithContext(ctx).Delete(&ArticleModel{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete article: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("article with ID %s not found", id)
	}

	return nil
}

// List retrieves articles based on filter criteria
func (r *PostgresBlogRepository) List(ctx context.Context, filter usecase.BlogFilter) ([]*entity.Article, error) {
	query := r.db.WithContext(ctx).Model(&ArticleModel{})

	// Apply filters
	if filter.ID != nil {
		query = query.Where("id = ?", *filter.ID)
	}

	if filter.AuthorID != nil {
		query = query.Where("author_id = ?", *filter.AuthorID)
	}

	if filter.PublishedOnly != nil && *filter.PublishedOnly {
		query = query.Where("published_at IS NOT NULL")
	}

	if filter.DraftsOnly != nil && *filter.DraftsOnly {
		query = query.Where("published_at IS NULL")
	}

	if filter.ScheduledOnly != nil && *filter.ScheduledOnly {
		query = query.Where("published_at > NOW()")
	}

	var models []ArticleModel
	if err := query.Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list articles: %w", err)
	}

	// Convert models to entities
	articles := make([]*entity.Article, len(models))
	for i, model := range models {
		articles[i] = model.toEntity()
	}

	return articles, nil
}
