package repository

import (
	"database/sql"
	"lang-portal/backend_go/internal/models"
)

type WordRepository struct {
	db *sql.DB
}

func NewWordRepository(db *sql.DB) *WordRepository {
	return &WordRepository{db: db}
}

func (r *WordRepository) GetAll(page, limit int) ([]models.Word, error) {
	offset := (page - 1) * limit
	rows, err := r.db.Query(`
		SELECT id, japanese, romaji, english, parts 
		FROM words 
		LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []models.Word
	for rows.Next() {
		var word models.Word
		err := rows.Scan(&word.ID, &word.Japanese, &word.Romaji, &word.English, &word.Parts)
		if err != nil {
			return nil, err
		}
		words = append(words, word)
	}
	return words, nil
}

func (r *WordRepository) GetByID(id int64) (*models.Word, error) {
	word := &models.Word{}
	err := r.db.QueryRow(`
		SELECT id, japanese, romaji, english, parts 
		FROM words 
		WHERE id = ?`, id).Scan(&word.ID, &word.Japanese, &word.Romaji, &word.English, &word.Parts)
	if err != nil {
		return nil, err
	}
	return word, nil
}

func (r *WordRepository) GetTotalCount() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM words").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
} 