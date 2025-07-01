# Real-time Q&A Board (Clean Architecture Example)

This project is a complete, working example of a Real-time Q\&A Board built using the **Go Clean Architecture Framework**. It demonstrates how to structure a modern Go application that handles both RESTful API requests and real-time WebSocket communication.

The goal is to provide a practical, easy-to-understand template for teams looking to adopt Go and Clean Architecture for building scalable and maintainable web services.

## Core Concepts Demonstrated

* **Clean Architecture:** A strict separation of concerns into four distinct layers (`entity`, `usecase`, `adapter`, `delivery`).
* **Hybrid API Design:** A combination of a RESTful API for session management and a WebSocket API for all real-time events.
* **Persist-First Strategy:** A robust flow for handling data submission, inspired by modern messaging architectures like Slack's.
* **Beginner-Friendly Code:** Clear and idiomatic Go code that avoids overly complex or dogmatic patterns.

-----

## Architecture Overview

The framework is organized into a `3+1` layer design to enforce the Dependency Rule, ensuring the business logic (`usecase`) remains independent of external concerns like databases or web frameworks.

1.  **Entity**: The core data structures of the application (e.g., `Session`, `Question`). These are pure Go structs.
2.  **Usecase**: Contains the business logic and orchestrates the flow of data. It defines interfaces that the outer layers must implement.
3.  **Adapter**: The "glue" layer. Implements the interfaces defined by the use cases, connecting to databases, caches, or other services.
4.  **Delivery**: The entrypoint layer that handles incoming requests (e.g., HTTP handlers, WebSocket controllers) and passes them to the use case layer.

-----

## API Design

The application exposes both REST and WebSocket endpoints to provide a responsive user experience.

### RESTful API (HTTP)

Used for initial data loading and session management.

| Method | Endpoint | Description | Audience |
| --- | --- | --- | --- |
| `POST` | `/sessions` | Creates a new Q\&A `Session`. | Presenter/Admin |
| `GET` | `/sessions/{sessionID}/questions` | Fetches all questions for a session. Used on initial page load. | Participant |
| `POST` | `/sessions/{sessionID}/questions` | Submits a new question to a session. | Participant |

\<br/\>

### WebSocket API

Used for all real-time events within a Q\&A session.

| Endpoint | Action |
| --- | --- |
| `GET /ws/{sessionID}` | A client joins a session and upgrades the connection to a WebSocket. |

**Client → Server Messages:**

| Type | Payload | Description |
| --- | --- | --- |
| `QUESTION_SUBMIT` | `{ "text": "..." }` | Submits a new question in real-time. <!-- ⚠️ Not implemented via WebSocket; handled via REST API instead. --> |
| `QUESTION_UPVOTE` | `{ "question_id": "..." }` | Upvotes an existing question. <!-- ⚠️ Not implemented via WebSocket; handled via REST API instead. --> |

**Server → Client Messages:**

| Type | Payload | Description |
| --- | --- | --- |
| `BOARD_STATE_UPDATED` | `{ "questions": [...] }` | Broadcasts the complete, sorted list of questions to all participants, keeping everyone in sync. |
| `ERROR` | `{ "message": "..." }` | Sends a specific error message to a single client (e.g., "Already upvoted"). <!-- ⚠️ Not implemented as a dedicated WebSocket message type in the current codebase. --> |

-----

## Core Flow: Submitting a Question

To ensure data integrity and a responsive feel, this project uses a "persist-first" design for creating new questions.

1.  The client submits a new question via an **HTTP POST** request.
2.  The **Web App** (delivery/adapter layers) immediately persists the question to the database.
3.  Upon successful persistence, the **Usecase** layer triggers a broadcast event.
4.  The **Channel Server** (WebSocket handler) receives this event, fetches the newly updated and sorted list of all questions, and pushes the complete state to all connected clients.

<!-- end list -->

```mermaid
graph TD
    subgraph Client
        A[Participant's Browser]
    end

    subgraph Backend
        B[Web App (REST Handler)]
        C[Database]
        D[Broadcaster]
        E[Channel Server (WebSocket Handler)]
    end

    A -- 1. POST /questions --> B
    B -- 2. Persist Question --> C
    B -- 3. Trigger Broadcast --> D
    D -- 4. Push "State Updated" Event --> E
    E -- 5. Broadcast New Board State --> A
```