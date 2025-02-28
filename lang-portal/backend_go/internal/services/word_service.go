package services

import (
	"database/sql"
	"github.com/mohawa/lang-portal/backend_go/internal/database"
	"github.com/mohawa/lang-portal/backend_go/internal/models"
)

type WordService struct {
	db *sql.DB
}

func NewWordService() *WordService {
	return &WordService{db: database.DB}
}

func (s *WordService) GetWords(page, perPage int) (*models.PaginatedResponse, error) {
	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM words").Scan(&total)
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * perPage
	rows, err := s.db.Query(`
		SELECT w.japanese, w.romaji, w.english,
			   COUNT(CASE WHEN wri.correct = 1 THEN 1 END) as correct_count,
			   COUNT(CASE WHEN wri.correct = 0 THEN 1 END) as wrong_count
		FROM words w
		LEFT JOIN word_review_items wri ON w.id = wri.word_id
		GROUP BY w.id
		LIMIT ? OFFSET ?
	`, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []models.WordWithStats
	for rows.Next() {
		var w models.WordWithStats
		if err := rows.Scan(&w.Japanese, &w.Romaji, &w.English, &w.CorrectCount, &w.WrongCount); err != nil {
			return nil, err
		}
		words = append(words, w)
	}

	return &models.PaginatedResponse{
		Items: words,
		Pagination: models.Pagination{
			CurrentPage:  page,
			TotalPages:   (total + perPage - 1) / perPage,
			TotalItems:   total,
			ItemsPerPage: perPage,
		},
	}, nil
}

func (s *WordService) GetWord(id int) (*models.WordResponse, error) {
	var word models.WordResponse
	err := s.db.QueryRow(`
		SELECT w.japanese, w.romaji, w.english,
			   COUNT(CASE WHEN wri.correct = 1 THEN 1 END) as correct_count,
			   COUNT(CASE WHEN wri.correct = 0 THEN 1 END) as wrong_count
		FROM words w
		LEFT JOIN word_review_items wri ON w.id = wri.word_id
		WHERE w.id = ?
		GROUP BY w.id
	`, id).Scan(&word.Japanese, &word.Romaji, &word.English, &word.Stats.CorrectCount, &word.Stats.WrongCount)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(`
		SELECT g.id, g.name
		FROM groups g
		JOIN words_groups wg ON g.id = wg.group_id
		WHERE wg.word_id = ?
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.ID, &group.Name); err != nil {
			return nil, err
		}
		word.Groups = append(word.Groups, group)
	}

	return &word, nil
} 