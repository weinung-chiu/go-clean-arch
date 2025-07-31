# Go Clean Architecture - Starter Template

A clean, well-structured Go project template following Clean Architecture principles. This repository serves as a starting point for new Go projects with proper architectural boundaries and best practices.

## What is Clean Architecture?

Clean Architecture, popularized by Robert C. Martin (Uncle Bob), is a software design philosophy that emphasizes:

- **Dependency Inversion**: Inner layers never depend on outer layers
- **Separation of Concerns**: Each layer has distinct responsibilities  
- **Testability**: Business logic is independent of frameworks and external systems
- **Maintainability**: Code is organized for long-term maintainability

**Source**: [The Clean Architecture by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

## Project Structure

This project follows Go Clean Architecture patterns with clear layer boundaries:

```
.
├── cmd/                    # Application entry points
│   └── api/               # HTTP API server
├── internal/              # Private application code
│   ├── entity/           # Business entities (innermost layer)
│   ├── usecase/          # Business logic and application services
│   ├── adapter/          # External interface implementations
│   ├── delivery/         # HTTP handlers and routing
│   │   └── api/
│   │       └── router.go
│   ├── platform/         # Common utilities (logger, etc.)
│   │   └── logger/
│   └── config/           # Configuration management
├── configs/              # Configuration files
├── go.mod               # Go module definition
└── CLAUDE.md            # Development guidelines
```

## Layer Responsibilities

### Entities (`internal/entity/`)
- Pure business objects with no external dependencies
- Core business rules and data structures
- Independent of frameworks, databases, or external concerns

### Use Cases (`internal/usecase/`)
- Business logic and application services
- Orchestrates entities and defines application-specific business rules
- Depends only on entities and interfaces

### Adapters (`internal/adapter/`)
- Implementations of interfaces defined in use case layer
- Database repositories, external service clients
- Handles technical details of external integrations

### Delivery (`internal/delivery/`)
- HTTP handlers, routing, request/response transformation
- Framework-specific code (Gin, middleware, etc.)
- Converts external requests to use case calls

## Getting Started

### Prerequisites
- Go 1.24+ 
- Basic understanding of Clean Architecture principles

### Running the Application

1. **Clone and setup**:
   ```bash
   git clone <repository-url>
   cd go-clean-arch
   ```

2. **Configure environment** (optional):
   ```bash
   # Set environment variables directly or create .env file in configs/
   export APP_ENV=dev
   export API_PORT=8080
   export APP_LOG_LEVEL=debug
   ```

3. **Run the API server**:
   ```bash
   go run cmd/api/main.go
   ```

4. **Test the health endpoint**:
   ```bash
   curl http://localhost:8080/health
   ```

### Development Commands

- **Build**: `go build ./...`
- **Test**: `go test ./...`
- **Format**: `go fmt ./...`
- **Vet**: `go vet ./...`

### Process Cleanup (if needed)
If you encounter port conflicts or orphaned processes:
```bash
# Kill orphaned go run processes
pkill -f "go run"

# Free up port 8080 (or your configured port)
lsof -ti:8080 | xargs kill -9

# Clean Go build cache
go clean -cache -modcache -testcache
```

## Architecture Guidelines

### Dependency Rule
- **Inner layers** (entities, use cases) never import outer layers
- **Outer layers** can import and depend on inner layers
- Use **interfaces** to invert dependencies when needed

### Implementation Patterns

1. **Start with entities** - Define your core business objects
2. **Define use case interfaces** - Specify what your application needs
3. **Implement adapters** - Create concrete implementations
4. **Wire in main()** - Dependency injection at application startup

### Adding New Features

1. **Entity Layer**: Define business objects in `internal/entity/`
2. **Use Case Layer**: Create business logic in `internal/usecase/`
3. **Adapter Layer**: Implement repositories/services in `internal/adapter/`
4. **Delivery Layer**: Add HTTP handlers in `internal/delivery/api/`
5. **Wire Dependencies**: Update `cmd/api/main.go` with new dependencies

## Configuration

The application uses environment-based configuration:

- `APP_ENV`: Environment (dev, prod) - defaults to "prod"
- `APP_LOG_LEVEL`: Log level (debug, info, warn, error) - defaults to "warn"  
- `API_PORT`: HTTP server port - defaults to 8080

See `internal/config/` for configuration management.

## Next Steps

This template provides a clean foundation. To build your application:

1. **Define your domain entities** in `internal/entity/`
2. **Implement business logic** in `internal/usecase/`
3. **Add data persistence** with repository implementations in `internal/adapter/`
4. **Create HTTP endpoints** in `internal/delivery/api/`
5. **Add tests** for each layer, especially business logic

## References

- [Clean Architecture by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [golang-standards/project-layout](https://github.com/golang-standards/project-layout)
- [Go Clean Architecture Examples](https://github.com/bxcodec/go-clean-arch)

---

*This project template prioritizes architectural clarity and maintainability over feature completeness. It's designed to be a solid foundation for building scalable Go applications.*