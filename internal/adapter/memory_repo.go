package adapter

import (
	"context"
	"errors"
	"go-clean-arch/internal/entity"
	"sync"
	"time"
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

// Add a map for direct question lookup by ID

type MemoryQuestionRepo struct {
	mu        sync.RWMutex
	questions map[string][]*entity.Question // sessionID -> questions
	byID      map[string]*entity.Question   // questionID -> question
}

func NewMemoryQuestionRepo() *MemoryQuestionRepo {
	return &MemoryQuestionRepo{
		questions: make(map[string][]*entity.Question),
		byID:      make(map[string]*entity.Question),
	}
}

func (r *MemoryQuestionRepo) CreateQuestion(ctx context.Context, q *entity.Question) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.questions[q.SessionID] = append(r.questions[q.SessionID], q)
	r.byID[q.ID] = q
	return nil
}

func (r *MemoryQuestionRepo) ListQuestionsBySession(ctx context.Context, sessionID string) ([]*entity.Question, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]*entity.Question{}, r.questions[sessionID]...), nil
}

func (r *MemoryQuestionRepo) GetQuestionByID(ctx context.Context, questionID string) (*entity.Question, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	q, ok := r.byID[questionID]
	if !ok {
		return nil, ErrNotFound
	}
	return q, nil
}

func (r *MemoryQuestionRepo) UpvoteQuestionByID(ctx context.Context, questionID, participantID, participantNickname string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	q, ok := r.byID[questionID]
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

// ---

type MemoryParticipantRepo struct {
	mu           sync.RWMutex
	participants map[string]map[string]*entity.Participant // sessionID -> participantID -> participant
	byNickname   map[string]map[string]*entity.Participant // sessionID -> nickname -> participant
}

func NewMemoryParticipantRepo() *MemoryParticipantRepo {
	return &MemoryParticipantRepo{
		participants: make(map[string]map[string]*entity.Participant),
		byNickname:   make(map[string]map[string]*entity.Participant),
	}
}

func (r *MemoryParticipantRepo) CreateParticipant(ctx context.Context, p *entity.Participant) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Initialize maps if they don't exist
	if r.participants[p.SessionID] == nil {
		r.participants[p.SessionID] = make(map[string]*entity.Participant)
	}
	if r.byNickname[p.SessionID] == nil {
		r.byNickname[p.SessionID] = make(map[string]*entity.Participant)
	}

	// Store by ID
	r.participants[p.SessionID][p.ID] = p
	// Store by nickname
	r.byNickname[p.SessionID][p.Nickname] = p

	return nil
}

func (r *MemoryParticipantRepo) GetParticipant(ctx context.Context, filter entity.ParticipantFilter) (*entity.Participant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// If ID is provided, search by ID first (most specific)
	if filter.ID != nil {
		for _, sessionParticipants := range r.participants {
			if p, ok := sessionParticipants[*filter.ID]; ok {
				return p, nil
			}
		}
		return nil, ErrNotFound
	}

	// If SessionID and Nickname are provided, search by both
	if filter.SessionID != nil && filter.Nickname != nil {
		if sessionNicknames, ok := r.byNickname[*filter.SessionID]; ok {
			if p, ok := sessionNicknames[*filter.Nickname]; ok {
				return p, nil
			}
		}
		return nil, ErrNotFound
	}

	// If only SessionID is provided, return first participant in session
	if filter.SessionID != nil {
		if sessionParticipants, ok := r.participants[*filter.SessionID]; ok {
			for _, p := range sessionParticipants {
				return p, nil // Return first participant found
			}
		}
		return nil, ErrNotFound
	}

	// If only Nickname is provided, search across all sessions
	if filter.Nickname != nil {
		for _, sessionNicknames := range r.byNickname {
			if p, ok := sessionNicknames[*filter.Nickname]; ok {
				return p, nil
			}
		}
		return nil, ErrNotFound
	}

	return nil, ErrNotFound
}

func (r *MemoryParticipantRepo) ListParticipants(ctx context.Context, filter entity.ParticipantFilter) ([]*entity.Participant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var participants []*entity.Participant

	// If SessionID is provided, list participants in that session
	if filter.SessionID != nil {
		if sessionParticipants, ok := r.participants[*filter.SessionID]; ok {
			for _, p := range sessionParticipants {
				participants = append(participants, p)
			}
		}
		return participants, nil
	}

	// If no SessionID provided, list all participants across all sessions
	for _, sessionParticipants := range r.participants {
		for _, p := range sessionParticipants {
			participants = append(participants, p)
		}
	}

	return participants, nil
}

func (r *MemoryParticipantRepo) UpdateParticipantLastSeen(ctx context.Context, participantID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Search through all sessions
	for _, sessionParticipants := range r.participants {
		if p, ok := sessionParticipants[participantID]; ok {
			p.LastSeenAt = time.Now()
			return nil
		}
	}

	return ErrNotFound
}

func (r *MemoryParticipantRepo) AddParticipant(ctx context.Context, sessionID string, p *entity.Participant) error {
	return r.CreateParticipant(ctx, p)
}

// UpvoteQuestion updates the upvote count and participant upvote set.
func (r *MemoryQuestionRepo) UpvoteQuestion(ctx context.Context, sessionID, questionID, participantID, participantNickname string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	qs := r.questions[sessionID]
	for _, q := range qs {
		if q.ID == questionID {
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
	}
	return false, ErrNotFound
}
