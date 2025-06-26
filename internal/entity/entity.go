package entity

import (
	"fmt"
	"time"
)

// Session represents a single Q&A event, like a conference talk or a meeting.
// 每個 Session 代表一個獨立的問答活動，例如一場演講或會議。
type Session struct {
	// ID is the system's unique identifier for the session (e.g., a UUID).
	// This is the primary key.
	// ID 是系統中此 Session 的唯一識別碼（例如 UUID），作為主鍵（Primary Key）。
	ID string

	// Title is the descriptive name of the session.
	// Title 是此 Session 的描述性標題。
	Title string

	// IsAcceptingQuestions is a flag to control if new questions can be submitted.
	// IsAcceptingQuestions 是一個旗標，用來控制是否開放提交新問題。
	IsAcceptingQuestions bool
}

// Question represents a single question submitted by a participant.
// Question 代表由參與者提交的單一問題。
type Question struct {
	// ID is the system's unique identifier for the question (e.g., a UUID).
	// This is the primary key.
	// ID 是系統中此 Question 的唯一識別碼（例如 UUID），作為主鍵（Primary Key）。
	ID string

	// SessionID is the foreign key linking this question to a Session.
	// SessionID 是將此問題關聯到一個 Session 的外鍵（Foreign Key）。
	SessionID string

	// Text is the actual content of the question.
	// Text 是問題的實際內容。
	Text string

	// AuthorNickname is the display name of the participant who submitted the question.
	// AuthorNickname 是提交此問題的參與者顯示的暱稱。
	AuthorNickname string

	// Upvotes is the total number of upvotes this question has received.
	// Upvotes 是此問題收到的「附議」總數。
	Upvotes int

	// CreatedAt is the timestamp when the question was submitted.
	// CreatedAt 是問題被提交時的時間戳。
	CreatedAt time.Time

	// UpvotedBy maps participant IDs to their nicknames for this question's upvotes.
	// UpvotedBy 將參與者 ID 對應到他們的暱稱，以記錄對此問題的附議情況。
	UpvotedBy map[string]string
}

func (q *Question) String() string {
	return fmt.Sprintf("[Partcipant: %s] %s", q.AuthorNickname, q.Text)
}

// Participant represents a user connected to a session.
// In this model, a participant is ephemeral and tied to a connection.
// Participant 代表一位連接到 Session 的使用者。在此模型中，參與者是暫時性的，並與一個連線綁定。
type Participant struct {
	// ID is the system's unique identifier for the participant (e.g., a WebSocket connection ID).
	// This is the primary key.
	// ID 是系統中此 Participant 的唯一識別碼（例如 WebSocket 連線 ID），作為主鍵（Primary Key）。
	ID string

	// SessionID is the foreign key linking this participant to a Session.
	// SessionID 是將此參與者關聯到一個 Session 的外鍵（Foreign Key）。
	SessionID string

	// Nickname is the display name chosen by the participant.
	// Nickname 是參與者選擇的顯示名稱。
	Nickname string

	// UpvotedQuestions stores a set of Question IDs that this participant has upvoted.
	// The map's key is the Question ID, and the boolean value indicates presence (true).
	// UpvotedQuestions 儲存了這位參與者已經「附議」過的問題 ID 集合。
	// Map 的鍵是問題 ID，布林值 true 代表存在於集合中。
	UpvotedQuestions map[string]bool

	// CreatedAt is the timestamp when the participant was created.
	// CreatedAt 是參與者被創建時的時間戳。
	CreatedAt time.Time

	// LastSeenAt is the timestamp when the participant was last active.
	// LastSeenAt 是參與者最後活躍的時間戳。
	LastSeenAt time.Time
}

// AuthToken represents a JWT token for participant authentication
type AuthToken struct {
	// Token is the JWT string
	Token string `json:"token"`

	// ExpiresAt is when the token expires
	ExpiresAt time.Time `json:"expires_at"`

	// ParticipantID is the ID of the participant this token belongs to
	ParticipantID string `json:"participant_id"`
}

// AuthClaims represents the claims in a JWT token
type AuthClaims struct {
	ParticipantID string `json:"participant_id"`
	SessionID     string `json:"session_id"`
	Nickname      string `json:"nickname"`
}
