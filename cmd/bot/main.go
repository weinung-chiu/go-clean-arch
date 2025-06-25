package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type Session struct {
	ID string `json:"id"`
}

type SessionsResponse struct {
	Data  []Session   `json:"data"`
	Error interface{} `json:"error"`
}

var questions = []string{
	"What is the topic?",
	"Can you explain more?",
	"How does this work?",
	"What are the next steps?",
	"Can you give an example?",
}

func getSessions(apiBase string) ([]Session, error) {
	resp, err := http.Get(apiBase + "/api/v1/sessions/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var sessionsResp SessionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&sessionsResp); err != nil {
		return nil, err
	}
	return sessionsResp.Data, nil
}

type QuestionRequest struct {
	Text           string `json:"text"`
	AuthorNickname string `json:"nickname"`
}

func submitQuestion(apiBase, sessionID, question string) error {
	url := fmt.Sprintf("%s/api/v1/sessions/%s/questions", apiBase, sessionID)
	qReq := QuestionRequest{
		Text:           question,
		AuthorNickname: "BotUser",
	}
	bodyBytes, err := json.Marshal(qReq)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", strings.NewReader(string(bodyBytes)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to submit question: %s", resp.Status)
	}
	return nil
}

type Question struct {
	ID             string            `json:"id"`
	SessionID      string            `json:"session_id"`
	Text           string            `json:"text"`
	AuthorNickname string            `json:"author_nickname"`
	Upvotes        int               `json:"upvotes"`
	UpvotedBy      map[string]string `json:"upvotedBy"`
}

type QuestionsResponse struct {
	Data  map[string][]Question `json:"data"`
	Error interface{}           `json:"error"`
}

type QuestionsListResponse struct {
	Data struct {
		Questions []Question `json:"questions"`
	} `json:"data"`
	Error interface{} `json:"error"`
}

func getQuestions(apiBase, sessionID string) ([]Question, error) {
	resp, err := http.Get(apiBase + "/api/v1/sessions/" + sessionID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var questionsResp QuestionsListResponse
	if err := json.NewDecoder(resp.Body).Decode(&questionsResp); err != nil {
		return nil, err
	}
	return questionsResp.Data.Questions, nil
}

func upvoteQuestion(apiBase, sessionID, questionID, botName string) error {
	url := fmt.Sprintf("%s/api/v1/questions/%s/upvote", apiBase, questionID)
	body := map[string]string{"participant_id": botName, "nickname": botName}
	bodyBytes, _ := json.Marshal(body)
	resp, err := http.Post(url, "application/json", strings.NewReader(string(bodyBytes)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to upvote question: %s", resp.Status)
	}
	return nil
}

func main() {
	apiBase := os.Getenv("API_BASE")
	if apiBase == "" {
		apiBase = "http://localhost"
	}
	interval := flag.Int("interval", 30, "Interval in seconds between questions")
	limit := flag.Int("limit", 30, "Maximum number of questions to submit per session")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	botName := "BotUser" + time.Now().Format(time.TimeOnly)

	for range *limit {
		sessions, err := getSessions(apiBase)
		if err != nil {
			fmt.Println("Error fetching sessions:", err)
			time.Sleep(10 * time.Second)
			continue
		}
		for _, s := range sessions {
			qList, err := getQuestions(apiBase, s.ID)
			if err != nil {
				fmt.Printf("Error fetching questions for session %s: %v\n", s.ID, err)
				continue
			}
			for _, q := range qList {
				if q.AuthorNickname == botName {
					continue // Don't upvote own questions
				}
				err := upvoteQuestion(apiBase, s.ID, q.ID, botName)
				if err != nil {
					fmt.Printf("Error upvoting question %s: %v\n", q.ID, err)
				} else {
					fmt.Printf("Upvoted question %s in session %s\n", q.ID, s.ID)
				}
			}
			q := questions[rand.Intn(len(questions))]
			err = submitQuestion(apiBase, s.ID, q)
			if err != nil {
				fmt.Printf("Error submitting question to session %s: %v\n", s.ID, err)
			} else {
				fmt.Printf("Submitted question to session %s: %s\n", s.ID, q)
			}
		}
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}
