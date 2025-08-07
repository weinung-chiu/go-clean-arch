package adapter

import (
	"context"
	"errors"
	"go-clean-arch/internal/config"
	"go-clean-arch/internal/entity"
	"go-clean-arch/internal/usecase"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// JWTAuthService implements AuthService interface using JWT tokens
type JWTAuthService struct {
	config   *config.AppConfig
	userRepo usecase.UserRepository
}

// NewJWTAuthService creates a new JWT authentication service
func NewJWTAuthService(cfg *config.AppConfig, userRepo usecase.UserRepository) *JWTAuthService {
	return &JWTAuthService{
		config:   cfg,
		userRepo: userRepo,
	}
}

// JWTClaims represents the JWT token claims with user data
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	jwt.RegisteredClaims
}

// LoginWithPassword authenticates a user with username/password and returns a JWT token
func (j *JWTAuthService) LoginWithPassword(username, password string) (*usecase.AuthResult, error) {
	ctx := context.Background()

	// Get user by username
	user, err := j.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Generate JWT token with user data in claims
	token, expiresAt, err := j.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &usecase.AuthResult{
		Token:     token,
		User:      user,
		ExpiresAt: expiresAt,
	}, nil
}

// RegisterWithPassword creates a new user account with username/password
func (j *JWTAuthService) RegisterWithPassword(username, password string) (*usecase.AuthResult, error) {
	ctx := context.Background()

	// Validate input
	if username == "" {
		return nil, errors.New("username is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}
	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters long")
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create new user
	now := time.Now()
	user := &entity.User{
		ID:           uuid.New().String(),
		Name:         username, // Use username as display name for simplicity
		Username:     username,
		PasswordHash: string(passwordHash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Store user
	if err := j.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate JWT token with user data in claims
	token, expiresAt, err := j.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &usecase.AuthResult{
		Token:     token,
		User:      user,
		ExpiresAt: expiresAt,
	}, nil
}

// ValidateToken validates a JWT token and returns user data from claims (no repo lookup)
func (j *JWTAuthService) ValidateToken(tokenString string) (*entity.User, error) {
	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(j.config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims and return user data (stateless - no repo lookup)
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return &entity.User{
			ID:       claims.UserID,
			Username: claims.Username,
			Name:     claims.Name,
			// Note: PasswordHash not included in token for security
			// CreatedAt/UpdatedAt not included to keep token minimal
		}, nil
	}

	return nil, errors.New("invalid token claims")
}

// generateToken creates a JWT token with user data embedded in claims
func (j *JWTAuthService) generateToken(user *entity.User) (string, time.Time, error) {
	// Calculate expiration time from config
	expirationMinutes := time.Duration(j.config.JWTExpiry) * time.Minute
	expiresAt := time.Now().Add(expirationMinutes)

	// Create claims with user data
	claims := JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Name:     user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "go-clean-arch",
		},
	}

	// Create and sign token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.config.JWTSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}
