package repository

import (
	"database/sql"
	"time"

	"lang-portal/backend_go/internal/models"
)

type StudyActivityRepository struct {
	db *sql.DB
}

func NewStudyActivityRepository(db *sql.DB) *StudyActivityRepository {
	return &StudyActivityRepository{db: db}
}

func (r *StudyActivityRepository) GetByID(id int64) (*models.StudyActivity, error) {
	activity := &models.StudyActivity{}
	err := r.db.QueryRow(`
		SELECT id, group_id, type, created_at, completed_at, settings
		FROM study_activities
		WHERE id = ?`, id).Scan(
		&activity.ID,
		&activity.GroupID,
		&activity.Type,
		&activity.CreatedAt,
		&activity.CompletedAt,
		&activity.Settings,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return activity, nil
}

func (r *StudyActivityRepository) GetActivitySessions(activityID int64, page, limit int) ([]models.StudySession, error) {
	offset := (page - 1) * limit
	rows, err := r.db.Query(`
		SELECT id, group_id, study_activity_id, created_at
		FROM study_sessions
		WHERE study_activity_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`,
		activityID, limit, offset)
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

func (r *StudyActivityRepository) Create(groupID int64, activityType string) (*models.StudyActivity, error) {
	result, err := r.db.Exec(`
		INSERT INTO study_activities (group_id, type, created_at)
		VALUES (?, ?, ?)`,
		groupID, activityType, time.Now())
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *StudyActivityRepository) GetTotalCount() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM study_activities").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *StudyActivityRepository) GetSessionCount(activityID int64) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) 
		FROM study_sessions 
		WHERE study_activity_id = ?`, activityID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *StudyActivityRepository) GetActivityStats(activityID int64) (*models.StudyActivityDetails, error) {
	var details models.StudyActivityDetails
	var groupName string

	// Get basic activity info and group name
	err := r.db.QueryRow(`
		SELECT 
			sa.id,
			g.name,
			sa.type,
			sa.created_at
		FROM study_activities sa
		JOIN groups g ON sa.group_id = g.id
		WHERE sa.id = ?`,
		activityID).Scan(&details.ID, &groupName, &details.Type, &details.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	// Get completion statistics
	err = r.db.QueryRow(`
		SELECT 
			COUNT(DISTINCT w.id) as total_words,
			COUNT(DISTINCT CASE WHEN wri.correct = 1 THEN w.id END) as completed_words,
			COALESCE(AVG(CASE WHEN wri.correct = 1 THEN 1.0 ELSE 0.0 END), 0) as completion_rate
		FROM study_sessions ss
		JOIN words_groups wg ON ss.group_id = wg.group_id
		JOIN words w ON wg.word_id = w.id
		LEFT JOIN word_review_items wri ON w.id = wri.word_id AND ss.id = wri.study_session_id
		WHERE ss.study_activity_id = ?`,
		activityID).Scan(&details.TotalWords, &details.CompletedWords, &details.CompletionRate)
	if err != nil {
		return nil, err
	}

	details.GroupName = groupName
	return &details, nil
} 