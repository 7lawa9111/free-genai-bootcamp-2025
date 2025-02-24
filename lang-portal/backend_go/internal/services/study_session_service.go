package services

import (
	"time"

	"lang-portal/backend_go/internal/errors"
	"lang-portal/backend_go/internal/models"
	"lang-portal/backend_go/internal/repository"
	"lang-portal/backend_go/internal/validation"
)

type StudySessionService struct {
	repo *repository.StudySessionRepository
}

func NewStudySessionService(repo *repository.StudySessionRepository) *StudySessionService {
	return &StudySessionService{repo: repo}
}

type StudySessionDetails struct {
	ID               int64     `json:"id"`
	ActivityName     string    `json:"activity_name"`
	GroupName        string    `json:"group_name"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	ReviewItemsCount int       `json:"review_items_count"`
}

func (s *StudySessionService) GetStudySessions(page, limit int) ([]models.StudySession, int, error) {
	page, limit, err := validation.ValidatePagination(page, limit)
	if err != nil {
		return nil, 0, err
	}

	sessions, err := s.repo.GetAll(page, limit)
	if err != nil {
		return nil, 0, errors.ErrDatabaseError
	}

	totalCount, err := s.repo.GetTotalCount()
	if err != nil {
		return nil, 0, errors.ErrDatabaseError
	}

	return sessions, totalCount, nil
}

func (s *StudySessionService) GetStudySessionByID(id int64) (*StudySessionDetails, error) {
	if err := validation.ValidateID(id); err != nil {
		return nil, err
	}

	session, err := s.repo.GetByID(id)
	if err == repository.ErrNotFound {
		return nil, errors.ErrStudySessionNotFound
	}
	if err != nil {
		return nil, errors.ErrDatabaseError
	}

	wordCount, err := s.repo.GetWordCount(id)
	if err != nil {
		return nil, errors.ErrDatabaseError
	}

	// For now, return static activity and group names
	// TODO: Get actual names from their respective repositories
	return &StudySessionDetails{
		ID:               session.ID,
		ActivityName:     "Vocabulary Quiz",
		GroupName:        "Basic Greetings",
		StartTime:        session.CreatedAt,
		EndTime:          session.CreatedAt.Add(10 * time.Minute), // Placeholder duration
		ReviewItemsCount: wordCount,
	}, nil
}

func (s *StudySessionService) GetStudySessionWords(sessionID int64, page, limit int) ([]models.Word, int, error) {
	if err := validation.ValidateID(sessionID); err != nil {
		return nil, 0, err
	}

	page, limit, err := validation.ValidatePagination(page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Verify session exists
	if _, err := s.repo.GetByID(sessionID); err != nil {
		if err == repository.ErrNotFound {
			return nil, 0, errors.ErrStudySessionNotFound
		}
		return nil, 0, errors.ErrDatabaseError
	}

	words, err := s.repo.GetSessionWords(sessionID, page, limit)
	if err != nil {
		return nil, 0, errors.ErrDatabaseError
	}

	wordCount, err := s.repo.GetWordCount(sessionID)
	if err != nil {
		return nil, 0, errors.ErrDatabaseError
	}

	return words, wordCount, nil
}

func (s *StudySessionService) ReviewWord(sessionID, wordID int64, correct bool) error {
	if err := validation.ValidateID(sessionID); err != nil {
		return err
	}
	if err := validation.ValidateID(wordID); err != nil {
		return err
	}

	// Verify session exists
	if _, err := s.repo.GetByID(sessionID); err != nil {
		if err == repository.ErrNotFound {
			return errors.ErrStudySessionNotFound
		}
		return errors.ErrDatabaseError
	}

	if err := s.repo.AddWordReview(sessionID, wordID, correct); err != nil {
		return errors.ErrDatabaseError
	}

	return nil
}

func (s *StudySessionService) ResetHistory() error {
	if err := s.repo.ResetHistory(); err != nil {
		return errors.ErrDatabaseError
	}
	return nil
}

func (s *StudySessionService) GetLatestSession() (*models.StudySession, error) {
	return s.repo.GetLatest()
}

func (s *StudySessionService) CreateStudySession(groupID, studyActivityID int64) (*models.StudySession, error) {
	if err := validation.ValidateID(groupID); err != nil {
		return nil, err
	}
	if err := validation.ValidateID(studyActivityID); err != nil {
		return nil, err
	}

	session, err := s.repo.Create(groupID, studyActivityID)
	if err != nil {
		return nil, errors.ErrDatabaseError
	}

	return session, nil
}

func (s *StudySessionService) GetStudyStats() (*models.StudyStatsResponse, error) {
	stats, err := s.repo.GetStats()
	if err != nil {
		return nil, errors.ErrDatabaseError
	}

	return stats, nil
}

func (s *StudySessionService) FullReset() error {
	if err := s.repo.FullReset(); err != nil {
		return errors.ErrDatabaseError
	}
	return nil
} 