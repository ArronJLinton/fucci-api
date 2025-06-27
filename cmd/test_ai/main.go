package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ArronJLinton/fucci-api/internal/ai"
)

// Mock cache for testing
type MockCache struct{}

func (m *MockCache) Get(ctx context.Context, key string, value interface{}) error {
	return fmt.Errorf("cache miss")
}

func (m *MockCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return nil
}

func (m *MockCache) Exists(ctx context.Context, key string) (bool, error) {
	return false, nil
}

func main() {
	// Check for required environment variables
	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	openAIBaseURL := os.Getenv("OPENAI_BASE_URL")
	if openAIBaseURL == "" {
		openAIBaseURL = "https://api.openai.com/v1"
	}

	// Initialize mock cache
	mockCache := &MockCache{}

	// Initialize AI prompt generator
	promptGenerator := ai.NewPromptGenerator(openAIKey, openAIBaseURL, mockCache)

	// Test match data
	testMatchData := ai.MatchData{
		MatchID:  "1321727",
		HomeTeam: "Manchester City",
		AwayTeam: "Liverpool",
		Date:     "2024-01-15T20:00:00Z",
		Status:   "NS", // Not Started
		Lineups: &ai.LineupData{
			HomeStarters: []ai.Player{
				{ID: 1, Name: "Ederson", Number: 31, Pos: "G", Photo: ""},
				{ID: 2, Name: "Kyle Walker", Number: 2, Pos: "D", Photo: ""},
				{ID: 3, Name: "Rúben Dias", Number: 3, Pos: "D", Photo: ""},
				{ID: 4, Name: "John Stones", Number: 5, Pos: "D", Photo: ""},
				{ID: 5, Name: "João Cancelo", Number: 7, Pos: "D", Photo: ""},
				{ID: 6, Name: "Rodri", Number: 16, Pos: "M", Photo: ""},
				{ID: 7, Name: "Kevin De Bruyne", Number: 17, Pos: "M", Photo: ""},
				{ID: 8, Name: "Bernardo Silva", Number: 20, Pos: "M", Photo: ""},
				{ID: 9, Name: "Phil Foden", Number: 47, Pos: "M", Photo: ""},
				{ID: 10, Name: "Jack Grealish", Number: 10, Pos: "F", Photo: ""},
				{ID: 11, Name: "Erling Haaland", Number: 9, Pos: "F", Photo: ""},
			},
			AwayStarters: []ai.Player{
				{ID: 12, Name: "Alisson", Number: 1, Pos: "G", Photo: ""},
				{ID: 13, Name: "Trent Alexander-Arnold", Number: 66, Pos: "D", Photo: ""},
				{ID: 14, Name: "Virgil van Dijk", Number: 4, Pos: "D", Photo: ""},
				{ID: 15, Name: "Joël Matip", Number: 32, Pos: "D", Photo: ""},
				{ID: 16, Name: "Andy Robertson", Number: 26, Pos: "D", Photo: ""},
				{ID: 17, Name: "Fabinho", Number: 3, Pos: "M", Photo: ""},
				{ID: 18, Name: "Jordan Henderson", Number: 14, Pos: "M", Photo: ""},
				{ID: 19, Name: "Thiago Alcântara", Number: 6, Pos: "M", Photo: ""},
				{ID: 20, Name: "Mohamed Salah", Number: 11, Pos: "F", Photo: ""},
				{ID: 21, Name: "Sadio Mané", Number: 10, Pos: "F", Photo: ""},
				{ID: 22, Name: "Diogo Jota", Number: 20, Pos: "F", Photo: ""},
			},
		},
		NewsHeadlines: []string{
			"Manchester City vs Liverpool: Premier League title race heats up",
			"Pep Guardiola praises Liverpool's pressing game",
			"Jurgen Klopp confident ahead of Etihad clash",
			"Erling Haaland in sensational form for City",
			"Mohamed Salah eyes Golden Boot in crucial match",
		},
		SocialSentiment: &ai.SocialSentiment{
			TwitterSentiment:     0.3, // Slightly positive
			RedditSentiment:      0.1, // Neutral
			TopTopics:            []string{"title race", "tactics", "form"},
			ControversialMoments: []string{"referee decisions", "VAR calls"},
		},
	}

	ctx := context.Background()

	fmt.Println("🤖 Testing AI Prompt Generation")
	fmt.Println("===============================")

	// Test pre-match prompt
	fmt.Println("\n📋 Testing Pre-Match Prompt Generation...")
	preMatchPrompt, err := promptGenerator.GeneratePreMatchPrompt(ctx, testMatchData)
	if err != nil {
		log.Printf("❌ Failed to generate pre-match prompt: %v", err)
	} else {
		fmt.Println("✅ Pre-match prompt generated successfully!")
		printPrompt("Pre-Match", preMatchPrompt)
	}

	// Test post-match prompt with different data
	postMatchData := testMatchData
	postMatchData.Status = "FT" // Finished
	postMatchData.Stats = &ai.MatchStats{
		HomeGoals:       2,
		AwayGoals:       1,
		HomeShots:       15,
		AwayShots:       8,
		HomePossession:  65,
		AwayPossession:  35,
		HomeFouls:       12,
		AwayFouls:       14,
		HomeYellowCards: 2,
		AwayYellowCards: 3,
		HomeRedCards:    0,
		AwayRedCards:    0,
	}

	fmt.Println("\n📋 Testing Post-Match Prompt Generation...")
	postMatchPrompt, err := promptGenerator.GeneratePostMatchPrompt(ctx, postMatchData)
	if err != nil {
		log.Printf("❌ Failed to generate post-match prompt: %v", err)
	} else {
		fmt.Println("✅ Post-match prompt generated successfully!")
		printPrompt("Post-Match", postMatchPrompt)
	}

	fmt.Println("\n🎉 Testing completed!")
}

func printPrompt(promptType string, prompt *ai.DebatePrompt) {
	fmt.Printf("\n%s Debate Prompt:\n", promptType)
	fmt.Printf("Headline: %s\n", prompt.Headline)
	fmt.Printf("Description: %s\n", prompt.Description)
	fmt.Printf("Cards (%d):\n", len(prompt.Cards))

	for i, card := range prompt.Cards {
		fmt.Printf("  %d. [%s] %s\n", i+1, card.Stance, card.Title)
		fmt.Printf("     Description: %s\n", card.Description)
	}

	// Pretty print JSON
	jsonData, _ := json.MarshalIndent(prompt, "", "  ")
	fmt.Printf("\nJSON Output:\n%s\n", string(jsonData))
}
