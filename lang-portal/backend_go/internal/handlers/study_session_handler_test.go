package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"lang-portal/backend_go/internal/models"
	"lang-portal/backend_go/internal/repository"
	"lang-portal/backend_go/internal/services"
	"lang-portal/backend_go/internal/testutil"
)

func setupTestRouter(t *testing.T) (*gin.Engine, *sql.DB) {
	db := testutil.SetupTestDB(t)
	
	repo := repository.NewStudySessionRepository(db)
	service := services.NewStudySessionService(repo)
	handler := NewStudySessionHandler(service)

	r := gin.Default()
	api := r.Group("/api")
	
	api.GET("/study_sessions", handler.GetStudySessions)
	api.GET("/study_sessions/:id", handler.GetStudySessionByID)
	api.GET("/study_sessions/:id/words", handler.GetStudySessionWords)
	api.POST("/study_sessions/:id/words/:word_id/review", handler.ReviewWord)

	return r, db
}

func TestGetStudySessions(t *testing.T) {
	r, db := setupTestRouter(t)

	// Create test data
	_, err := db.Exec(`
		INSERT INTO study_sessions (id, group_id, study_activity_id, created_at) VALUES
		(1, 1, 1, datetime('now')),
		(2, 1, 1, datetime('now', '-1 day'))`)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Make request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/study_sessions?page=1&limit=10", nil)
	r.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", w.Code)
	}

	var response ListResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	sessions := response.Items.([]interface{})
	if len(sessions) != 2 {
		t.Errorf("expected 2 sessions; got %d", len(sessions))
	}
}

func TestGetStudySessionByID(t *testing.T) {
	r, db := setupTestRouter(t)

	// Create test data
	_, err := db.Exec(`
		INSERT INTO study_sessions (id, group_id, study_activity_id, created_at)
		VALUES (1, 1, 1, datetime('now'))`)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Make request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/study_sessions/1", nil)
	r.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", w.Code)
	}

	var session models.StudySessionDetails
	err = json.Unmarshal(w.Body.Bytes(), &session)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if session.ID != 1 {
		t.Errorf("expected session ID 1; got %d", session.ID)
	}

	// Test non-existent session
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/study_sessions/999", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status NotFound; got %v", w.Code)
	}
} 