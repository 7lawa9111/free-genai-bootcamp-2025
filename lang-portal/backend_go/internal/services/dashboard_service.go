package services

import (
	"database/sql"
	"github.com/mohawa/lang-portal/backend_go/internal/database"
	"github.com/mohawa/lang-portal/backend_go/internal/models"
)

type DashboardService struct {
	db *sql.DB
}

func NewDashboardService() *DashboardService {
	return &DashboardService{db: database.DB}
}

func (s *DashboardService) GetLastStudySession() (*models.StudySession, error) {
	var session models.StudySession
	err := s.db.QueryRow(`
		SELECT 
			ss.id,
			ss.group_id,
			ss.created_at,
			ss.study_activity_id,
			g.name as group_name
		FROM study_sessions ss
		JOIN groups g ON ss.group_id = g.id
		ORDER BY ss.created_at DESC
		LIMIT 1
	`).Scan(
		&session.ID,
		&session.GroupID,
		&session.CreatedAt,
		&session.StudyActivityID,
		&session.GroupName,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *DashboardService) GetStudyProgress() (*models.StudyProgress, error) {
	var progress models.StudyProgress

	err := s.db.QueryRow(`
		SELECT COUNT(DISTINCT id) FROM words
	`).Scan(&progress.TotalAvailableWords)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRow(`
		SELECT COUNT(DISTINCT word_id) 
		FROM word_review_items
	`).Scan(&progress.TotalWordsStudied)
	if err != nil {
		return nil, err
	}

	return &progress, nil
}

func (s *DashboardService) GetQuickStats() (*models.DashboardStats, error) {
	var stats models.DashboardStats

	err := s.db.QueryRow(`
		SELECT 
			COALESCE(
				CAST(SUM(CASE WHEN correct = 1 THEN 1 ELSE 0 END) AS FLOAT) /
				NULLIF(COUNT(*), 0) * 100,
				0
			)
		FROM word_review_items
	`).Scan(&stats.SuccessRate)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM study_sessions
	`).Scan(&stats.TotalStudySessions)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRow(`
		SELECT COUNT(DISTINCT group_id) 
		FROM study_sessions
	`).Scan(&stats.TotalActiveGroups)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRow(`
		WITH RECURSIVE dates(date) AS (
			SELECT date(MAX(created_at)) FROM study_sessions
			UNION ALL
			SELECT date(date, '-1 day')
			FROM dates
			WHERE EXISTS (
				SELECT 1 FROM study_sessions 
				WHERE date(created_at) = date(dates.date, '-1 day')
			)
		)
		SELECT COUNT(*) FROM dates
	`).Scan(&stats.StudyStreakDays)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// ... rest of the implementation ... 