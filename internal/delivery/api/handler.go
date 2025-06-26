package api

import (
	"context"
	"encoding/json"
	"go-clean-arch/internal/entity"
	"go-clean-arch/internal/usecase"
	"go-clean-arch/internal/usecase/service/auth"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func respondWithError(c *gin.Context, err error) {
	// TODO: return correct HTTP status codes based on error type
	c.JSON(http.StatusInternalServerError, SessionAPIResp{Error: err.Error()})
}

// HandlerListSessions handles GET /sessions
func HandlerListSessions(app *usecase.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessions, err := app.ListSessions(c.Request.Context())
		if err != nil {
			respondWithError(c, err)
			return
		}
		var result []Session
		for _, s := range sessions {
			result = append(result, Session{
				ID:                   s.ID,
				Title:                s.Title,
				IsAcceptingQuestions: s.IsAcceptingQuestions,
			})
		}
		c.JSON(http.StatusOK, SessionAPIResp{Data: result})
	}
}

// HandlerNewSession handles POST /sessions
func HandlerNewSession(app *usecase.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Title string `json:"title"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, SessionAPIResp{Error: "invalid request body"})
			return
		}
		session, err := app.NewSession(c.Request.Context(), req.Title)
		if err != nil {
			respondWithError(c, err)
			return
		}
		c.JSON(http.StatusOK, SessionAPIResp{Data: Session{
			ID:                   session.ID,
			Title:                session.Title,
			IsAcceptingQuestions: session.IsAcceptingQuestions,
		}})
	}
}

// HandlerGetSession handles GET /sessions/:id
func HandlerGetSession(app *usecase.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("id")
		session, questions, err := app.GetSession(c.Request.Context(), sessionID)
		if err != nil {
			respondWithError(c, err)
			return
		}
		var qList []Question
		for _, q := range questions {
			qList = append(qList, Question{
				ID:             q.ID,
				SessionID:      q.SessionID,
				Text:           q.Text,
				AuthorNickname: q.AuthorNickname,
				Upvotes:        q.Upvotes,
				CreatedAt:      q.CreatedAt,
			})
		}
		c.JSON(http.StatusOK, SessionAPIResp{Data: gin.H{
			"session": Session{
				ID:                   session.ID,
				Title:                session.Title,
				IsAcceptingQuestions: session.IsAcceptingQuestions,
			},
			"questions": qList,
		}})
	}
}

// HandlerSubmitQuestion handles POST /sessions/:id/questions
func HandlerSubmitQuestion(app *usecase.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("id")

		// Get authenticated participant from context
		participantInterface, exists := c.Get("participant")
		if !exists {
			c.JSON(http.StatusUnauthorized, SessionAPIResp{Error: "Authentication required"})
			return
		}
		participant := participantInterface.(*entity.Participant)

		var req struct {
			Text string `json:"text"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, SessionAPIResp{Error: "invalid request body"})
			return
		}

		question, err := app.SubmitQuestion(c.Request.Context(), sessionID, req.Text, participant.Nickname)
		if err != nil {
			respondWithError(c, err)
			return
		}
		c.JSON(http.StatusOK, SessionAPIResp{Data: Question{
			ID:             question.ID,
			SessionID:      question.SessionID,
			Text:           question.Text,
			AuthorNickname: question.AuthorNickname,
			Upvotes:        question.Upvotes,
			CreatedAt:      question.CreatedAt,
		}})
	}
}

// HandlerWebSocketSession handles WebSocket connections for session events
func HandlerWebSocketSession(app *usecase.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("id")
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("WebSocket upgrade error:", err)
			return
		}
		defer conn.Close()

		ctx, cancel := context.WithCancel(c.Request.Context())
		defer cancel()

		eventCh, _ := app.SubscribeToClientEvent(ctx, sessionID)
		//defer unsubscribe()

		pongWait := 60 * time.Second
		pingPeriod := (pongWait * 9) / 10
		conn.SetReadDeadline(time.Now().Add(pongWait))
		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

		go func() {
			for {
				time.Sleep(pingPeriod)
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					cancel()
					return
				}
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-eventCh:
				if !ok {
					return
				}

				// Marshal the questions list to JSON and send as WS payload
				questions := make([]Question, 0, len(event.Questions))
				for _, q := range event.Questions {
					questions = append(questions, Question{
						ID:             q.ID,
						SessionID:      q.SessionID,
						Text:           q.Text,
						AuthorNickname: q.AuthorNickname,
						Upvotes:        q.Upvotes,
						CreatedAt:      q.CreatedAt,
						UpvotedBy:      q.UpvotedBy,
					})
				}
				payload, err := json.Marshal(gin.H{"questions": questions})
				if err != nil {
					return
				}
				if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
					return
				}
			}
		}
	}
}

// HandlerUpvoteQuestion handles POST /questions/:id/upvote
func HandlerUpvoteQuestion(app *usecase.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		questionID := c.Param("id")

		// Get authenticated participant from context
		participantInterface, exists := c.Get("participant")
		if !exists {
			c.JSON(http.StatusUnauthorized, SessionAPIResp{Error: "Authentication required"})
			return
		}
		participant := participantInterface.(*entity.Participant)

		err := app.UpvoteQuestion(c.Request.Context(), questionID, participant.ID, participant.Nickname)
		if err != nil {
			respondWithError(c, err)
			return
		}
		c.JSON(http.StatusOK, SessionAPIResp{Data: gin.H{"message": "Question upvoted successfully"}})
	}
}

// AuthMiddleware validates JWT tokens and adds participant to context
func AuthMiddleware(app *usecase.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, SessionAPIResp{Error: "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, SessionAPIResp{Error: "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]
		participant, err := app.AuthService.ValidateParticipantToken(c.Request.Context(), tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, SessionAPIResp{Error: "Invalid or expired token"})
			c.Abort()
			return
		}

		// Add participant to context
		c.Set("participant", participant)
		c.Next()
	}
}

// HandlerRegister handles POST /sessions/:id/register
func HandlerRegister(app *usecase.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("id")

		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, SessionAPIResp{Error: "Invalid request body"})
			return
		}

		token, err := app.AuthService.RegisterParticipant(c.Request.Context(), sessionID, req.Nickname)
		if err != nil {
			if err.Error() == auth.ErrNicknameAlreadyTaken {
				c.JSON(http.StatusConflict, SessionAPIResp{Error: "Nickname already taken in this session"})
				return
			}
			respondWithError(c, err)
			return
		}

		c.JSON(http.StatusCreated, SessionAPIResp{Data: AuthResponse{
			Token:         token.Token,
			ExpiresAt:     token.ExpiresAt,
			ParticipantID: token.ParticipantID,
			Nickname:      req.Nickname,
			SessionID:     sessionID,
		}})
	}
}

// HandlerLogin handles POST /sessions/:id/login
func HandlerLogin(app *usecase.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("id")

		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, SessionAPIResp{Error: "Invalid request body"})
			return
		}

		token, err := app.AuthService.LoginParticipant(c.Request.Context(), sessionID, req.Nickname)
		if err != nil {
			if err.Error() == auth.ErrParticipantNotFound {
				c.JSON(http.StatusNotFound, SessionAPIResp{Error: "Participant not found"})
				return
			}
			respondWithError(c, err)
			return
		}

		c.JSON(http.StatusOK, SessionAPIResp{Data: AuthResponse{
			Token:         token.Token,
			ExpiresAt:     token.ExpiresAt,
			ParticipantID: token.ParticipantID,
			Nickname:      req.Nickname,
			SessionID:     sessionID,
		}})
	}
}

// HandlerGetProfile handles GET /auth/profile (requires authentication)
func HandlerGetProfile(app *usecase.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		participantInterface, exists := c.Get("participant")
		if !exists {
			c.JSON(http.StatusUnauthorized, SessionAPIResp{Error: "Authentication required"})
			return
		}

		participant := participantInterface.(*entity.Participant)
		c.JSON(http.StatusOK, SessionAPIResp{Data: Participant{
			ID:               participant.ID,
			SessionID:        participant.SessionID,
			Nickname:         participant.Nickname,
			UpvotedQuestions: participant.UpvotedQuestions,
			CreatedAt:        participant.CreatedAt,
			LastSeenAt:       participant.LastSeenAt,
		}})
	}
}
