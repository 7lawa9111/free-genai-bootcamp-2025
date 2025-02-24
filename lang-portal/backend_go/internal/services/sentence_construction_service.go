package services

import (
	"lang-portal/backend_go/internal/errors"
	"lang-portal/backend_go/internal/models"
	"lang-portal/backend_go/internal/repository"
	"lang-portal/backend_go/internal/validation"
)

type SentenceConstructionService struct {
	repo *repository.SentenceConstructionRepository
}

func NewSentenceConstructionService(repo *repository.SentenceConstructionRepository) *SentenceConstructionService {
	return &SentenceConstructionService{repo: repo}
}

func (s *SentenceConstructionService) CreateActivity(groupID int64) (*models.SentenceConstruction, error) {
	if err := validation.ValidateID(groupID); err != nil {
		return nil, err
	}

	activity, err := s.repo.CreateActivity(groupID)
	if err != nil {
		return nil, errors.ErrDatabaseError
	}

	// Validate we got enough sentences
	if len(activity.Sentences) < 3 { // Minimum sentences for a meaningful activity
		return nil, errors.ErrInvalidInput
	}

	// Validate each sentence has required fields and scrambled words
	for _, sentence := range activity.Sentences {
		if sentence.Japanese == "" || sentence.English == "" {
			return nil, errors.ErrInvalidInput
		}
		if len(sentence.Words) < 2 { // Need at least 2 words to construct
			return nil, errors.ErrInvalidInput
		}
	}

	return activity, nil
}

func (s *SentenceConstructionService) SaveResult(result *models.SentenceResult) error {
	if err := validation.ValidateID(result.ActivityID); err != nil {
		return err
	}

	// Validate result data
	if err := s.validateResult(result); err != nil {
		return err
	}

	if err := s.repo.SaveResult(result); err != nil {
		return errors.ErrDatabaseError
	}

	return nil
}

func (s *SentenceConstructionService) GetActivityStats(activityID int64) (*models.SentenceResult, error) {
	if err := validation.ValidateID(activityID); err != nil {
		return nil, err
	}

	stats, err := s.repo.GetActivityStats(activityID)
	if err == repository.ErrNotFound {
		return nil, errors.ErrActivityNotFound
	}
	if err != nil {
		return nil, errors.ErrDatabaseError
	}

	return stats, nil
}

// Helper functions

func (s *SentenceConstructionService) validateResult(result *models.SentenceResult) error {
	if result.TimeTaken <= 0 {
		return errors.ErrInvalidInput
	}

	if result.CompletedCount < 0 || result.TotalSentences < result.CompletedCount {
		return errors.ErrInvalidInput
	}

	if result.TotalSentences == 0 {
		return errors.ErrInvalidInput
	}

	if result.Accuracy < 0 || result.Accuracy > 100 {
		return errors.ErrInvalidInput
	}

	// Get activity to verify it exists and is a sentence construction activity
	stats, err := s.repo.GetActivityStats(result.ActivityID)
	if err == repository.ErrNotFound {
		return errors.ErrActivityNotFound
	}
	if err != nil {
		return errors.ErrDatabaseError
	}

	// Verify the total sentences match what we expect
	if stats.TotalSentences != result.TotalSentences {
		return errors.ErrInvalidInput
	}

	return nil
}

// Additional helper methods for activity-specific logic
func (s *SentenceConstructionService) IsActivityComplete(activityID int64) (bool, error) {
	stats, err := s.GetActivityStats(activityID)
	if err != nil {
		return false, err
	}

	return stats.CompletedCount == stats.TotalSentences, nil
}

func (s *SentenceConstructionService) GetProgress(activityID int64) (float64, error) {
	stats, err := s.GetActivityStats(activityID)
	if err != nil {
		return 0, err
	}

	if stats.TotalSentences == 0 {
		return 0, nil
	}

	return float64(stats.CompletedCount) / float64(stats.TotalSentences), nil
} 