# AI Prompting Guide for the Go Q&A Board Framework

## 1. Core Philosophy & Design Goals

This framework implements a **Clean Architecture** style in Golang. Its primary goal is to provide a clear, maintainable, and scalable structure for building real-time web services. The design is intended to be easy for teams with diverse backgrounds to adopt and maintain.

The single most important principle is the **Dependency Rule**: all source code dependencies must point inwards. Outer layers can depend on inner layers, but inner layers must remain completely independent of the outer layers.

The application follows a **Persist-First Design**, inspired by Slack's architecture. When a user action (like submitting a question) occurs, the data is first saved to a persistent repository, and only *then* is a notification broadcast to clients.

## 2. Architectural Layers & Directory Mapping

> **Note:** The overall directory structure is heavily inspired by the community-standard [golang-standards/project-layout](https://github.com/golang-standards/project-layout).

Our framework is organized into four primary layers, which correspond to the directories in the `internal/` package.

#### Layer 1: `entity` (The Core)
- **Directory:** `internal/entity/`
- **Purpose:** Contains the core business objects and enterprise-wide business rules. These are the nouns of our application.
- **Rules:**
    - Has ZERO external dependencies. It cannot import any other layer in our project.
    - Contains plain Go structs (`Session`, `Question`, `Participant`).
- **Example:** `type Question struct { ID string; Text string; ... }`

#### Layer 2: `usecase` (The Application Logic)
- **Directory:** `internal/usecase/`
- **Purpose:** Orchestrates the flow of data and implements application-specific business logic. This is where we define *what* the application does.
- **Rules:**
    - Depends only on the `entity` layer.
    - Defines the interfaces that outer layers must implement. Key interfaces include repositories for data access (e.g., `QuestionRepository`) and broadcasters for real-time events (e.g., `EventBroadcaster`).
- **Example:** The `Application` struct in `application.go` holds the business logic for creating and upvoting questions.

#### Layer 3: `adapter` (The Implementations)
- **Directory:** `internal/adapter/`
- **Purpose:** Implements the interfaces defined in the `usecase` layer. This is *how* the application's dependencies are satisfied.
- **Rules:**
    - Depends on `usecase` for the interface definitions.
    - Contains all infrastructure-specific code. This is where database connections, external API clients, and in-memory stores live.
- **Example:** `memory_repo.go` provides a concrete implementation of a `usecase` repository interface, storing data in memory for the PoC.

#### Layer 4: `delivery` (The Entry Points)
- **Directory:** `internal/delivery/`
- **Purpose:** Exposes the application's use cases to the outside world. This is the primary entry point for users and external systems.
- **Rules:**
    - Depends on the `usecase` layer to execute application logic.
    - Handles concerns specific to the delivery mechanism, such as HTTP routing, request/response models (DTOs), and WebSocket connection management.
- **Example:** `delivery/api/` contains the RESTful API handlers and routers. A `delivery/ws/` would contain the WebSocket handlers.

#### Outer Layers & Infrastructure
- **Directories:** `cmd/`, `configs/`, `internal/platform/`
- **Purpose:** These are the outermost components responsible for wiring everything together.
- `cmd/`: Contains the `main.go` entry points for different applications (e.g., `api` server, `bot`). This is where the application is initialized and started.
- `configs/`: Handles loading and parsing application configuration.
- `internal/platform/`: Contains application-wide utilities like the structured logger.

## 3. Golden Path Example: A User Submits a Question

This flow demonstrates how the layers interact:

1.  **Request Ingress (`delivery`):** A client sends a `POST /sessions/{id}/questions` request. The router in `delivery/api/router.go` directs it to the appropriate handler in `handler.go`.
2.  **DTO & Validation (`delivery`):** The handler parses the request body into a DTO defined in `model.go` and performs initial validation.
3.  **Call Usecase (`delivery` -> `usecase`):** The handler calls the `CreateQuestion` method on the `usecase.Application` instance.
4.  **Execute Business Logic (`usecase`):**
    - The use case creates a new `entity.Question` from the input.
    - It calls the `SaveQuestion` method on its repository interface (`QuestionRepository`).
5.  **Persist Data (`adapter`):** The `adapter.MemoryRepository` (which implements the `QuestionRepository` interface) saves the question to its in-memory store.
6.  **Broadcast Event (`usecase`):** After the question is successfully saved, the use case calls the `BroadcastBoardState` method on its `EventBroadcaster` interface.
7.  **Push Real-time Update (`delivery`):** A component in the `delivery` layer (e.g., a WebSocket handler that implements the `EventBroadcaster` interface) receives the broadcast and sends the updated `BOARD_STATE_UPDATED` message to all connected WebSocket clients.

## 4. Do's and Don'ts

* **DO:** Define all repository and external service interfaces within the `usecase` layer.
* **DO:** Use plain `entity` objects for all core business logic.
* **DON'T:** Allow the `entity` or `usecase` layers to import `delivery` or `adapter`.
* **DON'T:** Place HTTP-specific code (like `gin.Context`) or SQL code inside a `usecase` function.
* **DO:** Use DTOs in the `delivery` layer to decouple API models from core `entity` models.