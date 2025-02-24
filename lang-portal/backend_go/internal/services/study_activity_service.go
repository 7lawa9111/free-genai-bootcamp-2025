package services

import (
	"lang-portal/backend_go/internal/errors"
	"lang-portal/backend_go/internal/models"
	"lang-portal/backend_go/internal/repository"
	"lang-portal/backend_go/internal/validation"
)

type StudyActivityService struct {
	repo *repository.StudyActivityRepository
}

func NewStudyActivityService(repo *repository.StudyActivityRepository) *StudyActivityService {
	return &StudyActivityService{repo: repo}
}

func (s *StudyActivityService) GetActivityByID(id int64) (*models.StudyActivityDetails, error) {
	if err := validation.ValidateID(id); err != nil {
		return nil, err
	}

	details, err := s.repo.GetActivityStats(id)
	if err == repository.ErrNotFound {
		return nil, errors.ErrActivityNotFound
	}
	if err != nil {
		return nil, errors.ErrDatabaseError
	}

	return details, nil
}

func (s *StudyActivityService) GetActivitySessions(activityID int64, page, limit int) ([]models.StudySession, int, error) {
	if err := validation.ValidateID(activityID); err != nil {
		return nil, 0, err
	}

	page, limit, err := validation.ValidatePagination(page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Verify activity exists
	if _, err := s.repo.GetByID(activityID); err != nil {
		if err == repository.ErrNotFound {
			return nil, 0, errors.ErrActivityNotFound
		}
		return nil, 0, errors.ErrDatabaseError
	}

	sessions, err := s.repo.GetActivitySessions(activityID, page, limit)
	if err != nil {
		return nil, 0, errors.ErrDatabaseError
	}

	totalCount, err := s.repo.GetSessionCount(activityID)
	if err != nil {
		return nil, 0, errors.ErrDatabaseError
	}

	return sessions, totalCount, nil
}

func (s *StudyActivityService) CreateStudyActivity(groupID int64, activityType string) (*models.StudyActivity, error) {
	if err := validation.ValidateID(groupID); err != nil {
		return nil, err
	}

	// Validate activity type
	if !isValidActivityType(activityType) {
		return nil, errors.ErrInvalidInput
	}

	activity, err := s.repo.Create(groupID, activityType)
	if err != nil {
		return nil, errors.ErrDatabaseError
	}

	return activity, nil
}

// Helper function to validate activity types
func isValidActivityType(activityType string) bool {
	validTypes := map[string]bool{
		"vocabulary_quiz": true,
		"word_matching":   true,
		"flash_cards":     true,
		// Add more activity types as needed
	}
	return validTypes[activityType]
}

// Additional helper functions for activity-specific logic
func (s *StudyActivityService) GetActivityProgress(activityID int64) (float64, error) {
	if err := validation.ValidateID(activityID); err != nil {
		return 0, err
	}

	details, err := s.repo.GetActivityStats(activityID)
	if err == repository.ErrNotFound {
		return 0, errors.ErrActivityNotFound
	}
	if err != nil {
		return 0, errors.ErrDatabaseError
	}

	return details.CompletionRate, nil
}

func (s *StudyActivityService) IsActivityComplete(activityID int64) (bool, error) {
	progress, err := s.GetActivityProgress(activityID)
	if err != nil {
		return false, err
	}

	return progress >= 1.0, nil
} 