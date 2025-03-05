package services

import (
	"database/sql"
	"github.com/mohawa/lang-portal/backend_go/internal/database"
)

type GroupService struct {
	db *sql.DB
}

func NewGroupService() *GroupService {
	return &GroupService{
		db: database.DB,
	}
}

func (s *GroupService) GetGroup(id int) (map[string]interface{}, error) {
	var name string
	err := s.db.QueryRow("SELECT name FROM groups WHERE id = ?", id).Scan(&name)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":   id,
		"name": name,
		"stats": map[string]interface{}{
			"total_word_count": 5,
		},
	}, nil
} 