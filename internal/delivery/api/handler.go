package api

import (
	"github.com/gin-gonic/gin"
	"go-clean-arch/internal/usecase"
	"net/http"
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
