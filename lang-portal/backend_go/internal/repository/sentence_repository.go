package repository

import (
	"lang-portal/backend_go/internal/models"
	"lang-portal/backend_go/internal/errors"
	"database/sql"
)

type SentenceRepository struct {
	db *sql.DB
}

func NewSentenceRepository(db *sql.DB) *SentenceRepository {
	return &SentenceRepository{db: db}
}

func (r *SentenceRepository) GetByID(id int64) (*models.Sentence, error) {
	sentence := &models.Sentence{}
	err := r.db.QueryRow(`
		SELECT id, japanese, english, words, hints
		FROM sentences 
		WHERE id = ?`, id).Scan(&sentence.ID, &sentence.Japanese, &sentence.English, &sentence.Words, &sentence.Hints)
	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return sentence, nil
} 