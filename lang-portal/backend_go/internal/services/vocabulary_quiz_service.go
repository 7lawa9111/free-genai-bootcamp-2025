package services

import (
	"fmt"
	"lang-portal/backend_go/internal/errors"
	"lang-portal/backend_go/internal/models"
	"lang-portal/backend_go/internal/repository"
	"lang-portal/backend_go/internal/validation"
	"strings"
)

type VocabularyQuizService struct {
	repo *repository.VocabularyQuizRepository
}

func NewVocabularyQuizService(repo *repository.VocabularyQuizRepository) *VocabularyQuizService {
	return &VocabularyQuizService{repo: repo}
}

func (s *VocabularyQuizService) CreateQuiz(groupID int64) (*models.VocabularyQuiz, error) {
	if err := validation.ValidateID(groupID); err != nil {
		return nil, err
	}

	return s.repo.CreateQuiz(groupID)
}

func (s *VocabularyQuizService) SaveResult(result *models.QuizResult) error {
	if err := validation.ValidateID(result.ActivityID); err != nil {
		return err
	}

	// Validate result data
	if err := s.validateResult(result); err != nil {
		return err
	}

	// Calculate score as percentage
	result.Score = float64(result.CorrectCount) / float64(result.TotalCount) * 100

	if err := s.repo.SaveResult(result); err != nil {
		return errors.ErrDatabaseError
	}

	return nil
}

func (s *VocabularyQuizService) GetQuizStats(activityID int64) (*models.QuizResult, error) {
	if err := validation.ValidateID(activityID); err != nil {
		return nil, err
	}

	stats, err := s.repo.GetQuizStats(activityID)
	if err == repository.ErrNotFound {
		return nil, errors.ErrActivityNotFound
	}
	if err != nil {
		return nil, errors.ErrDatabaseError
	}

	return stats, nil
}

// Helper functions

func (s *VocabularyQuizService) validateResult(result *models.QuizResult) error {
	// Validate IDs
	if result.ActivityID <= 0 {
		return fmt.Errorf("invalid activity ID: %d", result.ActivityID)
	}
	if result.WordID <= 0 {
		return fmt.Errorf("invalid word ID: %d", result.WordID)
	}

	// Validate answer
	if strings.TrimSpace(result.Answer) == "" {
		return fmt.Errorf("answer cannot be empty")
	}

	// Validate time taken
	if result.TimeTaken <= 0 {
		return fmt.Errorf("time taken must be positive")
	}

	return nil
}

// Additional helper methods for quiz-specific logic
func (s *VocabularyQuizService) IsQuizComplete(activityID int64) (bool, error) {
	stats, err := s.GetQuizStats(activityID)
	if err != nil {
		return false, err
	}

	return stats.CorrectCount == stats.TotalCount, nil
}

func (s *VocabularyQuizService) GetProgress(activityID int64) (float64, error) {
	stats, err := s.GetQuizStats(activityID)
	if err != nil {
		return 0, err
	}

	if stats.TotalCount == 0 {
		return 0, nil
	}

	return float64(stats.CorrectCount) / float64(stats.TotalCount), nil
}

func (s *VocabularyQuizService) GetQuizDebug(activityID int64) (map[string]interface{}, error) {
	return s.repo.GetQuizDebug(activityID)
} 