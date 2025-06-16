package adapter

import (
	"context"
	"errors"
	"go-clean-arch/internal/entity"
	"sync"
)

var ErrNotFound = errors.New("not found")

type MemorySessionRepo struct {
	mu       sync.RWMutex
	sessions map[string]*entity.Session
}

func NewMemorySessionRepo() *MemorySessionRepo {
	return &MemorySessionRepo{
		sessions: make(map[string]*entity.Session),
	}
}

func (r *MemorySessionRepo) CreateSession(ctx context.Context, session *entity.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[session.ID] = session
	return nil
}

func (r *MemorySessionRepo) GetSessionByID(ctx context.Context, sessionID string) (*entity.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.sessions[sessionID]
	if !ok {
		return nil, ErrNotFound
	}
	return s, nil
}

func (r *MemorySessionRepo) AddParticipantToSession(ctx context.Context, sessionID string, participant *entity.Participant) error {
	// Not implemented for demo
	return nil
}

func (r *MemorySessionRepo) ListSessions(ctx context.Context) ([]*entity.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []*entity.Session
	for _, s := range r.sessions {
		out = append(out, s)
	}
	return out, nil
}

// ---

type MemoryQuestionRepo struct {
	mu        sync.RWMutex
	questions map[string][]*entity.Question // sessionID -> questions
}

func NewMemoryQuestionRepo() *MemoryQuestionRepo {
	return &MemoryQuestionRepo{
		questions: make(map[string][]*entity.Question),
	}
}

func (r *MemoryQuestionRepo) CreateQuestion(ctx context.Context, q *entity.Question) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.questions[q.SessionID] = append(r.questions[q.SessionID], q)
	return nil
}

func (r *MemoryQuestionRepo) ListQuestionsBySession(ctx context.Context, sessionID string) ([]*entity.Question, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]*entity.Question{}, r.questions[sessionID]...), nil
}
