package repository

import (
	"database/sql"
	"fmt"

	"lang-portal/backend_go/internal/models"
)

type VocabularyQuizRepository struct {
	db *sql.DB
}

func NewVocabularyQuizRepository(db *sql.DB) *VocabularyQuizRepository {
	return &VocabularyQuizRepository{db: db}
}

func (r *VocabularyQuizRepository) CreateQuiz(groupID int64) (*models.VocabularyQuiz, error) {
	// First verify the group exists
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM groups WHERE id = ?)", groupID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("error checking group existence: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("group not found: %d", groupID)
	}

	// Get words for the quiz
	rows, err := r.db.Query(`
		SELECT w.id, w.japanese, w.english, w.romaji
		FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		WHERE wg.group_id = ?
		ORDER BY RANDOM()
		LIMIT 10`, groupID)
	if err != nil {
		return nil, fmt.Errorf("error fetching words: %v", err)
	}
	defer rows.Close()

	var words []models.Word
	for rows.Next() {
		var word models.Word
		err := rows.Scan(&word.ID, &word.Japanese, &word.English, &word.Romaji)
		if err != nil {
			return nil, fmt.Errorf("error scanning word: %v", err)
		}
		words = append(words, word)
	}

	if len(words) < 4 {
		return nil, fmt.Errorf("not enough words in group for quiz (minimum 4 required)")
	}

	// Create the quiz activity
	result, err := r.db.Exec(`
		INSERT INTO study_activities (group_id, activity_type, created_at)
		VALUES (?, 'vocabulary_quiz', DATETIME('now'))`,
		groupID)
	if err != nil {
		return nil, fmt.Errorf("error creating quiz activity: %v", err)
	}

	quizID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting quiz ID: %v", err)
	}

	return &models.VocabularyQuiz{
		ID:        quizID,
		GroupID:   groupID,
		Questions: words,
	}, nil
}

func (r *VocabularyQuizRepository) SaveResult(result *models.QuizResult) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// First verify the activity exists
	var exists bool
	err = tx.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM study_activities 
			WHERE id = ? AND activity_type = 'vocabulary_quiz'
		)`, result.ActivityID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking activity existence: %v", err)
	}
	if !exists {
		return fmt.Errorf("activity not found: %d", result.ActivityID)
	}

	// Get group_id for the study session
	var groupID int64
	err = tx.QueryRow(`
		SELECT group_id FROM study_activities WHERE id = ?`, 
		result.ActivityID).Scan(&groupID)
	if err != nil {
		return fmt.Errorf("error getting group ID: %v", err)
	}

	// Update study activity with score
	_, err = tx.Exec(`
		UPDATE study_activities 
		SET score = ?, completed_at = CURRENT_TIMESTAMP
		WHERE id = ? AND activity_type = 'vocabulary_quiz'`,
		result.Score, result.ActivityID)
	if err != nil {
		return fmt.Errorf("error updating activity score: %v", err)
	}

	// Create study session
	res, err := tx.Exec(`
		INSERT INTO study_sessions (group_id, study_activity_id, created_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)`,
		groupID, result.ActivityID)
	if err != nil {
		return fmt.Errorf("error creating study session: %v", err)
	}

	sessionID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting session ID: %v", err)
	}

	// Save quiz answers
	_, err = tx.Exec(`
		INSERT INTO vocabulary_quiz_answers (
			study_session_id, word_id, answer, correct, created_at
		) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)`,
		sessionID, result.WordID, result.Answer, result.Correct)
	if err != nil {
		return fmt.Errorf("error saving quiz answer: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (r *VocabularyQuizRepository) GetQuizStats(activityID int64) (*models.QuizResult, error) {
	// First verify the activity exists and get basic info
	var result models.QuizResult
	err := r.db.QueryRow(`
		SELECT id, COALESCE(score, 0)
		FROM study_activities
		WHERE id = ? AND activity_type = 'vocabulary_quiz'`,
		activityID).Scan(&result.ActivityID, &result.Score)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("activity not found: %d", activityID)
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching activity: %v", err)
	}

	// Get answer counts in a separate query
	err = r.db.QueryRow(`
		SELECT 
			COUNT(CASE WHEN correct = 1 THEN 1 END),
			COUNT(*)
		FROM vocabulary_quiz_answers
		WHERE study_session_id IN (
			SELECT id FROM study_sessions 
			WHERE study_activity_id = ?
		)`,
		activityID).Scan(&result.CorrectCount, &result.TotalCount)
	
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error fetching answer counts: %v", err)
	}

	// If no answers yet, initialize counts to 0
	if err == sql.ErrNoRows {
		result.CorrectCount = 0
		result.TotalCount = 0
	}

	return &result, nil
}

func (r *VocabularyQuizRepository) GetQuizDebug(activityID int64) (map[string]interface{}, error) {
	debug := make(map[string]interface{})

	// Check activity
	var activity struct {
		ID          int64
		GroupID     int64
		Type        string
		CreatedAt   string
		CompletedAt sql.NullString
		Score       sql.NullFloat64
	}
	err := r.db.QueryRow(`
		SELECT id, group_id, activity_type, created_at, completed_at, score 
		FROM study_activities 
		WHERE id = ?`, activityID).Scan(
			&activity.ID, &activity.GroupID, &activity.Type,
			&activity.CreatedAt, &activity.CompletedAt, &activity.Score)
	if err != nil {
		return nil, fmt.Errorf("error checking activity: %v", err)
	}
	debug["activity"] = activity

	// Check sessions
	rows, err := r.db.Query(`
		SELECT id, created_at, completed_at 
		FROM study_sessions 
		WHERE study_activity_id = ?`, activityID)
	if err != nil {
		return nil, fmt.Errorf("error checking sessions: %v", err)
	}
	defer rows.Close()

	var sessions []interface{}
	for rows.Next() {
		var session struct {
			ID          int64
			CreatedAt   string
			CompletedAt sql.NullString
		}
		if err := rows.Scan(&session.ID, &session.CreatedAt, &session.CompletedAt); err != nil {
			return nil, fmt.Errorf("error scanning session: %v", err)
		}
		sessions = append(sessions, session)
	}
	debug["sessions"] = sessions

	return debug, nil
} 