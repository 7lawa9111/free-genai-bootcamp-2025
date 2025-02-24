package models

import "time"

type UserStats struct {
	TotalStudyTime      int     `json:"total_study_time_minutes"`
	TotalActivities     int     `json:"total_activities"`
	CompletedActivities int     `json:"completed_activities"`
	AverageAccuracy     float64 `json:"average_accuracy"`
	WordsLearned        int     `json:"words_learned"`
	LastStudyDate       *time.Time `json:"last_study_date,omitempty"`
}

type ActivityStats struct {
	ActivityType    string    `json:"activity_type"`
	CompletionRate  float64   `json:"completion_rate"`
	AverageAccuracy float64   `json:"average_accuracy"`
	TotalTime       int       `json:"total_time_minutes"`
	LastAttempt     time.Time `json:"last_attempt"`
}

type StudyProgress struct {
	DailyStreak    int                    `json:"daily_streak"`
	WeeklyProgress []DailyProgress        `json:"weekly_progress"`
	ByActivityType map[string]ActivityStats `json:"by_activity_type"`
}

type DailyProgress struct {
	Date      time.Time `json:"date"`
	Minutes   int       `json:"minutes_studied"`
	Activities int      `json:"activities_completed"`
} 