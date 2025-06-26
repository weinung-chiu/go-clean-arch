package adapter

import (
	"context"
	"fmt"
	"go-clean-arch/internal/entity"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthService struct {
	logger      *slog.Logger
	secretKey   []byte
	tokenExpiry time.Duration
}

func NewJWTAuthService(logger *slog.Logger, secretKey string, tokenExpiry time.Duration) *JWTAuthService {
	return &JWTAuthService{
		logger:      logger.With("component", "jwt_auth_service"),
		secretKey:   []byte(secretKey),
		tokenExpiry: tokenExpiry,
	}
}

func (s *JWTAuthService) GenerateToken(ctx context.Context, participant *entity.Participant) (*entity.AuthToken, error) {
	s.logger.DebugContext(ctx, "Generating token", "participant_id", participant.ID)

	now := time.Now()
	expiresAt := now.Add(s.tokenExpiry)

	claims := &entity.AuthClaims{
		ParticipantID: participant.ID,
		SessionID:     participant.SessionID,
		Nickname:      participant.Nickname,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"participant_id": claims.ParticipantID,
		"session_id":     claims.SessionID,
		"nickname":       claims.Nickname,
		"exp":            expiresAt.Unix(),
		"iat":            now.Unix(),
		"iss":            "go-clean-arch",
		"aud":            "participants",
	})

	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to sign token", "error", err)
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return &entity.AuthToken{
		Token:         tokenString,
		ExpiresAt:     expiresAt,
		ParticipantID: participant.ID,
	}, nil
}

func (s *JWTAuthService) ValidateToken(ctx context.Context, tokenString string) (*entity.AuthClaims, error) {
	s.logger.DebugContext(ctx, "Validating token")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to parse token", "error", err)
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		s.logger.ErrorContext(ctx, "Invalid token")
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		s.logger.ErrorContext(ctx, "Invalid token claims")
		return nil, fmt.Errorf("invalid token claims")
	}

	// Extract claims
	participantID, ok := claims["participant_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid participant_id claim")
	}

	sessionID, ok := claims["session_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid session_id claim")
	}

	nickname, ok := claims["nickname"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid nickname claim")
	}

	return &entity.AuthClaims{
		ParticipantID: participantID,
		SessionID:     sessionID,
		Nickname:      nickname,
	}, nil
}
