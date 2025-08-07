package adapter

import (
	"context"
	"errors"
	"go-clean-arch/internal/entity"
	"sync"
)

// MemoryUserRepository provides an in-memory implementation of UserRepository
type MemoryUserRepository struct {
	users          map[string]*entity.User // Key: user ID
	userByUsername map[string]*entity.User // Key: username
	mu             sync.RWMutex
}

// NewMemoryUserRepository creates a new in-memory user repository
func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users:          make(map[string]*entity.User),
		userByUsername: make(map[string]*entity.User),
	}
}

// Create stores a new user in memory
func (r *MemoryUserRepository) Create(ctx context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if user ID already exists
	if _, exists := r.users[user.ID]; exists {
		return errors.New("user with this ID already exists")
	}

	// Check if username already exists
	if _, exists := r.userByUsername[user.Username]; exists {
		return errors.New("username already taken")
	}

	// Store user by both ID and username
	r.users[user.ID] = user
	r.userByUsername[user.Username] = user

	return nil
}

// GetByUsername retrieves a user by username
func (r *MemoryUserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.userByUsername[username]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (r *MemoryUserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// Update updates an existing user
func (r *MemoryUserRepository) Update(ctx context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if user exists
	existingUser, exists := r.users[user.ID]
	if !exists {
		return errors.New("user not found")
	}

	// If username is changing, check availability and update index
	if existingUser.Username != user.Username {
		if _, usernameExists := r.userByUsername[user.Username]; usernameExists {
			return errors.New("username already taken")
		}

		// Remove old username mapping
		delete(r.userByUsername, existingUser.Username)
		// Add new username mapping
		r.userByUsername[user.Username] = user
	}

	// Update user
	r.users[user.ID] = user

	return nil
}
