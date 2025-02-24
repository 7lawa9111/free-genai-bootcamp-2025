package models

type Flashcard struct {
	ID       int64  `json:"id"`
	Japanese string `json:"japanese"`
	English  string `json:"english"`
	Romaji   string `json:"romaji"`
	Revealed bool   `json:"revealed"`
}

type FlashcardActivity struct {
	ID        int64       `json:"id"`
	GroupID   int64       `json:"group_id"`
	Cards     []Flashcard `json:"cards"`
	Direction string      `json:"direction"` // "ja_to_en" or "en_to_ja"
}

type FlashcardResult struct {
	ActivityID    int64   `json:"activity_id"`
	CardsReviewed int     `json:"cards_reviewed"`
	TotalCards    int     `json:"total_cards"`
	TimeTaken     int     `json:"time_taken_seconds"`
	Confidence    float64 `json:"confidence_score"` // User's self-reported confidence 0-1
} 