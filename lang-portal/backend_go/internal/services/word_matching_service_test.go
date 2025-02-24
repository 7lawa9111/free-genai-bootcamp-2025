package services

import (
	"testing"
	"time"

	"lang-portal/backend_go/internal/models"
	"lang-portal/backend_go/internal/repository"
	"lang-portal/backend_go/internal/testutil"
)

func setupWordMatchingService(t *testing.T) (*WordMatchingService, *sql.DB) {
	db := testutil.SetupTestDB(t)
	repo := repository.NewWordMatchingRepository(db)
	service := NewWordMatchingService(repo)
	return service, db
}

func TestCreateActivity(t *testing.T) {
	service, db := setupWordMatchingService(t)

	// Create test data
	_, err := db.Exec(`
		INSERT INTO groups (id, name) VALUES (1, 'Test Group');
		INSERT INTO words (id, japanese, english) VALUES 
		(1, 'こんにちは', 'hello'),
		(2, 'さようなら', 'goodbye'),
		(3, 'ありがとう', 'thank you'),
		(4, 'おはよう', 'good morning');
		INSERT INTO words_groups (word_id, group_id) VALUES (1, 1), (2, 1), (3, 1), (4, 1);`)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Test creating activity
	activity, err := service.CreateActivity(1)
	if err != nil {
		t.Errorf("CreateActivity failed: %v", err)
	}

	if activity.GroupID != 1 {
		t.Errorf("expected group ID 1; got %d", activity.GroupID)
	}

	if len(activity.WordPairs) < 4 {
		t.Errorf("expected at least 4 word pairs; got %d", len(activity.WordPairs))
	}

	// Test creating activity with invalid group
	_, err = service.CreateActivity(999)
	if err == nil {
		t.Error("expected error for invalid group ID")
	}
}

func TestSaveResult(t *testing.T) {
	service, db := setupWordMatchingService(t)

	// Create test data
	_, err := db.Exec(`
		INSERT INTO study_activities (id, group_id, type, created_at) 
		VALUES (1, 1, 'word_matching', ?);`,
		time.Now())
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Test valid result
	result := &models.WordMatchingResult{
		ActivityID:   1,
		Correct:      true,
		TimeTaken:    60,
		MatchedPairs: 5,
		TotalPairs:   10,
	}
	err = service.SaveResult(result)
	if err != nil {
		t.Errorf("SaveResult failed: %v", err)
	}

	// Test invalid result data
	invalidResults := []models.WordMatchingResult{
		{ActivityID: 1, TimeTaken: -1, MatchedPairs: 5, TotalPairs: 10},    // Invalid time
		{ActivityID: 1, TimeTaken: 60, MatchedPairs: -1, TotalPairs: 10},   // Invalid matched pairs
		{ActivityID: 1, TimeTaken: 60, MatchedPairs: 11, TotalPairs: 10},   // More matched than total
		{ActivityID: 999, TimeTaken: 60, MatchedPairs: 5, TotalPairs: 10},  // Invalid activity ID
	}

	for _, invalid := range invalidResults {
		if err := service.SaveResult(&invalid); err == nil {
			t.Errorf("expected error for invalid result: %+v", invalid)
		}
	}
}

func TestGetProgress(t *testing.T) {
	service, db := setupWordMatchingService(t)

	// Create test data
	_, err := db.Exec(`
		INSERT INTO study_activities (id, group_id, type, created_at) 
		VALUES (1, 1, 'word_matching', ?);
		INSERT INTO study_sessions (id, group_id, study_activity_id, created_at)
		VALUES (1, 1, 1, ?);
		INSERT INTO word_review_items (word_id, study_session_id, correct, created_at)
		VALUES 
		(1, 1, 1, ?),
		(2, 1, 1, ?),
		(3, 1, 0, ?),
		(4, 1, 1, ?);`,
		time.Now(), time.Now(), time.Now(), time.Now(), time.Now(), time.Now())
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Test progress calculation
	progress, err := service.GetProgress(1)
	if err != nil {
		t.Errorf("GetProgress failed: %v", err)
	}

	expectedProgress := 0.75 // 3 correct out of 4 total
	if progress != expectedProgress {
		t.Errorf("expected progress %.2f; got %.2f", expectedProgress, progress)
	}

	// Test completion status
	isComplete, err := service.IsActivityComplete(1)
	if err != nil {
		t.Errorf("IsActivityComplete failed: %v", err)
	}

	if !isComplete {
		t.Error("expected activity to be complete")
	}
} 