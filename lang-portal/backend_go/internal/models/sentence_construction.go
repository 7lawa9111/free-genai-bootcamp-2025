package models

type SentenceConstruction struct {
	ID        int64       `json:"id"`
	GroupID   int64       `json:"group_id"`
	Sentences []Sentence  `json:"sentences"`
}

type Sentence struct {
	ID         int64    `json:"id"`
	Japanese   string   `json:"japanese"`    // Complete Japanese sentence
	English    string   `json:"english"`     // English translation
	Words      []string `json:"words"`       // Japanese words in scrambled order
	UserAnswer []string `json:"user_answer,omitempty"`
	Correct    bool     `json:"correct,omitempty"`
	Hints      []string `json:"hints"`       // Optional grammar/context hints
}

type SentenceResult struct {
	ActivityID      int64   `json:"activity_id"`
	CompletedCount  int     `json:"completed_count"`
	TotalSentences  int     `json:"total_sentences"`
	TimeTaken       int     `json:"time_taken_seconds"`
	Accuracy        float64 `json:"accuracy"` // Percentage of correct sentences
} 