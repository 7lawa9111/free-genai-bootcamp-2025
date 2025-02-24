package repository

import (
	"database/sql"
	"encoding/json"
	"time"

	"lang-portal/backend_go/internal/models"
)

type SentenceConstructionRepository struct {
	db *sql.DB
}

func NewSentenceConstructionRepository(db *sql.DB) *SentenceConstructionRepository {
	return &SentenceConstructionRepository{db: db}
}

func (r *SentenceConstructionRepository) CreateActivity(groupID int64) (*models.SentenceConstruction, error) {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create study activity
	result, err := tx.Exec(`
		INSERT INTO study_activities (group_id, type, created_at)
		VALUES (?, 'sentence_construction', ?)`,
		groupID, time.Now())
	if err != nil {
		return nil, err
	}

	activityID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Get sentences from the group
	rows, err := tx.Query(`
		SELECT 
			s.id, 
			s.japanese, 
			s.english,
			s.words,
			s.hints
		FROM sentences s
		JOIN sentences_groups sg ON s.id = sg.sentence_id
		WHERE sg.group_id = ?
		ORDER BY RANDOM()
		LIMIT 5`, // Get 5 random sentences for practice
		groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sentences []models.Sentence
	for rows.Next() {
		var sentence models.Sentence
		var wordsJSON, hintsJSON []byte
		if err := rows.Scan(&sentence.ID, &sentence.Japanese, &sentence.English, &wordsJSON, &hintsJSON); err != nil {
			return nil, err
		}
		// Parse words and hints from JSON
		if err := json.Unmarshal(wordsJSON, &sentence.Words); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(hintsJSON, &sentence.Hints); err != nil {
			sentence.Hints = []string{} // Default to empty hints if parsing fails
		}
		sentences = append(sentences, sentence)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.SentenceConstruction{
		ID:        activityID,
		GroupID:   groupID,
		Sentences: sentences,
	}, nil
}

func (r *SentenceConstructionRepository) SaveResult(result *models.SentenceResult) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create study session for this activity
	sqlResult, err := tx.Exec(`
		INSERT INTO study_sessions (group_id, study_activity_id, created_at)
		SELECT group_id, ?, ?
		FROM study_activities
		WHERE id = ?`,
		result.ActivityID, time.Now(), result.ActivityID)
	if err != nil {
		return err
	}

	sessionID, err := sqlResult.LastInsertId()
	if err != nil {
		return err
	}

	// Update activity with accuracy score
	_, err = tx.Exec(`
		UPDATE study_activities
		SET accuracy_score = ?, completed_at = ?
		WHERE id = ?`,
		result.Accuracy, time.Now(), result.ActivityID)
	if err != nil {
		return err
	}

	// Record individual sentence results
	_, err = tx.Exec(`
		INSERT INTO sentence_review_items (study_session_id, sentence_id, correct, created_at)
		SELECT ?, s.id, ?, ?
		FROM sentences s
		JOIN sentences_groups sg ON s.id = sg.sentence_id
		JOIN study_activities sa ON sg.group_id = sa.group_id
		WHERE sa.id = ?
		LIMIT ?`,
		sessionID, true, time.Now(), result.ActivityID, result.CompletedCount)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *SentenceConstructionRepository) GetActivityStats(activityID int64) (*models.SentenceResult, error) {
	var result models.SentenceResult
	err := r.db.QueryRow(`
		SELECT 
			a.id as activity_id,
			COUNT(DISTINCT sri.sentence_id) as completed_count,
			(SELECT COUNT(DISTINCT s.id)
			 FROM sentences s
			 JOIN sentences_groups sg ON s.id = sg.sentence_id
			 WHERE sg.group_id = a.group_id) as total_sentences,
			COALESCE(a.accuracy_score, 0) as accuracy
		FROM study_activities a
		LEFT JOIN study_sessions s ON a.id = s.study_activity_id
		LEFT JOIN sentence_review_items sri ON s.id = sri.study_session_id
		WHERE a.id = ? AND a.type = 'sentence_construction'
		GROUP BY a.id`,
		activityID).Scan(&result.ActivityID, &result.CompletedCount, &result.TotalSentences, &result.Accuracy)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &result, nil
} 