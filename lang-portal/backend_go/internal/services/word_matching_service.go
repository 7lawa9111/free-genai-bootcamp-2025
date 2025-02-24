package services

import (
	"lang-portal/backend_go/internal/errors"
	"lang-portal/backend_go/internal/models"
	"lang-portal/backend_go/internal/repository"
	"lang-portal/backend_go/internal/validation"
)

type WordMatchingService struct {
	repo *repository.WordMatchingRepository
}

func NewWordMatchingService(repo *repository.WordMatchingRepository) *WordMatchingService {
	return &WordMatchingService{repo: repo}
}

func (s *WordMatchingService) CreateActivity(groupID int64) (*models.WordMatchingActivity, error) {
	if err := validation.ValidateID(groupID); err != nil {
		return nil, err
	}

	activity, err := s.repo.CreateActivity(groupID)
	if err != nil {
		return nil, errors.ErrDatabaseError
	}

	// Validate we got enough word pairs
	if len(activity.WordPairs) < 4 { // Minimum pairs for a meaningful activity
		return nil, errors.ErrInvalidInput
	}

	return activity, nil
}

func (s *WordMatchingService) SaveResult(result *models.WordMatchingResult) error {
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

func (s *WordMatchingService) GetActivityStats(activityID int64) (*models.WordMatchingResult, error) {
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

func (s *WordMatchingService) validateResult(result *models.WordMatchingResult) error {
	if result.TimeTaken <= 0 {
		return errors.ErrInvalidInput
	}

	if result.MatchedPairs < 0 || result.TotalPairs < result.MatchedPairs {
		return errors.ErrInvalidInput
	}

	if result.TotalPairs == 0 {
		return errors.ErrInvalidInput
	}

	// Get activity to verify it exists and is a word matching activity
	stats, err := s.repo.GetActivityStats(result.ActivityID)
	if err == repository.ErrNotFound {
		return errors.ErrActivityNotFound
	}
	if err != nil {
		return errors.ErrDatabaseError
	}

	// Verify the total pairs match what we expect
	if stats.TotalPairs != result.TotalPairs {
		return errors.ErrInvalidInput
	}

	return nil
}

// Additional helper methods for activity-specific logic
func (s *WordMatchingService) IsActivityComplete(activityID int64) (bool, error) {
	stats, err := s.GetActivityStats(activityID)
	if err != nil {
		return false, err
	}

	return stats.MatchedPairs == stats.TotalPairs, nil
}

func (s *WordMatchingService) GetProgress(activityID int64) (float64, error) {
	stats, err := s.GetActivityStats(activityID)
	if err != nil {
		return 0, err
	}

	if stats.TotalPairs == 0 {
		return 0, nil
	}

	return float64(stats.MatchedPairs) / float64(stats.TotalPairs), nil
} 