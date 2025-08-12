# Clean Architecture Demo - Go Implementation

A practical demonstration of Clean Architecture principles in Go, implemented as a blog system to showcase architectural patterns and design decisions.

## Overview

**Goal**: Provide a clean, maintainable, and minimally opinionated starting point 

Our company has adopted a new policy requiring teams to use a consistent architecture for new products and legacy codebase refactoring. Since each team has different experiences, styles, preferences, and legacy codebases, we need a standardized starting point that provides flexibility for different teams and products while maintaining architectural consistency.

To achieve this, our design baseline should be:

* **Beginner-friendly**: Easy to adopt for Go newcomers
* **Development-ready**: Minimal architectural assumptions, maximum flexibility
* **Standards-compliant**: Follows industry and community standards without personal styling choices
* **Well-documented**: Clear documentation that supports developers, search queries, and AI assistance
* **Maintainability-focused**: Prioritizes long-term maintenance over rapid iteration based on product lifecycle needs

We have adopted [Clean Architecture by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) as our architectural foundation. This project demonstrates how Clean Architecture enables **maintainable, testable, and framework-independent** Go applications through proper dependency management and layer separation.

*Inspired by many excellent repositories and resources, particularly [Crescendo Lab's Go Clean Architecture implementation](https://github.com/chatbotgang/go-clean-arch).*



### Core Principles Demonstrated

✅ **Dependency Inversion**: Inner layers never depend on outer layers  
✅ **Separation of Concerns**: Each layer has distinct, well-defined responsibilities  
✅ **Framework Independence**: Business logic works without HTTP framework, databases, or external tools  
✅ **Interface Segregation**: Small, focused interfaces define layer contracts  


## Key Architectural Patterns

### 1. Repository Pattern (Dependency Inversion)

Clean separation between business logic and data persistence through interface contracts:

**Use Case Layer defines interfaces:**
```go
// BlogRepository defines the contract for blog data access operations
type BlogRepository interface {
	Create(ctx context.Context, article *entity.Article) error
	GetByID(ctx context.Context, id string) (*entity.Article, error)
	Update(ctx context.Context, article *entity.Article) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter BlogFilter) ([]*entity.Article, error)
}
```

**Adapter Layer implements interfaces:**
```go
// MemoryBlogRepository provides an in-memory implementation
type MemoryBlogRepository struct {
	articles map[string]*entity.Article
	mu       sync.RWMutex
}

func (r *MemoryBlogRepository) Create(ctx context.Context, article *entity.Article) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.articles[article.ID] = article
	return nil
}
```

**Business logic only depends on interfaces:**
```go
func (app *Application) CreateArticle(ctx context.Context, article *entity.Article) error {
	return app.blogRepo.Create(ctx, article) // Interface method
}
```

**Result**: Multiple implementations can be swapped without changing business logic (MemoryBlogRepository ↔ PostgresBlogRepository ↔ CachedBlogRepository)

### 2. Same Usecase, Different Input Interfaces

The same business logic can be invoked through different interfaces:

**HTTP API (REST):**
```go
// HTTP handler in delivery layer
func PublishArticle(app *usecase.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		articleID := c.Param("id")
		article, err := app.PublishArticle(c.Request.Context(), articleID)
		c.JSON(http.StatusOK, gin.H{"article": article})
	}
}
```

**CLI Admin Tool:**
```go
// CLI command in cmd/admin
func handlePublish(ctx context.Context, app *usecase.Application, args []string) {
	articleID := args[0]
	article, err := app.PublishArticle(ctx, articleID)
	fmt.Printf("✅ Successfully published article: %s\n", article.Title)
}
```

**Result**: Both HTTP `POST /articles/:id/publish` and CLI `admin publish <id>` execute the same `app.PublishArticle()` business logic.

## Project Structure

```
├── cmd/                    # Entry points (Dependency Injection)
│   ├── api/               # HTTP server
│   └── admin/             # CLI admin tool
├── internal/
│   ├── entity/           # Business Entities (innermost)
│   ├── usecase/          # Application Business Rules (where we open maximum flexibility to teams)
│   ├── adapter/          # Interface Adapters
│   ├── delivery/         # External Interfaces
│   ├── platform/         # Infrastructure utilities
│   └── config/           # Configuration management
└── configs/              # Configuration files
```


## Future Enhancements

- **Error Handling**: This project uses Go's native error handling. For a more sophisticated error design, see [Crescendo Lab's implementation](https://github.com/chatbotgang/go-clean-arch)
- **Testing**: Comprehensive testing strategies including unit tests for business logic and integration tests for adapters
- **Validation**: Request validation in handlers, business validation in use cases
- **Caching**: Cache adapter wrapping repository
- **Monitoring**: Observability in infrastructure layer


## References

- [Clean Architecture by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [golang-standards/project-layout](https://github.com/golang-standards/project-layout)
