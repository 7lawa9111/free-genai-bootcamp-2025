package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"lang-portal/backend_go/internal/models"
	"lang-portal/backend_go/internal/repository"
	"lang-portal/backend_go/internal/services"
	"lang-portal/backend_go/internal/testutil"
)

func setupWordMatchingTestRouter(t *testing.T) (*gin.Engine, *sql.DB) {
	db := testutil.SetupTestDB(t)
	
	repo := repository.NewWordMatchingRepository(db)
	service := services.NewWordMatchingService(repo)
	handler := NewWordMatchingHandler(service)

	r := gin.Default()
	api := r.Group("/api/word-matching")
	{
		api.POST("/activities", handler.CreateActivity)
		api.POST("/activities/:id/result", handler.SaveResult)
		api.GET("/activities/:id/stats", handler.GetActivityStats)
		api.GET("/activities/:id/progress", handler.GetProgress)
	}

	return r, db
}

func TestCreateActivity(t *testing.T) {
	r, db := setupWordMatchingTestRouter(t)

	// Create test data
	_, err := db.Exec(`
		INSERT INTO groups (id, name) VALUES (1, 'Test Group');
		INSERT INTO words (id, japanese, english) VALUES 
		(1, 'こんにちは', 'hello'),
		(2, 'さようなら', 'goodbye'),
		(3, 'ありがとう', 'thank you');
		INSERT INTO words_groups (word_id, group_id) VALUES (1, 1), (2, 1), (3, 1);`)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Make request
	reqBody := CreateActivityRequest{GroupID: 1}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/word-matching/activities", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusCreated {
		t.Errorf("expected status Created; got %v", w.Code)
	}

	var response models.WordMatchingActivity
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response.GroupID != 1 {
		t.Errorf("expected group ID 1; got %d", response.GroupID)
	}

	if len(response.WordPairs) == 0 {
		t.Error("expected word pairs; got none")
	}
}

func TestSaveResult(t *testing.T) {
	r, db := setupWordMatchingTestRouter(t)

	// Create test data
	_, err := db.Exec(`
		INSERT INTO study_activities (id, group_id, type, created_at) 
		VALUES (1, 1, 'word_matching', ?);`,
		time.Now())
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Make request
	result := models.WordMatchingResult{
		ActivityID:   1,
		Correct:      true,
		TimeTaken:    60,
		MatchedPairs: 5,
		TotalPairs:   10,
	}
	body, _ := json.Marshal(result)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/word-matching/activities/1/result", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", w.Code)
	}

	// Verify result was saved
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM word_review_items WHERE study_session_id IN (SELECT id FROM study_sessions WHERE study_activity_id = 1)").Scan(&count)
	if err != nil {
		t.Fatalf("failed to verify result: %v", err)
	}
	if count == 0 {
		t.Error("expected word review items; got none")
	}
}

func TestGetActivityStats(t *testing.T) {
	r, db := setupWordMatchingTestRouter(t)

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
		(3, 1, 0, ?);`,
		time.Now(), time.Now(), time.Now(), time.Now(), time.Now())
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Make request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/word-matching/activities/1/stats", nil)
	r.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", w.Code)
	}

	var stats models.WordMatchingResult
	err = json.Unmarshal(w.Body.Bytes(), &stats)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if stats.ActivityID != 1 {
		t.Errorf("expected activity ID 1; got %d", stats.ActivityID)
	}

	if stats.MatchedPairs != 2 {
		t.Errorf("expected 2 matched pairs; got %d", stats.MatchedPairs)
	}

	if stats.TotalPairs != 3 {
		t.Errorf("expected 3 total pairs; got %d", stats.TotalPairs)
	}
}

func TestGetProgress(t *testing.T) {
	r, db := setupWordMatchingTestRouter(t)

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

	// Make request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/word-matching/activities/1/progress", nil)
	r.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", w.Code)
	}

	var response ProgressResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	expectedProgress := 0.75 // 3 correct out of 4 total
	if response.Progress != expectedProgress {
		t.Errorf("expected progress %.2f; got %.2f", expectedProgress, response.Progress)
	}

	if !response.IsComplete {
		t.Error("expected activity to be complete")
	}
}

func TestCreateActivity_InvalidGroup(t *testing.T) {
	r, db := setupWordMatchingTestRouter(t)

	// Make request with non-existent group
	reqBody := CreateActivityRequest{GroupID: 999}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/word-matching/activities", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest; got %v", w.Code)
	}
}

func TestSaveResult_InvalidData(t *testing.T) {
	r, _ := setupWordMatchingTestRouter(t)

	// Make request with invalid data
	result := models.WordMatchingResult{
		ActivityID:   1,
		TimeTaken:    -1, // Invalid time
		MatchedPairs: 5,
		TotalPairs:   3, // Invalid: more matched than total
	}
	body, _ := json.Marshal(result)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/word-matching/activities/1/result", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest; got %v", w.Code)
	}
} 