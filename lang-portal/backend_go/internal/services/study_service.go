package services

import (
	"database/sql"
	"time"
	"github.com/mohawa/lang-portal/backend_go/internal/database"
	"github.com/mohawa/lang-portal/backend_go/internal/models"
)

type StudyService struct {
	db *sql.DB
}

func NewStudyService() *StudyService {
	return &StudyService{db: database.DB}
}

func (s *StudyService) GetStudySessions(page, perPage int) (*models.PaginatedResponse, error) {
	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM study_sessions").Scan(&total)
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * perPage
	rows, err := s.db.Query(`
		SELECT 
			ss.id,
			sa.name as activity_name,
			g.name as group_name,
			ss.created_at as start_time,
			ss.created_at as end_time,
			COUNT(wri.word_id) as review_items_count
		FROM study_sessions ss
		JOIN groups g ON ss.group_id = g.id
		JOIN study_activities sa ON ss.study_activity_id = sa.id
		LEFT JOIN word_review_items wri ON ss.id = wri.study_session_id
		GROUP BY ss.id
		ORDER BY ss.created_at DESC
		LIMIT ? OFFSET ?
	`, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.StudySession
	for rows.Next() {
		var s models.StudySession
		if err := rows.Scan(
			&s.ID,
			&s.ActivityName,
			&s.GroupName,
			&s.CreatedAt,
			&s.CreatedAt,
			&s.ReviewItemCount,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}

	return &models.PaginatedResponse{
		Items: sessions,
		Pagination: models.Pagination{
			CurrentPage:  page,
			TotalPages:   (total + perPage - 1) / perPage,
			TotalItems:   total,
			ItemsPerPage: perPage,
		},
	}, nil
}

func (s *StudyService) CreateStudyActivity(groupID, studyActivityID int) (*models.StudySession, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.Exec(`
		INSERT INTO study_sessions (group_id, study_activity_id, created_at)
		VALUES (?, ?, ?)
	`, groupID, studyActivityID, time.Now())
	if err != nil {
		return nil, err
	}

	sessionID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &models.StudySession{
		ID:              int(sessionID),
		GroupID:         groupID,
		StudyActivityID: studyActivityID,
	}, nil
}

func (s *StudyService) GetStudyActivity(id int) (*models.StudyActivity, error) {
	var activity models.StudyActivity
	err := s.db.QueryRow(`
		SELECT id, name, thumbnail_url, description
		FROM study_activities
		WHERE id = ?
	`, id).Scan(&activity.ID, &activity.Name, &activity.ThumbnailURL, &activity.Description)
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

func (s *StudyService) ReviewWord(sessionID, wordID int, correct bool) error {
	_, err := s.db.Exec(`
		INSERT INTO word_review_items (word_id, study_session_id, correct, created_at)
		VALUES (?, ?, ?, ?)
	`, wordID, sessionID, correct, time.Now())
	return err
} 