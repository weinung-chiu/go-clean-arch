package api

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-clean-arch/internal/usecase"
	"log"
	"net/http"
	"time"
)

func respondWithError(c *gin.Context, err error) {
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
		var req struct {
			Text     string `json:"text"`
			Nickname string `json:"nickname"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, SessionAPIResp{Error: "invalid request body"})
			return
		}
		question, err := app.SubmitQuestion(c.Request.Context(), sessionID, req.Text, req.Nickname)
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

		eventCh, _ := app.SubscribeToQuestionUpdatedEvents(ctx, sessionID)
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
