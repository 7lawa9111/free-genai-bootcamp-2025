package repository

import (
	"database/sql"
	"encoding/json"
	"time"

	"lang-portal/backend_go/internal/models"
)

type WritingPracticeRepository struct {
	db *sql.DB
}

func NewWritingPracticeRepository(db *sql.DB) *WritingPracticeRepository {
	return &WritingPracticeRepository{db: db}
}

func (r *WritingPracticeRepository) CreateActivity(groupID int64) (*models.WritingPractice, error) {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create study activity
	result, err := tx.Exec(`
		INSERT INTO study_activities (group_id, type, created_at)
		VALUES (?, 'writing_practice', ?)`,
		groupID, time.Now())
	if err != nil {
		return nil, err
	}

	activityID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Get words from the group with their hints
	rows, err := tx.Query(`
		SELECT 
			w.id, 
			w.japanese, 
			w.english, 
			w.romaji,
			w.parts as hints
		FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		WHERE wg.group_id = ?
		ORDER BY RANDOM()
		LIMIT 10`, // Get 10 random words for writing practice
		groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []models.WriteExercise
	for rows.Next() {
		var exercise models.WriteExercise
		var hintsJSON []byte
		if err := rows.Scan(&exercise.ID, &exercise.Japanese, &exercise.English, &exercise.Romaji, &hintsJSON); err != nil {
			return nil, err
		}
		// Parse hints from JSON
		if err := json.Unmarshal(hintsJSON, &exercise.Hints); err != nil {
			exercise.Hints = []string{} // Default to empty hints if parsing fails
		}
		exercises = append(exercises, exercise)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.WritingPractice{
		ID:        activityID,
		GroupID:   groupID,
		Exercises: exercises,
	}, nil
}

func (r *WritingPracticeRepository) SaveResult(result *models.WritingResult) error {
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

	// Use sessionID for recording writing results
	_, err = tx.Exec(`
		INSERT INTO writing_practice_results (study_session_id, word_id, written_text, accuracy)
		VALUES (?, ?, ?, ?)`,
		sessionID, result.WordID, result.WrittenText, result.Accuracy)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *WritingPracticeRepository) GetActivityStats(activityID int64) (*models.WritingResult, error) {
	var result models.WritingResult
	err := r.db.QueryRow(`
		SELECT 
			a.id as activity_id,
			COUNT(DISTINCT wri.word_id) as exercises_done,
			(SELECT COUNT(DISTINCT w.id)
			 FROM words w
			 JOIN words_groups wg ON w.id = wg.word_id
			 WHERE wg.group_id = a.group_id) as total_exercises,
			COALESCE(a.accuracy_score, 0) as accuracy
		FROM study_activities a
		LEFT JOIN study_sessions s ON a.id = s.study_activity_id
		LEFT JOIN word_review_items wri ON s.id = wri.study_session_id
		WHERE a.id = ? AND a.type = 'writing_practice'
		GROUP BY a.id`,
		activityID).Scan(&result.ActivityID, &result.ExercisesDone, &result.TotalExercises, &result.Accuracy)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &result, nil
} 