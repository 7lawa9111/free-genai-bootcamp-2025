package repository

import (
	"database/sql"
	"time"

	"lang-portal/backend_go/internal/models"
)

type WordMatchingRepository struct {
	db *sql.DB
}

func NewWordMatchingRepository(db *sql.DB) *WordMatchingRepository {
	return &WordMatchingRepository{db: db}
}

func (r *WordMatchingRepository) CreateActivity(groupID int64) (*models.WordMatchingActivity, error) {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create study activity
	result, err := tx.Exec(`
		INSERT INTO study_activities (group_id, type, created_at)
		VALUES (?, 'word_matching', ?)`,
		groupID, time.Now())
	if err != nil {
		return nil, err
	}

	activityID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Get random word pairs from the group
	rows, err := tx.Query(`
		SELECT w.id, w.japanese, w.english
		FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		WHERE wg.group_id = ?
		ORDER BY RANDOM()
		LIMIT 10`, // Get 10 random words for matching
		groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wordPairs []models.WordPair
	for rows.Next() {
		var pair models.WordPair
		if err := rows.Scan(&pair.ID, &pair.Japanese, &pair.English); err != nil {
			return nil, err
		}
		wordPairs = append(wordPairs, pair)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.WordMatchingActivity{
		ID:        activityID,
		GroupID:   groupID,
		WordPairs: wordPairs,
	}, nil
}

func (r *WordMatchingRepository) SaveResult(result *models.WordMatchingResult) error {
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

	// Record word review results
	_, err = tx.Exec(`
		INSERT INTO word_review_items (word_id, study_session_id, correct, created_at)
		SELECT w.id, ?, ?, ?
		FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		JOIN study_activities sa ON wg.group_id = sa.group_id
		WHERE sa.id = ?`,
		sessionID, result.Correct, time.Now(), result.ActivityID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *WordMatchingRepository) GetActivityStats(activityID int64) (*models.WordMatchingResult, error) {
	var result models.WordMatchingResult
	err := r.db.QueryRow(`
		SELECT 
			a.id as activity_id,
			COUNT(DISTINCT CASE WHEN wri.correct = 1 THEN wri.word_id END) as matched_pairs,
			COUNT(DISTINCT wri.word_id) as total_pairs,
			COALESCE(AVG(CASE WHEN wri.correct = 1 THEN 1 ELSE 0 END), 0) > 0.8 as correct
		FROM study_activities a
		LEFT JOIN study_sessions s ON a.id = s.study_activity_id
		LEFT JOIN word_review_items wri ON s.id = wri.study_session_id
		WHERE a.id = ? AND a.type = 'word_matching'
		GROUP BY a.id`,
		activityID).Scan(&result.ActivityID, &result.MatchedPairs, &result.TotalPairs, &result.Correct)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &result, nil
} 