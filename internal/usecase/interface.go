package usecase

import (
	"context"
	"go-clean-arch/internal/entity"
)

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

type ParticipantRepository interface {
	CreateParticipant(ctx context.Context, participant *entity.Participant) error
	GetParticipant(ctx context.Context, filter entity.ParticipantFilter) (*entity.Participant, error)
	UpdateParticipantLastSeen(ctx context.Context, participantID string) error
	ListParticipants(ctx context.Context, filter entity.ParticipantFilter) ([]*entity.Participant, error)
}
