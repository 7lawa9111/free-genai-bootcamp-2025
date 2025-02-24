package models

type ListeningComprehension struct {
	ID        int64       `json:"id"`
	GroupID   int64       `json:"group_id"`
	Exercises []Listening `json:"exercises"`
}

type Listening struct {
	ID         int64    `json:"id"`
	Japanese   string   `json:"japanese"`    // Japanese text
	English    string   `json:"english"`     // English translation
	AudioURL   string   `json:"audio_url"`   // URL to audio file
	UserAnswer string   `json:"user_answer,omitempty"`
	Correct    bool     `json:"correct,omitempty"`
	Hints      []string `json:"hints"`       // Optional hints
}

type ListeningResult struct {
	ActivityID      int64   `json:"activity_id"`
	CompletedCount  int     `json:"completed_count"`
	TotalExercises  int     `json:"total_exercises"`
	TimeTaken       int     `json:"time_taken_seconds"`
	Accuracy        float64 `json:"accuracy"` // Percentage of correct answers
} 