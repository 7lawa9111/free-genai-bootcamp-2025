package services

import (
	"lang-portal/backend_go/internal/errors"
	"lang-portal/backend_go/internal/models"
	"lang-portal/backend_go/internal/repository"
	"time"
)

type StatisticsService struct {
	repo *repository.StatisticsRepository
}

func NewStatisticsService(repo *repository.StatisticsRepository) *StatisticsService {
	return &StatisticsService{repo: repo}
}

func (s *StatisticsService) GetUserStats() (*models.UserStats, error) {
	stats, err := s.repo.GetUserStats()
	if err != nil {
		return nil, errors.ErrDatabaseError
	}

	return stats, nil
}

func (s *StatisticsService) GetActivityStats() (map[string]models.ActivityStats, error) {
	stats, err := s.repo.GetActivityStats()
	if err != nil {
		return nil, errors.ErrDatabaseError
	}

	// Validate activity types
	validTypes := map[string]bool{
		"word_matching":          true,
		"vocabulary_quiz":        true,
		"flashcards":            true,
		"writing_practice":      true,
		"sentence_construction": true,
	}

	for actType := range stats {
		if !validTypes[actType] {
			delete(stats, actType)
		}
	}

	return stats, nil
}

func (s *StatisticsService) GetStudyProgress() (*models.StudyProgress, error) {
	progress, err := s.repo.GetStudyProgress()
	if err != nil {
		return nil, errors.ErrDatabaseError
	}

	// Ensure weekly progress has entries for all 7 days
	if len(progress.WeeklyProgress) < 7 {
		fullProgress := make([]models.DailyProgress, 7)
		progressMap := make(map[string]models.DailyProgress)

		// Map existing progress by date string
		for _, dp := range progress.WeeklyProgress {
			dateStr := dp.Date.Format("2006-01-02")
			progressMap[dateStr] = dp
		}

		// Fill in missing days with zero values
		for i := 0; i < 7; i++ {
			date := time.Now().AddDate(0, 0, -i)
			dateStr := date.Format("2006-01-02")
			if dp, exists := progressMap[dateStr]; exists {
				fullProgress[i] = dp
			} else {
				fullProgress[i] = models.DailyProgress{
					Date:       date,
					Minutes:   0,
					Activities: 0,
				}
			}
		}

		progress.WeeklyProgress = fullProgress
	}

	return progress, nil
}

// Helper methods for calculating additional statistics
func (s *StatisticsService) GetAverageStudyTime() (float64, error) {
	stats, err := s.repo.GetUserStats()
	if err != nil {
		return 0, errors.ErrDatabaseError
	}

	if stats.TotalActivities == 0 {
		return 0, nil
	}

	return float64(stats.TotalStudyTime) / float64(stats.TotalActivities), nil
}

func (s *StatisticsService) GetCompletionRate() (float64, error) {
	stats, err := s.repo.GetUserStats()
	if err != nil {
		return 0, errors.ErrDatabaseError
	}

	if stats.TotalActivities == 0 {
		return 0, nil
	}

	return float64(stats.CompletedActivities) / float64(stats.TotalActivities) * 100, nil
} 