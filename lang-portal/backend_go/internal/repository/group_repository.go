package repository

import (
	"database/sql"
	"lang-portal/backend_go/internal/models"
	"lang-portal/backend_go/internal/errors"
	"fmt"
)

type GroupRepository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) GetAll(page, limit int) ([]models.Group, error) {
	offset := (page - 1) * limit
	rows, err := r.db.Query(`
		SELECT id, name 
		FROM groups 
		LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var group models.Group
		err := rows.Scan(&group.ID, &group.Name)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (r *GroupRepository) GetByID(id int64) (*models.Group, error) {
	group := &models.Group{}
	err := r.db.QueryRow(`
		SELECT id, name 
		FROM groups 
		WHERE id = ?`, id).Scan(&group.ID, &group.Name)
	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (r *GroupRepository) GetTotalCount() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM groups").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *GroupRepository) GetWordCount(groupID int64) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) 
		FROM words_groups 
		WHERE group_id = ?`, groupID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *GroupRepository) GetGroupWords(groupID int64, page, limit int) ([]models.Word, error) {
	offset := (page - 1) * limit
	query := `
		SELECT w.id, w.japanese, w.english, w.romaji, w.parts
		FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		WHERE wg.group_id = ?
		LIMIT ? OFFSET ?`

	fmt.Printf("Debug - Query: %s\n", query)
	fmt.Printf("Debug - Params: groupID=%d, limit=%d, offset=%d\n", groupID, limit, offset)

	rows, err := r.db.Query(query, groupID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	var words []models.Word
	for rows.Next() {
		var word models.Word
		var parts sql.NullString
		err := rows.Scan(&word.ID, &word.Japanese, &word.English, &word.Romaji, &parts)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %v", err)
		}
		if parts.Valid {
			word.Parts = &parts.String
		}
		words = append(words, word)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %v", err)
	}

	fmt.Printf("Debug - Found %d words\n", len(words))
	return words, nil
}

func (r *GroupRepository) GetGroupStudySessions(groupID int64, page, limit int) ([]models.StudySession, error) {
	offset := (page - 1) * limit
	rows, err := r.db.Query(`
		SELECT id, group_id, study_activity_id, created_at
		FROM study_sessions
		WHERE group_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`,
		groupID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.StudySession
	for rows.Next() {
		var session models.StudySession
		err := rows.Scan(&session.ID, &session.GroupID, &session.StudyActivityID, &session.CreatedAt)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (r *GroupRepository) GetStudySessionCount(groupID int64) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) 
		FROM study_sessions 
		WHERE group_id = ?`, groupID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *GroupRepository) GetGroupWordCount(groupID int64) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM words_groups
		WHERE group_id = ?`, groupID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count error: %v", err)
	}
	return count, nil
} 