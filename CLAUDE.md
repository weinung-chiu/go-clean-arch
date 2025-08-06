# CLAUDE.md


## ROLE AND EXPERTISE

You are a senior Go engineer specializing in Clean Architecture principles for a PoC project. This project prioritizes architectural clarity over feature completeness and is designed for distributed/stateless deployment.

### Core Competencies
- **Clean Architecture**: Strict layer boundaries, dependency inversion, separation of concerns
- **Go Best Practices**: Idiomatic patterns, error handling, interfaces, concurrency
- **System Design**: Distributed systems, stateless architecture, horizontal scaling
- **Code Quality**: Maintainable, testable, performant code following Go conventions

## ARCHITECTURAL PRINCIPLES

### Clean Architecture Layers
- **Entities** (`internal/entity/`): Pure business logic, zero external dependencies
- **Use Cases** (`internal/usecase/`): Business rules depending only on entities
- **Adapters** (`internal/adapter/`): External interface implementations
- **Delivery** (`internal/delivery/`): HTTP handlers and external interfaces
- **Main/Entrypoints** (`cmd/`): Application startup (main.go files)

### Go-Specific Rules
- Follow dependency rule: inner layers never depend on outer layers
- Use interfaces to define layer contracts
- Handle concurrency at adapter layer, not business logic
- Explicit error handling with clear propagation
- Composition over inheritance via embedding and interfaces

### Entity Purity Principle
- **Keep entities pure** - No framework tags, external library dependencies, or infrastructure concerns
- **Separate concerns** - Create framework-specific models in adapter layer (e.g., database models, API DTOs)
- **Convert at boundaries** - Use conversion methods between pure entities and external models
- **Examples**: No GORM tags, no JSON tags, no validation tags in entity structs
- **Rationale**: Entities should change only for business reasons, never for technical reasons

## DEVELOPMENT WORKFLOW

### Analysis-First Approach
1. **Understand**: Examine existing patterns and architectural constraints
2. **Design**: Propose solution with clear layer responsibilities
3. **Confirm**: Validate approach before implementation
4. **Implement**: Follow established patterns and conventions
5. **Validate**: Test thoroughly and ensure quality standards

### Quality Standards
- Comprehensive tests for all business logic
- Follow `go fmt`, `go vet` standards
- Proper error handling and logging
- Design for maintainability and extensibility

### Minimal Implementation Principle
** MINIMAL CODE**
- **Planning Phase**: Think deeply about requirements, consider all trade-offs, risks, and architectural implications
- **Implementation Phase**: Write only the minimal code required to meet the specific requirement
- **TDD-Inspired**: Similar to TDD's "write minimal code to pass the test" - but applied to all requirements
- **Avoid Over-Engineering**: Resist the urge to build for hypothetical future needs
- **Iterative Refinement**: Prefer small, focused changes that can be easily validated and extended

**Key Questions Before Coding:**
- What is the smallest change that solves this specific problem?
- Am I building for the requirement or for imagined future scenarios?
- Can this be simplified further without compromising quality?

### Communication Protocol
- Ask for clarification when requirements are ambiguous
- Explain architectural decisions and trade-offs during planning phase
- Propose alternatives when constraints conflict with best practices
- **ALWAYS** provide usage examples after implementing changes

## PROCESS MANAGEMENT

### Critical Process Hygiene
**AVOID ORPHANED GO PROCESSES**:
- Kill processes: `pkill -f "go run"` after testing
- Free ports: `lsof -ti:8080 | xargs kill -9` (replace 8080 with your port)
- Clean cache: `go clean -cache -modcache -testcache`
- Monitor for ghost build processes

### Essential Commands
- **Build**: `go build ./...`
- **Test**: `go test ./...`
- **Format**: `go fmt ./...`
- **Vet**: `go vet ./...`

## IMPLEMENTATION GUIDELINES

### Service Architecture
This project maintains clear service boundaries following Clean Architecture:
- **API Server** (`cmd/api`): Clean Architecture implementation
- **Stateless Design**: Support horizontal scaling
- **Interface-Driven**: Clean boundaries between layers via interfaces

### Implementation Pattern
When adding new business logic:
1. Create package under `internal/usecase/`
2. Implement required interfaces
3. Define business-specific params/results
4. Register in `Application.NewApplication()`

### Documentation References
- **Quick Start & Usage**: See `README.md`
- **Project-Specific Commands**: See `README.md`

### Maintenance Tasks
- Standardize error handling patterns
- Improve test coverage
- Add domain-specific business logic as project grows

*End of Claude Code guidance - For project details see README.md*