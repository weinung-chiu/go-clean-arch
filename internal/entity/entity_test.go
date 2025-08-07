package entity

import (
	"testing"
	"time"
)

func TestArticle_IsPublished(t *testing.T) {
	tests := []struct {
		name        string
		publishedAt *time.Time
		want        bool
	}{
		{
			name:        "nil published date - not published",
			publishedAt: nil,
			want:        false,
		},
		{
			name:        "published in the past - is published",
			publishedAt: timePtr(time.Now().Add(-1 * time.Hour)),
			want:        true,
		},
		{
			name:        "published now - is published",
			publishedAt: timePtr(time.Now().Add(-1 * time.Second)),
			want:        true,
		},
		{
			name:        "scheduled for future - not yet published",
			publishedAt: timePtr(time.Now().Add(1 * time.Hour)),
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article := &Article{
				ID:          "test-id",
				Title:       "Test Article",
				Content:     "Test content",
				AuthorID:    "test-author",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				PublishedAt: tt.publishedAt,
			}

			if got := article.IsPublished(); got != tt.want {
				t.Errorf("Article.IsPublished() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArticle_IsScheduled(t *testing.T) {
	tests := []struct {
		name        string
		publishedAt *time.Time
		want        bool
	}{
		{
			name:        "nil published date - not scheduled",
			publishedAt: nil,
			want:        false,
		},
		{
			name:        "published in the past - not scheduled",
			publishedAt: timePtr(time.Now().Add(-1 * time.Hour)),
			want:        false,
		},
		{
			name:        "scheduled for future - is scheduled",
			publishedAt: timePtr(time.Now().Add(1 * time.Hour)),
			want:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article := &Article{
				ID:          "test-id",
				Title:       "Test Article",
				Content:     "Test content",
				AuthorID:    "test-author",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				PublishedAt: tt.publishedAt,
			}

			if got := article.IsScheduled(); got != tt.want {
				t.Errorf("Article.IsScheduled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArticle_IsDraft(t *testing.T) {
	tests := []struct {
		name        string
		publishedAt *time.Time
		want        bool
	}{
		{
			name:        "nil published date - is draft",
			publishedAt: nil,
			want:        true,
		},
		{
			name:        "published in the past - not draft",
			publishedAt: timePtr(time.Now().Add(-1 * time.Hour)),
			want:        false,
		},
		{
			name:        "scheduled for future - not draft",
			publishedAt: timePtr(time.Now().Add(1 * time.Hour)),
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article := &Article{
				ID:          "test-id",
				Title:       "Test Article",
				Content:     "Test content",
				AuthorID:    "test-author",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				PublishedAt: tt.publishedAt,
			}

			if got := article.IsDraft(); got != tt.want {
				t.Errorf("Article.IsDraft() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to create time pointers
func timePtr(t time.Time) *time.Time {
	return &t
}
