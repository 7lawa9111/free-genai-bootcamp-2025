package services

import (
	"time"

	"lang-portal/backend_go/internal/models"
	"lang-portal/backend_go/internal/repository"
)

type DashboardService struct {
	studySessionRepo *repository.StudySessionRepository
	wordRepo         *repository.WordRepository
	groupRepo        *repository.GroupRepository
}

func (s *DashboardService) GetStats() (*models.DashboardStats, error) {
	totalSessions, err := s.studySessionRepo.GetTotalCount()
	if err != nil {
		return nil, err
	}
	// ... implement rest of stats gathering
	return &models.DashboardStats{
		TotalSessions: totalSessions,
	}, nil
}

func NewDashboardService(
	studySessionRepo *repository.StudySessionRepository,
	wordRepo *repository.WordRepository,
	groupRepo *repository.GroupRepository,
) *DashboardService {
	return &DashboardService{
		studySessionRepo: studySessionRepo,
		wordRepo:         wordRepo,
		groupRepo:        groupRepo,
	}
}

type LastStudySessionResponse struct {
	ID              int64     `json:"id"`
	GroupID         int64     `json:"group_id"`
	CreatedAt       time.Time `json:"created_at"`
	StudyActivityID int64     `json:"study_activity_id"`
	GroupName       string    `json:"group_name"`
}

type StudyProgressResponse struct {
	TotalWordsStudied    int `json:"total_words_studied"`
	TotalAvailableWords int `json:"total_available_words"`
}

type QuickStatsResponse struct {
	SuccessRate        float64 `json:"success_rate"`
	TotalStudySessions int     `json:"total_study_sessions"`
	TotalActiveGroups  int     `json:"total_active_groups"`
	StudyStreakDays    int     `json:"study_streak_days"`
}

func (s *DashboardService) GetLastStudySession() (*LastStudySessionResponse, error) {
	session, err := s.studySessionRepo.GetLatest()
	if err != nil {
		return nil, err
	}

	group, err := s.groupRepo.GetByID(session.GroupID)
	if err != nil {
		return nil, err
	}

	return &LastStudySessionResponse{
		ID:              session.ID,
		GroupID:         session.GroupID,
		CreatedAt:       session.CreatedAt,
		StudyActivityID: session.StudyActivityID,
		GroupName:       group.Name,
	}, nil
}

func (s *DashboardService) GetStudyProgress() (*StudyProgressResponse, error) {
	// TODO: Implement actual counting logic from repositories
	// For now, return placeholder data
	return &StudyProgressResponse{
		TotalWordsStudied:    3,
		TotalAvailableWords: 124,
	}, nil
}

func (s *DashboardService) GetQuickStats() (*QuickStatsResponse, error) {
	totalSessions, err := s.studySessionRepo.GetTotalCount()
	if err != nil {
		return nil, err
	}
	// ... rest of the implementation
	return &QuickStatsResponse{
		TotalStudySessions: totalSessions,
		// ... other fields
	}, nil
}

func (s *DashboardService) calculateSuccessRate() (float64, error) {
	// TODO: Implement actual success rate calculation
	// SELECT (COUNT(CASE WHEN correct = true THEN 1 END) * 100.0 / COUNT(*))
	// FROM word_review_items
	return 80.0, nil
}

func (s *DashboardService) calculateStudyStreak() (int, error) {
	// TODO: Implement study streak calculation
	// This would involve checking consecutive days with study sessions
	return 4, nil
}

func (s *DashboardService) getActiveGroupsCount() (int, error) {
	// TODO: Implement active groups count
	// This could be groups with study sessions in the last X days
	return 3, nil
} 