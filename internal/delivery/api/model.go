package api

import "time"

// SessionAPIResp is the response model for session-related APIs.
type SessionAPIResp struct {
	Data  any `json:"data"`
	Error any `json:"error"`
}

type Session struct {
	ID                   string `json:"id"`
	Title                string `json:"title"`
	IsAcceptingQuestions bool   `json:"is_accepting_questions"`
}

type Question struct {
	ID             string            `json:"id"`
	SessionID      string            `json:"session_id"`
	Text           string            `json:"text"`
	AuthorNickname string            `json:"author_nickname"`
	Upvotes        int               `json:"upvotes"`
	CreatedAt      time.Time         `json:"created_at"`
	UpvotedBy      map[string]string `json:"upvoted_by"` // map of user ID to nickname
}

// Authentication DTOs
type RegisterRequest struct {
	Nickname string `json:"nickname" binding:"required"`
}

type LoginRequest struct {
	Nickname string `json:"nickname" binding:"required"`
}

type AuthResponse struct {
	Token         string    `json:"token"`
	ExpiresAt     time.Time `json:"expires_at"`
	ParticipantID string    `json:"participant_id"`
	Nickname      string    `json:"nickname"`
	SessionID     string    `json:"session_id"`
}

type Participant struct {
	ID               string          `json:"id"`
	SessionID        string          `json:"session_id"`
	Nickname         string          `json:"nickname"`
	UpvotedQuestions map[string]bool `json:"upvoted_questions"`
	CreatedAt        time.Time       `json:"created_at"`
	LastSeenAt       time.Time       `json:"last_seen_at"`
}
