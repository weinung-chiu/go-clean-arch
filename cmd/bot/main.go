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

// API Response structures
type Session struct {
	ID                   string `json:"id"`
	Title                string `json:"title"`
	IsAcceptingQuestions bool   `json:"is_accepting_questions"`
}

type SessionsResponse struct {
	Data  []Session   `json:"data"`
	Error interface{} `json:"error"`
}

type AuthResponse struct {
	Token         string    `json:"token"`
	ExpiresAt     time.Time `json:"expires_at"`
	ParticipantID string    `json:"participant_id"`
	Nickname      string    `json:"nickname"`
	SessionID     string    `json:"session_id"`
}

type AuthAPIResponse struct {
	Data  AuthResponse `json:"data"`
	Error interface{}  `json:"error"`
}

type Question struct {
	ID             string            `json:"id"`
	SessionID      string            `json:"session_id"`
	Text           string            `json:"text"`
	AuthorNickname string            `json:"author_nickname"`
	Upvotes        int               `json:"upvotes"`
	UpvotedBy      map[string]string `json:"upvoted_by"`
}

type QuestionsListResponse struct {
	Data struct {
		Questions []Question `json:"questions"`
	} `json:"data"`
	Error interface{} `json:"error"`
}

type GenericAPIResponse struct {
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}

// Bot configuration
type BotConfig struct {
	APIBase  string
	Interval int
	Limit    int
	BotName  string
	Token    string
}

// Bot client with authentication
type BotClient struct {
	config BotConfig
	client *http.Client
}

func NewBotClient(config BotConfig) *BotClient {
	return &BotClient{
		config: config,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Authentication methods
func (b *BotClient) authenticate(sessionID string) error {
	// Try to register first
	err := b.register(sessionID)
	if err != nil {
		// If registration fails (e.g., nickname taken), try to login
		err = b.login(sessionID)
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}
	}
	return nil
}

func (b *BotClient) register(sessionID string) error {
	url := fmt.Sprintf("%s/api/v1/sessions/%s/register", b.config.APIBase, sessionID)
	reqBody := map[string]string{"nickname": b.config.BotName}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	resp, err := b.client.Post(url, "application/json", strings.NewReader(string(bodyBytes)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("registration failed with status: %d", resp.StatusCode)
	}

	var authResp AuthAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return err
	}

	b.config.Token = authResp.Data.Token
	fmt.Printf("Bot registered successfully for session %s\n", sessionID)
	return nil
}

func (b *BotClient) login(sessionID string) error {
	url := fmt.Sprintf("%s/api/v1/sessions/%s/login", b.config.APIBase, sessionID)
	reqBody := map[string]string{"nickname": b.config.BotName}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	resp, err := b.client.Post(url, "application/json", strings.NewReader(string(bodyBytes)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status: %d", resp.StatusCode)
	}

	var authResp AuthAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return err
	}

	b.config.Token = authResp.Data.Token
	fmt.Printf("Bot logged in successfully for session %s\n", sessionID)
	return nil
}

// API methods with authentication
func (b *BotClient) getSessions() ([]Session, error) {
	resp, err := b.client.Get(b.config.APIBase + "/api/v1/sessions/")
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

func (b *BotClient) getQuestions(sessionID string) ([]Question, error) {
	resp, err := b.client.Get(fmt.Sprintf("%s/api/v1/sessions/%s", b.config.APIBase, sessionID))
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

func (b *BotClient) submitQuestion(sessionID, question string) error {
	url := fmt.Sprintf("%s/api/v1/sessions/%s/questions", b.config.APIBase, sessionID)
	reqBody := map[string]string{"text": question}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(bodyBytes)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+b.config.Token)

	resp, err := b.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to submit question: %s", resp.Status)
	}
	return nil
}

func (b *BotClient) upvoteQuestion(questionID string) error {
	url := fmt.Sprintf("%s/api/v1/questions/%s/upvote", b.config.APIBase, questionID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+b.config.Token)

	resp, err := b.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to upvote question: %s", resp.Status)
	}
	return nil
}

// Bot behavior
func (b *BotClient) run() {
	questions := []string{
		"What is the topic?",
		"Can you explain more?",
		"How does this work?",
		"What are the next steps?",
		"Can you give an example?",
		"What are the benefits?",
		"How is this different from other approaches?",
		"What challenges might we face?",
		"Can you provide more details?",
		"What's the timeline for this?",
	}

	questionCount := 0

	for questionCount < b.config.Limit {
		// Get all sessions
		sessions, err := b.getSessions()
		if err != nil {
			fmt.Printf("Error fetching sessions: %v\n", err)
			time.Sleep(10 * time.Second)
			continue
		}

		if len(sessions) == 0 {
			fmt.Println("No sessions available, waiting...")
			time.Sleep(time.Duration(b.config.Interval) * time.Second)
			continue
		}

		// Process each session
		for _, session := range sessions {
			if !session.IsAcceptingQuestions {
				fmt.Printf("Session %s is not accepting questions, skipping\n", session.ID)
				continue
			}

			// Authenticate for this session
			if err := b.authenticate(session.ID); err != nil {
				fmt.Printf("Failed to authenticate for session %s: %v\n", session.ID, err)
				continue
			}

			// Get questions for this session
			questionsList, err := b.getQuestions(session.ID)
			if err != nil {
				fmt.Printf("Error fetching questions for session %s: %v\n", session.ID, err)
				continue
			}

			// Upvote other people's questions
			for _, q := range questionsList {
				if q.AuthorNickname == b.config.BotName {
					continue // Don't upvote own questions
				}

				err := b.upvoteQuestion(q.ID)
				if err != nil {
					fmt.Printf("Error upvoting question %s: %v\n", q.ID, err)
				} else {
					fmt.Printf("Upvoted question '%s' in session %s\n", q.Text[:min(len(q.Text), 30)], session.ID)
				}
			}

			// Submit a random question
			randomQuestion := questions[rand.Intn(len(questions))]
			err = b.submitQuestion(session.ID, randomQuestion)
			if err != nil {
				fmt.Printf("Error submitting question to session %s: %v\n", session.ID, err)
			} else {
				fmt.Printf("Submitted question to session %s: %s\n", session.ID, randomQuestion)
				questionCount++
			}
		}

		fmt.Printf("Completed cycle %d/%d, waiting %d seconds...\n", questionCount, b.config.Limit, b.config.Interval)
		time.Sleep(time.Duration(b.config.Interval) * time.Second)
	}

	fmt.Println("Bot completed its task!")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	apiBase := os.Getenv("API_BASE")
	if apiBase == "" {
		apiBase = "http://localhost"
	}

	interval := flag.Int("interval", 30, "Interval in seconds between question cycles")
	limit := flag.Int("limit", 10, "Maximum number of questions to submit across all sessions")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	botName := "BotUser" + time.Now().Format("150405") // HHMMSS format

	config := BotConfig{
		APIBase:  apiBase,
		Interval: *interval,
		Limit:    *limit,
		BotName:  botName,
	}

	bot := NewBotClient(config)
	fmt.Printf("Starting bot '%s' with API base: %s\n", botName, apiBase)
	fmt.Printf("Will submit up to %d questions with %d second intervals\n", *limit, *interval)

	bot.run()
}
