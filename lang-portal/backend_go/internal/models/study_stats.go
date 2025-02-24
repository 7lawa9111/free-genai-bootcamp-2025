package models

type StudyStatsResponse struct {
    TotalSessions     int     `json:"total_sessions"`
    TotalTimeMinutes  float64 `json:"total_time_minutes"`
    AverageScore      float64 `json:"average_score"`
    CompletionRate    float64 `json:"completion_rate"`
} 