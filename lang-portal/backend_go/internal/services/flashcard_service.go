package services

import (
	"lang-portal/backend_go/internal/errors"
	"lang-portal/backend_go/internal/models"
	"lang-portal/backend_go/internal/repository"
	"lang-portal/backend_go/internal/validation"
)

type FlashcardService struct {
	repo *repository.FlashcardRepository
}

func NewFlashcardService(repo *repository.FlashcardRepository) *FlashcardService {
	return &FlashcardService{repo: repo}
}

func (s *FlashcardService) CreateActivity(groupID int64, direction string) (*models.FlashcardActivity, error) {
	if err := validation.ValidateID(groupID); err != nil {
		return nil, err
	}

	if err := s.validateDirection(direction); err != nil {
		return nil, err
	}

	activity, err := s.repo.CreateActivity(groupID, direction)
	if err != nil {
		return nil, errors.ErrDatabaseError
	}

	// Validate we got enough cards
	if len(activity.Cards) < 5 { // Minimum cards for a meaningful activity
		return nil, errors.ErrInvalidInput
	}

	return activity, nil
}

func (s *FlashcardService) SaveResult(result *models.FlashcardResult) error {
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

func (s *FlashcardService) GetActivityStats(activityID int64) (*models.FlashcardResult, error) {
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

func (s *FlashcardService) validateDirection(direction string) error {
	validDirections := map[string]bool{
		"ja_to_en": true,
		"en_to_ja": true,
	}

	if !validDirections[direction] {
		return errors.ErrInvalidInput
	}

	return nil
}

func (s *FlashcardService) validateResult(result *models.FlashcardResult) error {
	if result.TimeTaken <= 0 {
		return errors.ErrInvalidInput
	}

	if result.CardsReviewed < 0 || result.TotalCards < result.CardsReviewed {
		return errors.ErrInvalidInput
	}

	if result.TotalCards == 0 {
		return errors.ErrInvalidInput
	}

	if result.Confidence < 0 || result.Confidence > 1 {
		return errors.ErrInvalidInput
	}

	// Get activity to verify it exists and is a flashcard activity
	stats, err := s.repo.GetActivityStats(result.ActivityID)
	if err == repository.ErrNotFound {
		return errors.ErrActivityNotFound
	}
	if err != nil {
		return errors.ErrDatabaseError
	}

	// Verify the total cards match what we expect
	if stats.TotalCards != result.TotalCards {
		return errors.ErrInvalidInput
	}

	return nil
}

// Additional helper methods for activity-specific logic
func (s *FlashcardService) IsActivityComplete(activityID int64) (bool, error) {
	stats, err := s.GetActivityStats(activityID)
	if err != nil {
		return false, err
	}

	return stats.CardsReviewed == stats.TotalCards, nil
}

func (s *FlashcardService) GetProgress(activityID int64) (float64, error) {
	stats, err := s.GetActivityStats(activityID)
	if err != nil {
		return 0, err
	}

	if stats.TotalCards == 0 {
		return 0, nil
	}

	return float64(stats.CardsReviewed) / float64(stats.TotalCards), nil
} 