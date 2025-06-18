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

func main() {
	apiBase := os.Getenv("API_BASE")
	if apiBase == "" {
		apiBase = "http://localhost"
	}
	interval := flag.Int("interval", 30, "Interval in seconds between questions")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	for {
		sessions, err := getSessions(apiBase)
		if err != nil {
			fmt.Println("Error fetching sessions:", err)
			time.Sleep(10 * time.Second)
			continue
		}
		for _, s := range sessions {
			q := questions[rand.Intn(len(questions))]
			err := submitQuestion(apiBase, s.ID, q)
			if err != nil {
				fmt.Printf("Error submitting question to session %s: %v\n", s.ID, err)
			} else {
				fmt.Printf("Submitted question to session %s: %s\n", s.ID, q)
			}
		}
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}
