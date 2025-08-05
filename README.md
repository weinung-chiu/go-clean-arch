# Clean Architecture Demo - Go Implementation

A practical demonstration of Clean Architecture principles in Go, implemented as a blog system to showcase architectural patterns and design decisions.

## Architecture Overview

This project demonstrates how Clean Architecture enables **maintainable, testable, and framework-independent** Go applications through proper dependency management and layer separation.

### Core Principles Demonstrated

✅ **Dependency Inversion**: Inner layers never depend on outer layers  
✅ **Separation of Concerns**: Each layer has distinct, well-defined responsibilities  
✅ **Framework Independence**: Business logic works without Gin, databases, or external tools  
✅ **Interface Segregation**: Small, focused interfaces define layer contracts  
✅ **Single Responsibility**: Each component has one reason to change

## Layer Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Frameworks & Drivers                    │
│  (HTTP Server, CLI, Database, External APIs)               │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                Interface Adapters                   │   │
│  │     (Controllers, Presenters, Repositories)        │   │
│  │                                                     │   │
│  │  ┌─────────────────────────────────────────────┐   │   │
│  │  │              Application Business Rules     │   │   │
│  │  │              (Use Cases, Interactors)       │   │   │
│  │  │                                             │   │   │
│  │  │  ┌─────────────────────────────────────┐   │   │   │
│  │  │  │        Enterprise Business Rules    │   │   │   │
│  │  │  │              (Entities)             │   │   │   │
│  │  │  └─────────────────────────────────────┘   │   │   │
│  │  └─────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### Implementation Layers

| Layer | Package | Responsibility | Dependencies |
|-------|---------|----------------|--------------|
| **Entities** | `internal/entity/` | Business objects & rules | None |
| **Use Cases** | `internal/usecase/` | Application business logic | Entities only |
| **Adapters** | `internal/adapter/` | Interface implementations | Use Cases + External |
| **Delivery** | `internal/delivery/` | External interfaces | Use Cases |
| **Main** | `cmd/*/` | Dependency injection | All layers |

## Key Architectural Patterns

### 1. Dependency Inversion Pattern

**Use Case Layer defines interfaces:**
```go
// internal/usecase/interface.go
type BlogRepository interface {
    Create(ctx context.Context, article *entity.Article) error
    GetByID(ctx context.Context, id string) (*entity.Article, error)
    // ... more methods
}
```

**Adapter Layer implements interfaces:**
```go
// internal/adapter/memory_blog_repository.go
type MemoryBlogRepository struct { ... }

func (r *MemoryBlogRepository) Create(ctx context.Context, article *entity.Article) error {
    // Implementation details...
}
```

**Result**: Business logic never depends on specific database implementations.

### 2. Repository Pattern

Clean separation between business logic and data persistence:

```go
// Use Case depends on interface, not implementation
type Application struct {
    blogRepo BlogRepository  // Interface, not concrete type
}

// Easy to swap implementations:
// - MemoryBlogRepository (current)
// - PostgreSQLRepository (future)
// - MockRepository (testing)
```

### 3. Multiple Entry Points

Different interfaces sharing the same business logic:

```
cmd/api/     → HTTP handlers → Application methods
cmd/admin/   → CLI commands  → Application methods  
cmd/dev/     → Test scripts  → Application methods
```

Same business rules, different interfaces.

## Project Structure

```
├── cmd/                    # Entry points (Dependency Injection)
│   ├── api/               # HTTP server
│   ├── admin/             # CLI admin tool
│   └── dev/               # Development utilities
├── internal/
│   ├── entity/           # 🔵 Business Entities (innermost)
│   ├── usecase/          # 🟢 Application Business Rules
│   ├── adapter/          # 🟡 Interface Adapters
│   ├── delivery/         # 🔴 External Interfaces
│   ├── platform/         # Infrastructure utilities
│   └── config/           # Configuration management
└── configs/              # Configuration files
```

## Architectural Benefits Demonstrated

### 1. **Framework Independence**
- Business logic has zero external dependencies
- Can swap from Gin to Echo without changing business rules
- Database-agnostic through repository interfaces

### 2. **Testability**  
- Pure business logic can be unit tested without frameworks
- Mock repositories for testing use cases
- Integration tests at delivery layer

### 3. **Multiple Interfaces**
- Same business logic supports HTTP API AND CLI commands
- Easy to add GraphQL, gRPC, or messaging interfaces

### 4. **Maintainability**
- Changes in one layer don't cascade to others
- Clear boundaries make code easy to understand
- Dependency direction prevents architectural decay

## Quick Demo

### Test the Architecture

```bash
# 1. Start HTTP API
go run cmd/api/main.go

# 2. Create and publish article via HTTP
curl -X POST localhost:80/api/v1/articles \
  -H "Content-Type: application/json" \
  -d '{"title":"Test","content":"Content","author_id":"user"}'

curl -X POST localhost:80/api/v1/articles/{id}/publish

# 3. Use CLI admin (same business logic, different interface)
go run cmd/admin/main.go list

# 4. Run development utilities
go run cmd/dev/main.go
```

## Development Commands

```bash
# Build and test
go build ./...
go test ./...
go fmt ./...
go vet ./...

# Clean up processes if needed
pkill -f "go run"
```

## Key Design Decisions

### Why Blog Domain?
- **Familiar**: Everyone understands articles and publishing
- **Rich Logic**: Draft→Published workflow demonstrates entity methods
- **Multiple Operations**: CRUD + business operations (publish, delete)

### Why In-Memory Repository?
- **Focus on Architecture**: Not distracted by database setup
- **Interface Demonstration**: Easy to swap for real database
- **Zero Dependencies**: Clone and run immediately

### Why Multiple Entry Points?
- **Framework Independence**: Same logic, different interfaces  
- **Dependency Injection**: Shows how to wire Clean Architecture
- **Real-World Pattern**: Common in production applications

## Next Steps

Replace components while keeping architecture intact:

1. **Database**: Swap `MemoryBlogRepository` for `PostgreSQLRepository`
2. **Authentication**: Add JWT middleware in delivery layer
3. **Validation**: Request validation in handlers, business validation in use cases
4. **Caching**: Cache adapter wrapping repository
5. **Monitoring**: Observability in infrastructure layer

## References

- [Clean Architecture by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [golang-standards/project-layout](https://github.com/golang-standards/project-layout)

---

*Clean Architecture isn't just theory—it's a practical approach to building maintainable Go applications. This project shows how proper separation of concerns makes code easy to understand, test, and evolve.*