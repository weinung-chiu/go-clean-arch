package usecase

import (
	"context"
	"github.com/google/uuid"
	"go-clean-arch/internal/entity"
	"log/slog"
	"time"
)

type Application struct {
	logger                 *slog.Logger
	sessionRepo            SessionRepository
	questionRepo           QuestionRepository
	clientEventBroadcaster ClientEventBroadcaster
}

func NewApplication(params NewApplicationParams) (*Application, error) {
	return &Application{
		logger:                 params.Logger.With("component", "application"),
		sessionRepo:            params.SessionRepo,
		questionRepo:           params.QuestionRepo,
		clientEventBroadcaster: newClientEventBroadcaster(),
	}, nil
}

type NewApplicationParams struct {
	Logger       *slog.Logger
	SessionRepo  SessionRepository
	QuestionRepo QuestionRepository
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
	GetQuestionByID(ctx context.Context, questionID string) (*entity.Question, error)
	UpvoteQuestionByID(ctx context.Context, questionID, participantID, participantNickname string) (bool, error)
}

type ClientEventBroadcaster interface {
	Broadcast(ctx context.Context, event *ClientEventQuestionUpdated) error
	Subscribe(ctx context.Context, sessionID string) (<-chan *ClientEventQuestionUpdated, error)
}

// ClientEventQuestionUpdated is the event type for when a question is updated in a session.
// for simplicity, this is only type of event we handle in this example.
type ClientEventQuestionUpdated struct {
	SessionID string
	Timestamp time.Time

	// payloads
	Questions    []*entity.Question
	Participants []*entity.Participant
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
// TODO: should pass *entity.Participant instead of nickname
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
	questions, err := a.questionRepo.ListQuestionsBySession(ctx, sessionID)
	if err != nil {
		a.logger.ErrorContext(ctx, "Failed to fetch updated questions", "error", err)
		return nil, err
	}
	event := &ClientEventQuestionUpdated{
		SessionID: sessionID,
		Timestamp: time.Now(),

		Questions:    questions,
		Participants: nil, // Assuming we don't track participants in this example
	}
	err = a.clientEventBroadcaster.Broadcast(ctx, event)
	if err != nil {
		a.logger.ErrorContext(ctx, "Failed to publish question submitted event", "error", err)
	}

	return question, nil
}

// SubscribeToClientEvent subscribes to session events.
func (a *Application) SubscribeToClientEvent(ctx context.Context, sessionID string) (<-chan *ClientEventQuestionUpdated, error) {
	a.logger.DebugContext(ctx, "Subscribing to session events", "session_id", sessionID)
	return a.clientEventBroadcaster.Subscribe(ctx, sessionID)
}

// UpvoteQuestion Upvote a question by its ID only (no sessionID required)
// TODO: should pass *entity.Participant instead of nickname
func (a *Application) UpvoteQuestion(ctx context.Context, questionID, participantID, participantNickname string) error {
	success, err := a.questionRepo.UpvoteQuestionByID(ctx, questionID, participantID, participantNickname)
	if err != nil {
		a.logger.ErrorContext(ctx, "Failed to upvote question by ID", "error", err)
		return err
	}
	if !success {
		return nil
	}
	// Broadcast updated question's session's questions
	q, err := a.questionRepo.GetQuestionByID(ctx, questionID)
	if err != nil {
		return err
	}
	questions, err := a.questionRepo.ListQuestionsBySession(ctx, q.SessionID)
	if err != nil {
		return err
	}
	event := &ClientEventQuestionUpdated{
		SessionID:    q.SessionID,
		Timestamp:    time.Now(),
		Questions:    questions,
		Participants: nil,
	}
	_ = a.clientEventBroadcaster.Broadcast(ctx, event)
	return nil
}
