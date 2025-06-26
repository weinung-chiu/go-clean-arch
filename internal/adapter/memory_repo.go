package adapter

import (
	"context"
	"errors"
	"go-clean-arch/internal/entity"
	"sync"
	"time"
)

var ErrNotFound = errors.New("not found")

// Merged in-memory repository for sessions, questions, and participants
// Implements SessionRepository, QuestionRepository, ParticipantRepository

type MemoryRepo struct {
	mu sync.RWMutex
	// Sessions
	sessions map[string]*entity.Session
	// Questions
	questions     map[string][]*entity.Question // sessionID -> questions
	questionsByID map[string]*entity.Question   // questionID -> question
	// Participants
	participants map[string]*entity.Participant // participantID -> participant
	byNickname   map[string]*entity.Participant // nickname -> participant
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		sessions:      make(map[string]*entity.Session),
		questions:     make(map[string][]*entity.Question),
		questionsByID: make(map[string]*entity.Question),
		participants:  make(map[string]*entity.Participant),
		byNickname:    make(map[string]*entity.Participant),
	}
}

// --- SessionRepository methods ---

func (r *MemoryRepo) CreateSession(ctx context.Context, session *entity.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[session.ID] = session
	return nil
}

func (r *MemoryRepo) GetSessionByID(ctx context.Context, sessionID string) (*entity.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.sessions[sessionID]
	if !ok {
		return nil, ErrNotFound
	}
	return s, nil
}

func (r *MemoryRepo) AddParticipantToSession(ctx context.Context, sessionID string, participant *entity.Participant) error {
	// Not implemented for demo
	return nil
}

func (r *MemoryRepo) ListSessions(ctx context.Context) ([]*entity.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []*entity.Session
	for _, s := range r.sessions {
		out = append(out, s)
	}
	return out, nil
}

// --- QuestionRepository methods ---

func (r *MemoryRepo) CreateQuestion(ctx context.Context, q *entity.Question) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.questions[q.SessionID] = append(r.questions[q.SessionID], q)
	r.questionsByID[q.ID] = q
	return nil
}

func (r *MemoryRepo) ListQuestionsBySession(ctx context.Context, sessionID string) ([]*entity.Question, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]*entity.Question{}, r.questions[sessionID]...), nil
}

func (r *MemoryRepo) GetQuestionByID(ctx context.Context, questionID string) (*entity.Question, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	q, ok := r.questionsByID[questionID]
	if !ok {
		return nil, ErrNotFound
	}
	return q, nil
}

func (r *MemoryRepo) UpvoteQuestionByID(ctx context.Context, questionID, participantID, participantNickname string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	q, ok := r.questionsByID[questionID]
	if !ok {
		return false, ErrNotFound
	}
	if q.AuthorNickname == participantNickname {
		return false, errors.New("cannot upvote your own question")
	}
	if q.UpvotedBy == nil {
		q.UpvotedBy = make(map[string]string)
	}
	if _, already := q.UpvotedBy[participantID]; already {
		return false, errors.New("already upvoted")
	}
	q.Upvotes++
	q.UpvotedBy[participantID] = participantNickname
	return true, nil
}

// --- ParticipantRepository methods ---

func (r *MemoryRepo) CreateParticipant(ctx context.Context, p *entity.Participant) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.byNickname[p.Nickname]; exists {
		return errors.New("nickname already taken")
	}
	r.participants[p.ID] = p
	r.byNickname[p.Nickname] = p
	return nil
}

func (r *MemoryRepo) GetParticipant(ctx context.Context, filter entity.ParticipantFilter) (*entity.Participant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if filter.ID != nil {
		if p, ok := r.participants[*filter.ID]; ok {
			return p, nil
		}
		return nil, ErrNotFound
	}
	if filter.Nickname != nil {
		if p, ok := r.byNickname[*filter.Nickname]; ok {
			return p, nil
		}
		return nil, ErrNotFound
	}
	return nil, ErrNotFound
}

func (r *MemoryRepo) ListParticipants(ctx context.Context, filter entity.ParticipantFilter) ([]*entity.Participant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var participants []*entity.Participant
	for _, p := range r.participants {
		participants = append(participants, p)
	}
	return participants, nil
}

func (r *MemoryRepo) UpdateParticipantLastSeen(ctx context.Context, participantID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p, ok := r.participants[participantID]; ok {
		p.LastSeenAt = time.Now()
		return nil
	}
	return ErrNotFound
}
