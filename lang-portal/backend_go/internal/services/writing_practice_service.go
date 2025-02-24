package services

import (
	"lang-portal/backend_go/internal/errors"
	"lang-portal/backend_go/internal/models"
	"lang-portal/backend_go/internal/repository"
	"lang-portal/backend_go/internal/validation"
)

type WritingPracticeService struct {
	repo *repository.WritingPracticeRepository
}

func NewWritingPracticeService(repo *repository.WritingPracticeRepository) *WritingPracticeService {
	return &WritingPracticeService{repo: repo}
}

func (s *WritingPracticeService) CreateActivity(groupID int64) (*models.WritingPractice, error) {
	if err := validation.ValidateID(groupID); err != nil {
		return nil, err
	}

	activity, err := s.repo.CreateActivity(groupID)
	if err != nil {
		return nil, errors.ErrDatabaseError
	}

	// Validate we got enough exercises
	if len(activity.Exercises) < 5 { // Minimum exercises for a meaningful activity
		return nil, errors.ErrInvalidInput
	}

	// Validate each exercise has required fields
	for _, ex := range activity.Exercises {
		if ex.Japanese == "" || ex.English == "" {
			return nil, errors.ErrInvalidInput
		}
	}

	return activity, nil
}

func (s *WritingPracticeService) SaveResult(result *models.WritingResult) error {
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

func (s *WritingPracticeService) GetActivityStats(activityID int64) (*models.WritingResult, error) {
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

func (s *WritingPracticeService) validateResult(result *models.WritingResult) error {
	if result.TimeTaken <= 0 {
		return errors.ErrInvalidInput
	}

	if result.ExercisesDone < 0 || result.TotalExercises < result.ExercisesDone {
		return errors.ErrInvalidInput
	}

	if result.TotalExercises == 0 {
		return errors.ErrInvalidInput
	}

	if result.Accuracy < 0 || result.Accuracy > 100 {
		return errors.ErrInvalidInput
	}

	// Get activity to verify it exists and is a writing practice activity
	stats, err := s.repo.GetActivityStats(result.ActivityID)
	if err == repository.ErrNotFound {
		return errors.ErrActivityNotFound
	}
	if err != nil {
		return errors.ErrDatabaseError
	}

	// Verify the total exercises match what we expect
	if stats.TotalExercises != result.TotalExercises {
		return errors.ErrInvalidInput
	}

	return nil
}

// Additional helper methods for activity-specific logic
func (s *WritingPracticeService) IsActivityComplete(activityID int64) (bool, error) {
	stats, err := s.GetActivityStats(activityID)
	if err != nil {
		return false, err
	}

	return stats.ExercisesDone == stats.TotalExercises, nil
}

func (s *WritingPracticeService) GetProgress(activityID int64) (float64, error) {
	stats, err := s.GetActivityStats(activityID)
	if err != nil {
		return 0, err
	}

	if stats.TotalExercises == 0 {
		return 0, nil
	}

	return float64(stats.ExercisesDone) / float64(stats.TotalExercises), nil
} 