package models

type DashboardStats struct {
    TotalSessions     int     `json:"total_sessions"`
    TotalWords        int     `json:"total_words"`
    CompletionRate    float64 `json:"completion_rate"`
    AverageScore      float64 `json:"average_score"`
} 