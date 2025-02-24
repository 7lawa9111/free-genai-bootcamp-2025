package repository

import (
	"database/sql"
	"time"

	"lang-portal/backend_go/internal/models"
)

type FlashcardRepository struct {
	db *sql.DB
}

func NewFlashcardRepository(db *sql.DB) *FlashcardRepository {
	return &FlashcardRepository{db: db}
}

func (r *FlashcardRepository) CreateActivity(groupID int64, direction string) (*models.FlashcardActivity, error) {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create study activity
	result, err := tx.Exec(`
		INSERT INTO study_activities (group_id, type, created_at, settings)
		VALUES (?, 'flashcards', ?, ?)`,
		groupID, time.Now(), direction)
	if err != nil {
		return nil, err
	}

	activityID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Get words from the group
	rows, err := tx.Query(`
		SELECT w.id, w.japanese, w.english, w.romaji
		FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		WHERE wg.group_id = ?
		ORDER BY RANDOM()
		LIMIT 20`, // Get 20 random words for flashcards
		groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []models.Flashcard
	for rows.Next() {
		var card models.Flashcard
		if err := rows.Scan(&card.ID, &card.Japanese, &card.English, &card.Romaji); err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.FlashcardActivity{
		ID:        activityID,
		GroupID:   groupID,
		Cards:     cards,
		Direction: direction,
	}, nil
}

func (r *FlashcardRepository) SaveResult(result *models.FlashcardResult) error {
	// Start transaction
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

	// Update activity with confidence score
	_, err = tx.Exec(`
		UPDATE study_activities
		SET confidence_score = ?, completed_at = ?
		WHERE id = ?`,
		result.Confidence, time.Now(), result.ActivityID)
	if err != nil {
		return err
	}

	// Record review status for each card
	// Note: In a real implementation, you might want to track individual card results
	_, err = tx.Exec(`
		INSERT INTO word_review_items (study_session_id, word_id, correct, created_at)
		SELECT ?, w.id, ?, ?
		FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		JOIN study_activities sa ON wg.group_id = sa.group_id
		WHERE sa.id = ?
		LIMIT ?`,
		sessionID, true, time.Now(), result.ActivityID, result.CardsReviewed)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *FlashcardRepository) GetActivityStats(activityID int64) (*models.FlashcardResult, error) {
	var result models.FlashcardResult
	err := r.db.QueryRow(`
		SELECT 
			a.id as activity_id,
			COUNT(DISTINCT wri.word_id) as cards_reviewed,
			(SELECT COUNT(DISTINCT w.id)
			 FROM words w
			 JOIN words_groups wg ON w.id = wg.word_id
			 WHERE wg.group_id = a.group_id) as total_cards,
			COALESCE(a.confidence_score, 0) as confidence
		FROM study_activities a
		LEFT JOIN study_sessions s ON a.id = s.study_activity_id
		LEFT JOIN word_review_items wri ON s.id = wri.study_session_id
		WHERE a.id = ? AND a.type = 'flashcards'
		GROUP BY a.id`,
		activityID).Scan(&result.ActivityID, &result.CardsReviewed, &result.TotalCards, &result.Confidence)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &result, nil
} 