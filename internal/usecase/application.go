package usecase

import (
	"context"
	"github.com/google/uuid"
	"go-clean-arch/internal/entity"
	"log/slog"
	"time"
)

type Application struct {
	logger          *slog.Logger
	sessionRepo     SessionRepository
	questionRepo    QuestionRepository
	sessionEventBus SessionEventBus

	latestSessionID string // for testing purposes
}

func NewApplication(params NewApplicationParams) (*Application, error) {
	return &Application{
		logger:          params.Logger.With("component", "application"),
		sessionRepo:     params.SessionRepo,
		questionRepo:    params.QuestionRepo,
		sessionEventBus: params.EventBus,
	}, nil
}

type NewApplicationParams struct {
	Logger       *slog.Logger
	SessionRepo  SessionRepository
	QuestionRepo QuestionRepository
	EventBus     SessionEventBus
}

type SessionRepository interface {
	CreateSession(ctx context.Context, session *entity.Session) error
	GetSessionByID(ctx context.Context, sessionID string) (*entity.Session, error)
	AddParticipantToSession(ctx context.Context, sessionID string, participant *entity.Participant) error
	ListSessions(ctx context.Context) ([]*entity.Session, error)
}

type QuestionRepository interface {
	CreateQuestion(ctx context.Context, question *entity.Question) error
	ListQuestionsBySession(ctx context.Context, sessionID string) ([]*entity.Question, error)
}

// SessionEventType defines the type of event in a session.
type SessionEventType string

const (
	SessionEventQuestionSubmitted SessionEventType = "question_submitted"
)

// SessionEvent represents an event in a session.
type SessionEvent struct {
	Type      SessionEventType
	SessionID string
	Payload   any
	Timestamp time.Time
}

type SessionEventBus interface {
	Publish(ctx context.Context, event *SessionEvent) error
	Subscribe(ctx context.Context, sessionID string) (<-chan *SessionEvent, error)
}

func (a *Application) NewSession(ctx context.Context, name string) (*entity.Session, error) {
	a.logger.DebugContext(ctx, "Creating new session", "name", name)

	session := &entity.Session{
		ID:                   uuid.NewString(),
		Title:                name,
		IsAcceptingQuestions: true,
	}

	if err := a.sessionRepo.CreateSession(ctx, session); err != nil {
		a.logger.ErrorContext(ctx, "Failed to create session", "error", err)
		return nil, err
	}

	a.logger.InfoContext(ctx, "New session created", "session_id", session.ID)
	a.latestSessionID = session.ID
	return session, nil
}

// ListSessions returns all sessions.
func (a *Application) ListSessions(ctx context.Context) ([]*entity.Session, error) {
	a.logger.DebugContext(ctx, "Listing all sessions")
	return a.sessionRepo.ListSessions(ctx)
}

// GetSession returns a session and its questions.
func (a *Application) GetSession(ctx context.Context, sessionID string) (*entity.Session, []*entity.Question, error) {
	session, err := a.sessionRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		a.logger.ErrorContext(ctx, "Failed to get session", "error", err)
		return nil, nil, err
	}
	questions, err := a.questionRepo.ListQuestionsBySession(ctx, sessionID)
	if err != nil {
		a.logger.ErrorContext(ctx, "Failed to list questions", "error", err)
		return nil, nil, err
	}
	return session, questions, nil
}

// SubmitQuestion creates a new question for a session.
func (a *Application) SubmitQuestion(ctx context.Context, sessionID, text, nickname string) (*entity.Question, error) {
	question := &entity.Question{
		ID:             uuid.NewString(),
		SessionID:      sessionID,
		Text:           text,
		AuthorNickname: nickname,
		Upvotes:        0,
		CreatedAt:      time.Now(),
	}
	if err := a.questionRepo.CreateQuestion(ctx, question); err != nil {
		a.logger.ErrorContext(ctx, "Failed to create question", "error", err)
		return nil, err
	}
	a.logger.InfoContext(ctx, "New question submitted", "question_id", question.ID)

	event := &SessionEvent{
		Type:      SessionEventQuestionSubmitted,
		SessionID: sessionID,
		Payload:   question,
		Timestamp: time.Now(),
	}
	err := a.sessionEventBus.Publish(ctx, event)
	if err != nil {
		a.logger.ErrorContext(ctx, "Failed to publish question submitted event", "error", err)
	}

	return question, nil
}

// SubscribeToSessionEvents subscribes to session events.
func (a *Application) SubscribeToSessionEvents(ctx context.Context, sessionID string) (<-chan *SessionEvent, error) {
	a.logger.DebugContext(ctx, "Subscribing to session events", "session_id", sessionID)
	return a.sessionEventBus.Subscribe(ctx, sessionID)
}

func (a *Application) SubscribeToLastestSessionEvents(ctx context.Context) (<-chan *SessionEvent, error) {
	return a.SubscribeToSessionEvents(ctx, a.latestSessionID)
}
