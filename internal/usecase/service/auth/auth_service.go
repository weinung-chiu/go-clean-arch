package auth

import (
	"context"
	"errors"
	"go-clean-arch/internal/entity"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	participantRepo ParticipantRepository
	sessionRepo     SessionRepository
	authServer      AuthServer
}

type ParticipantRepository interface {
	CreateParticipant(ctx context.Context, participant *entity.Participant) error
	GetParticipant(ctx context.Context, filter entity.ParticipantFilter) (*entity.Participant, error)
	UpdateParticipantLastSeen(ctx context.Context, participantID string) error
}

type SessionRepository interface {
	GetSessionByID(ctx context.Context, sessionID string) (*entity.Session, error)
}

type AuthServer interface {
	GenerateToken(ctx context.Context, participant *entity.Participant) (*entity.AuthToken, error)
	ValidateToken(ctx context.Context, tokenString string) (*entity.AuthClaims, error)
}

var (
	ErrNicknameAlreadyTaken = "nickname already taken in this session"
	ErrParticipantNotFound  = "participant not found"
	ErrInvalidToken         = "invalid or expired token"
)

func NewAuthService(participantRepo ParticipantRepository, sessionRepo SessionRepository, authServer AuthServer) *Service {
	return &Service{
		participantRepo: participantRepo,
		sessionRepo:     sessionRepo,
		authServer:      authServer,
	}
}

func (s *Service) RegisterParticipant(ctx context.Context, sessionID, nickname string) (*entity.AuthToken, error) {
	_, err := s.sessionRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	existingParticipant, err := s.participantRepo.GetParticipant(ctx, entity.ParticipantFilter{
		SessionID: &sessionID,
		Nickname:  &nickname,
	})
	if err == nil && existingParticipant != nil {
		return nil, errors.New(ErrNicknameAlreadyTaken)
	}

	participant := &entity.Participant{
		ID:               uuid.NewString(),
		SessionID:        sessionID,
		Nickname:         nickname,
		UpvotedQuestions: make(map[string]bool),
		CreatedAt:        time.Now(),
		LastSeenAt:       time.Now(),
	}

	if err := s.participantRepo.CreateParticipant(ctx, participant); err != nil {
		return nil, err
	}

	token, err := s.authServer.GenerateToken(ctx, participant)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *Service) LoginParticipant(ctx context.Context, sessionID, nickname string) (*entity.AuthToken, error) {
	participant, err := s.participantRepo.GetParticipant(ctx, entity.ParticipantFilter{
		SessionID: &sessionID,
		Nickname:  &nickname,
	})
	if err != nil {
		return nil, errors.New(ErrParticipantNotFound)
	}

	if err := s.participantRepo.UpdateParticipantLastSeen(ctx, participant.ID); err != nil {
		// log warning if needed
	}

	token, err := s.authServer.GenerateToken(ctx, participant)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *Service) ValidateParticipantToken(ctx context.Context, tokenString string) (*entity.Participant, error) {
	claims, err := s.authServer.ValidateToken(ctx, tokenString)
	if err != nil {
		return nil, errors.New(ErrInvalidToken)
	}

	participant, err := s.participantRepo.GetParticipant(ctx, entity.ParticipantFilter{
		ID: &claims.ParticipantID,
	})
	if err != nil {
		return nil, errors.New(ErrParticipantNotFound)
	}

	if err := s.participantRepo.UpdateParticipantLastSeen(ctx, participant.ID); err != nil {
		// log warning if needed
	}

	return participant, nil
}
