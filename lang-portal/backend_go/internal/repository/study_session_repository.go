package repository

import (
	"database/sql"
	"time"

	"lang-portal/backend_go/internal/models"
)

type StudySessionRepository struct {
	db *sql.DB
}

func NewStudySessionRepository(db *sql.DB) *StudySessionRepository {
	return &StudySessionRepository{db: db}
}

func (r *StudySessionRepository) GetAll(page, limit int) ([]models.StudySession, error) {
	offset := (page - 1) * limit
	rows, err := r.db.Query(`
		SELECT id, group_id, study_activity_id, created_at 
		FROM study_sessions 
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.StudySession
	for rows.Next() {
		var session models.StudySession
		if err := rows.Scan(&session.ID, &session.GroupID, &session.StudyActivityID, &session.CreatedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (r *StudySessionRepository) GetByID(id int64) (*models.StudySession, error) {
	session := &models.StudySession{}
	err := r.db.QueryRow(`
		SELECT id, group_id, study_activity_id, created_at 
		FROM study_sessions 
		WHERE id = ?`, id).Scan(&session.ID, &session.GroupID, &session.StudyActivityID, &session.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (r *StudySessionRepository) GetSessionWords(sessionID int64, page, limit int) ([]models.Word, error) {
	offset := (page - 1) * limit
	rows, err := r.db.Query(`
		SELECT w.id, w.japanese, w.romaji, w.english, w.parts
		FROM words w
		JOIN word_review_items wri ON w.id = wri.word_id
		WHERE wri.study_session_id = ?
		ORDER BY wri.created_at
		LIMIT ? OFFSET ?`, sessionID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []models.Word
	for rows.Next() {
		var word models.Word
		if err := rows.Scan(&word.ID, &word.Japanese, &word.Romaji, &word.English, &word.Parts); err != nil {
			return nil, err
		}
		words = append(words, word)
	}

	return words, nil
}

func (r *StudySessionRepository) GetWordCount(sessionID int64) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) 
		FROM word_review_items 
		WHERE study_session_id = ?`, sessionID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *StudySessionRepository) AddWordReview(sessionID, wordID int64, correct bool) error {
	_, err := r.db.Exec(`
		INSERT INTO word_review_items (word_id, study_session_id, correct, created_at)
		VALUES (?, ?, ?, ?)`,
		wordID, sessionID, correct, time.Now())
	return err
}

func (r *StudySessionRepository) GetLatest() (*models.StudySession, error) {
	session := &models.StudySession{}
	err := r.db.QueryRow(`
		SELECT id, group_id, study_activity_id, created_at 
		FROM study_sessions 
		ORDER BY created_at DESC 
		LIMIT 1`).Scan(&session.ID, &session.GroupID, &session.StudyActivityID, &session.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (r *StudySessionRepository) GetTotalCount() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM study_sessions").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *StudySessionRepository) ResetHistory() error {
	_, err := r.db.Exec("DELETE FROM word_review_items")
	if err != nil {
		return err
	}
	_, err = r.db.Exec("DELETE FROM study_sessions")
	return err
}

func (r *StudySessionRepository) Create(groupID, studyActivityID int64) (*models.StudySession, error) {
	result, err := r.db.Exec(`
		INSERT INTO study_sessions (group_id, study_activity_id, created_at)
		VALUES (?, ?, ?)`,
		groupID, studyActivityID, time.Now())
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *StudySessionRepository) GetStats() (*models.StudyStatsResponse, error) {
	stats := &models.StudyStatsResponse{}
	err := r.db.QueryRow(`
		SELECT 
			COUNT(*) as total_sessions,
			COALESCE(AVG(TIMESTAMPDIFF(MINUTE, created_at, completed_at)), 0) as avg_time,
			COALESCE(AVG(CASE WHEN completed_at IS NOT NULL THEN 1 ELSE 0 END), 0) as completion_rate,
			COALESCE(AVG(CASE WHEN score IS NOT NULL THEN score ELSE 0 END), 0) as avg_score
		FROM study_sessions`).Scan(
		&stats.TotalSessions,
		&stats.TotalTimeMinutes,
		&stats.CompletionRate,
		&stats.AverageScore,
	)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *StudySessionRepository) FullReset() error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	tables := []string{
		"word_review_items",
		"study_sessions",
		"study_activities",
		"words_groups",
		"words",
		"groups",
	}

	for _, table := range tables {
		if _, err := tx.Exec("DELETE FROM " + table); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
} 