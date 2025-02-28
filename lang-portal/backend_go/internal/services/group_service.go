package services

import (
	"database/sql"
	"github.com/mohawa/lang-portal/backend_go/internal/database"
	"github.com/mohawa/lang-portal/backend_go/internal/models"
)

type GroupService struct {
	db *sql.DB
}

func NewGroupService() *GroupService {
	return &GroupService{db: database.DB}
}

func (s *GroupService) GetGroups(page, perPage int) (*models.PaginatedResponse, error) {
	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM groups").Scan(&total)
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * perPage
	rows, err := s.db.Query(`
		SELECT g.id, g.name, COUNT(wg.word_id) as word_count
		FROM groups g
		LEFT JOIN words_groups wg ON g.id = wg.group_id
		GROUP BY g.id
		LIMIT ? OFFSET ?
	`, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var g models.Group
		if err := rows.Scan(&g.ID, &g.Name, &g.WordCount); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}

	return &models.PaginatedResponse{
		Items: groups,
		Pagination: models.Pagination{
			CurrentPage:  page,
			TotalPages:   (total + perPage - 1) / perPage,
			TotalItems:   total,
			ItemsPerPage: perPage,
		},
	}, nil
}

func (s *GroupService) GetGroup(id int) (*models.GroupResponse, error) {
	var group models.GroupResponse
	err := s.db.QueryRow(`
		SELECT g.id, g.name, COUNT(wg.word_id) as total_word_count
		FROM groups g
		LEFT JOIN words_groups wg ON g.id = wg.group_id
		WHERE g.id = ?
		GROUP BY g.id
	`, id).Scan(&group.ID, &group.Name, &group.Stats.TotalWordCount)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (s *GroupService) GetGroupWords(groupID, page, perPage int) (*models.PaginatedResponse, error) {
	var total int
	err := s.db.QueryRow(`
		SELECT COUNT(*) 
		FROM words_groups 
		WHERE group_id = ?
	`, groupID).Scan(&total)
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * perPage
	rows, err := s.db.Query(`
		SELECT w.japanese, w.romaji, w.english,
			   COUNT(CASE WHEN wri.correct = 1 THEN 1 END) as correct_count,
			   COUNT(CASE WHEN wri.correct = 0 THEN 1 END) as wrong_count
		FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		LEFT JOIN word_review_items wri ON w.id = wri.word_id
		WHERE wg.group_id = ?
		GROUP BY w.id
		LIMIT ? OFFSET ?
	`, groupID, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []models.WordWithStats
	for rows.Next() {
		var w models.WordWithStats
		if err := rows.Scan(&w.Japanese, &w.Romaji, &w.English, &w.CorrectCount, &w.WrongCount); err != nil {
			return nil, err
		}
		words = append(words, w)
	}

	return &models.PaginatedResponse{
		Items: words,
		Pagination: models.Pagination{
			CurrentPage:  page,
			TotalPages:   (total + perPage - 1) / perPage,
			TotalItems:   total,
			ItemsPerPage: perPage,
		},
	}, nil
}

func (s *GroupService) GetStudySessions(page, perPage int) (*models.PaginatedResponse, error) {
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