package models

import "time"

type StudySession struct {
	ID              int64     `json:"id"`
	GroupID         int64     `json:"group_id"`
	StudyActivityID int64     `json:"study_activity_id"`
	CreatedAt       time.Time `json:"created_at"`
}

type StudySessionDetails struct {
	ID               int64     `json:"id"`
	ActivityName     string    `json:"activity_name"`
	GroupName        string    `json:"group_name"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	ReviewItemsCount int       `json:"review_items_count"`
}

type WordReview struct {
	WordID    int64     `json:"word_id"`
	SessionID int64     `json:"session_id"`
	Correct   bool      `json:"correct"`
	CreatedAt time.Time `json:"created_at"`
} 