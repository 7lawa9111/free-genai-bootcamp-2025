package services

import (
	"database/sql"
	"github.com/mohawa/lang-portal/backend_go/internal/database"
)

type ResetService struct {
	db *sql.DB
}

func NewResetService() *ResetService {
	return &ResetService{db: database.DB}
}

func (s *ResetService) ResetHistory() error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clear study history tables
	_, err = tx.Exec("DELETE FROM word_review_items")
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM study_sessions")
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM study_activities")
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *ResetService) FullReset() error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clear all tables in correct order to respect foreign keys
	tables := []string{
		"word_review_items",
		"study_sessions",
		"study_activities",
		"words_groups",
		"words",
		"groups",
	}

	for _, table := range tables {
		_, err = tx.Exec("DELETE FROM " + table)
		if err != nil {
			return err
		}
	}

	// Reset auto-increment counters
	for _, table := range tables {
		_, err = tx.Exec("DELETE FROM sqlite_sequence WHERE name = ?", table)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
} 