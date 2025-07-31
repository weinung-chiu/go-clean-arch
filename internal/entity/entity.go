package entity

import "time"

// This file contains core domain entities for the Clean Architecture implementation.
// All entities should be pure business objects with no external dependencies.
// Follow Clean Architecture principles: entities should contain business rules
// and be independent of frameworks, databases, or external concerns.

type User struct {
	ID   string
	Name string
}

type Article struct {
	ID          string
	Title       string
	Content     string
	AuthorID    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PublishedAt *time.Time
}

// IsPublished returns true if the article is currently published
func (a *Article) IsPublished() bool {
	return a.PublishedAt != nil && a.PublishedAt.Before(time.Now())
}

// IsScheduled returns true if the article is scheduled for future publication
func (a *Article) IsScheduled() bool {
	return a.PublishedAt != nil && a.PublishedAt.After(time.Now())
}

// IsDraft returns true if the article is not published or scheduled
func (a *Article) IsDraft() bool {
	return a.PublishedAt == nil
}
