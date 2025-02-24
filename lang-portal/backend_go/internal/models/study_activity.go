package models

import "time"

type StudyActivity struct {
	ID          int64       `json:"id"`
	GroupID     int64       `json:"group_id"`
	Type        string      `json:"type"`      // e.g., "vocabulary_quiz"
	CreatedAt   time.Time   `json:"created_at"`
	CompletedAt *time.Time  `json:"completed_at,omitempty"`
	Settings    *string     `json:"settings,omitempty"`
}

type StudyActivityDetails struct {
	ID               int64     `json:"id"`
	GroupName        string    `json:"group_name"`
	Type             string    `json:"type"`
	CreatedAt        time.Time `json:"created_at"`
	CompletionRate   float64   `json:"completion_rate"`
	TotalWords       int       `json:"total_words"`
	CompletedWords   int       `json:"completed_words"`
} 