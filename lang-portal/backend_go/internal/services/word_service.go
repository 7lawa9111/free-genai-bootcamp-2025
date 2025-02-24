package services

import (
	"lang-portal/backend_go/internal/models"
	"lang-portal/backend_go/internal/repository"
	"lang-portal/backend_go/internal/validation"
	"lang-portal/backend_go/internal/errors"
)

type WordService struct {
	wordRepo *repository.WordRepository
}

func NewWordService(wordRepo *repository.WordRepository) *WordService {
	return &WordService{wordRepo: wordRepo}
}

type WordDetails struct {
	ID           int64  `json:"id"`
	Japanese     string `json:"japanese"`
	Romaji       string `json:"romaji"`
	English      string `json:"english"`
	Parts        string `json:"parts,omitempty"`
	CorrectCount int    `json:"correct_count"`
	WrongCount   int    `json:"wrong_count"`
}

type WordResponse struct {
	ID       int64   `json:"id"`
	Japanese string  `json:"japanese"`
	Romaji   string  `json:"romaji"`
	English  string  `json:"english"`
	Parts    *string `json:"parts,omitempty"`
}

func (s *WordService) GetWords(page, limit int) ([]models.Word, int, error) {
	words, err := s.wordRepo.GetAll(page, limit)
	if err != nil {
		return nil, 0, err
	}

	totalCount, err := s.wordRepo.GetTotalCount()
	if err != nil {
		return nil, 0, err
	}

	return words, totalCount, nil
}

func (s *WordService) GetWordByID(id int64) (*WordResponse, error) {
	// Validate ID
	if err := validation.ValidateID(id); err != nil {
		return nil, err
	}

	word, err := s.wordRepo.GetByID(id)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, errors.ErrWordNotFound
		}
		return nil, errors.ErrDatabaseError
	}

	return &WordResponse{
		ID:       word.ID,
		Japanese: word.Japanese,
		Romaji:   word.Romaji,
		English:  word.English,
		Parts:    word.Parts,  // Now we can use the pointer directly
	}, nil
} 