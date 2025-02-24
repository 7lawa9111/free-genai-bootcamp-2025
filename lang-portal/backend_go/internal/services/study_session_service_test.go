package services

import (
	"testing"
	"time"

	"lang-portal/backend_go/internal/models"
	"lang-portal/backend_go/internal/repository"
	"lang-portal/backend_go/internal/testutil"
)

func TestStudySessionService_GetStudySessions(t *testing.T) {
	// Setup test database
	db := testutil.SetupTestDB(t)

	// Create test data
	_, err := db.Exec(`
		INSERT INTO study_sessions (id, group_id, study_activity_id, created_at) VALUES
		(1, 1, 1, ?),
		(2, 1, 1, ?)`,
		time.Now(), time.Now().Add(-24*time.Hour))
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Create repository and service
	repo := repository.NewStudySessionRepository(db)
	service := NewStudySessionService(repo)

	// Test GetStudySessions
	sessions, total, err := service.GetStudySessions(1, 10)
	if err != nil {
		t.Errorf("GetStudySessions failed: %v", err)
	}

	if len(sessions) != 2 {
		t.Errorf("expected 2 sessions, got %d", len(sessions))
	}

	if total != 2 {
		t.Errorf("expected total of 2, got %d", total)
	}
}

func TestStudySessionService_GetStudySessionByID(t *testing.T) {
	// Setup test database
	db := testutil.SetupTestDB(t)

	// Create test data
	createdAt := time.Now()
	_, err := db.Exec(`
		INSERT INTO study_sessions (id, group_id, study_activity_id, created_at)
		VALUES (1, 1, 1, ?)`,
		createdAt)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Create repository and service
	repo := repository.NewStudySessionRepository(db)
	service := NewStudySessionService(repo)

	// Test GetStudySessionByID
	session, err := service.GetStudySessionByID(1)
	if err != nil {
		t.Errorf("GetStudySessionByID failed: %v", err)
	}

	if session.ID != 1 {
		t.Errorf("expected session ID 1, got %d", session.ID)
	}

	// Test non-existent session
	_, err = service.GetStudySessionByID(999)
	if err == nil {
		t.Error("expected error for non-existent session")
	}
}

func TestStudySessionService_ReviewWord(t *testing.T) {
	// Setup test database
	db := testutil.SetupTestDB(t)

	// Create test data
	_, err := db.Exec(`
		INSERT INTO study_sessions (id, group_id, study_activity_id, created_at) VALUES (1, 1, 1, ?);
		INSERT INTO words (id, japanese, romaji, english) VALUES (1, 'こんにちは', 'konnichiwa', 'hello');`,
		time.Now())
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Create repository and service
	repo := repository.NewStudySessionRepository(db)
	service := NewStudySessionService(repo)

	// Test ReviewWord
	err = service.ReviewWord(1, 1, true)
	if err != nil {
		t.Errorf("ReviewWord failed: %v", err)
	}

	// Verify review was recorded
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM word_review_items WHERE study_session_id = 1 AND word_id = 1").Scan(&count)
	if err != nil {
		t.Errorf("failed to verify review: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 review, got %d", count)
	}
} 