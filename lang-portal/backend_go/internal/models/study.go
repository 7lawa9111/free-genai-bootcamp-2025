package models

import "time"

type StudySession struct {
	ID              int       `json:"id"`
	GroupID         int       `json:"group_id"`
	CreatedAt       time.Time `json:"created_at"`
	StudyActivityID int       `json:"study_activity_id"`
	ActivityName    string    `json:"activity_name,omitempty"`
	GroupName       string    `json:"group_name,omitempty"`
	ReviewItemCount int       `json:"review_items_count,omitempty"`
}

type StudyActivity struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ThumbnailURL string `json:"thumbnail_url"`
	Description  string `json:"description"`
}

type WordReviewItem struct {
	WordID         int       `json:"word_id"`
	StudySessionID int       `json:"study_session_id"`
	Correct        bool      `json:"correct"`
	CreatedAt      time.Time `json:"created_at"`
} 