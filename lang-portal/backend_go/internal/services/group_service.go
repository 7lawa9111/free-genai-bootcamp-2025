package services

import (
	"fmt"
	"lang-portal/backend_go/internal/errors"
	"lang-portal/backend_go/internal/models"
	"lang-portal/backend_go/internal/repository"
	"lang-portal/backend_go/internal/validation"
)

type GroupService struct {
	groupRepo        *repository.GroupRepository
	wordRepo         *repository.WordRepository
	studySessionRepo *repository.StudySessionRepository
}

func NewGroupService(
	groupRepo *repository.GroupRepository,
	wordRepo *repository.WordRepository,
	studySessionRepo *repository.StudySessionRepository,
) *GroupService {
	return &GroupService{
		groupRepo:        groupRepo,
		wordRepo:         wordRepo,
		studySessionRepo: studySessionRepo,
	}
}

type GroupResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	WordCount int    `json:"word_count"`
}

type GroupDetailsResponse struct {
	ID        int64 `json:"id"`
	Name      string `json:"name"`
	Stats     struct {
		TotalWordCount int `json:"total_word_count"`
	} `json:"stats"`
}

type GroupWordResponse struct {
	Japanese     string `json:"japanese"`
	Romaji       string `json:"romaji"`
	English      string `json:"english"`
	CorrectCount int    `json:"correct_count"`
	WrongCount   int    `json:"wrong_count"`
}

func (s *GroupService) GetGroups(page, limit int) ([]GroupResponse, int, error) {
	// Validate pagination parameters
	page, limit, err := validation.ValidatePagination(page, limit)
	if err != nil {
		return nil, 0, err
	}

	groups, err := s.groupRepo.GetAll(page, limit)
	if err != nil {
		return nil, 0, errors.ErrDatabaseError
	}

	totalCount, err := s.groupRepo.GetTotalCount()
	if err != nil {
		return nil, 0, errors.ErrDatabaseError
	}

	var response []GroupResponse
	for _, group := range groups {
		wordCount, err := s.groupRepo.GetWordCount(group.ID)
		if err != nil {
			return nil, 0, errors.ErrDatabaseError
		}

		response = append(response, GroupResponse{
			ID:        group.ID,
			Name:      group.Name,
			WordCount: wordCount,
		})
	}

	return response, totalCount, nil
}

func (s *GroupService) GetGroupByID(id int64) (*GroupDetailsResponse, error) {
	// Validate ID
	if err := validation.ValidateID(id); err != nil {
		return nil, err
	}

	group, err := s.groupRepo.GetByID(id)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, errors.ErrGroupNotFound
		}
		return nil, errors.ErrDatabaseError
	}

	wordCount, err := s.groupRepo.GetWordCount(group.ID)
	if err != nil {
		return nil, errors.ErrDatabaseError
	}

	return &GroupDetailsResponse{
		ID:   group.ID,
		Name: group.Name,
		Stats: struct {
			TotalWordCount int `json:"total_word_count"`
		}{
			TotalWordCount: wordCount,
		},
	}, nil
}

func (s *GroupService) GetGroupWords(groupID int64, page, limit int) ([]models.Word, int, error) {
	// Validate inputs
	if err := validation.ValidateID(groupID); err != nil {
		return nil, 0, err
	}

	// Get validated page and limit values
	var err error
	page, limit, err = validation.ValidatePagination(page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Get total count first
	totalCount, err := s.groupRepo.GetGroupWordCount(groupID)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting word count: %v", err)
	}

	// Get words
	words, err := s.groupRepo.GetGroupWords(groupID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting words: %v", err)
	}

	return words, totalCount, nil
}

func (s *GroupService) GetGroupStudySessions(groupID int64, page, limit int) ([]models.StudySession, int, error) {
	// Validate inputs
	if err := validation.ValidateID(groupID); err != nil {
		return nil, 0, err
	}

	page, limit, err := validation.ValidatePagination(page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Check if group exists
	if _, err := s.groupRepo.GetByID(groupID); err != nil {
		if err == repository.ErrNotFound {
			return nil, 0, errors.ErrGroupNotFound
		}
		return nil, 0, errors.ErrDatabaseError
	}

	sessions, err := s.groupRepo.GetGroupStudySessions(groupID, page, limit)
	if err != nil {
		return nil, 0, errors.ErrDatabaseError
	}

	totalCount, err := s.groupRepo.GetStudySessionCount(groupID)
	if err != nil {
		return nil, 0, errors.ErrDatabaseError
	}

	return sessions, totalCount, nil
} 