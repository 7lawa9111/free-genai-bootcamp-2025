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

func setupFlashcardTestRouter(t *testing.T) (*gin.Engine, *sql.DB) {
	db := testutil.SetupTestDB(t)
	
	repo := repository.NewFlashcardRepository(db)
	service := services.NewFlashcardService(repo)
	handler := NewFlashcardHandler(service)

	r := gin.Default()
	api := r.Group("/api/flashcards")
	{
		api.POST("/activities", handler.CreateActivity)
		api.POST("/activities/:id/result", handler.SaveResult)
		api.GET("/activities/:id/stats", handler.GetActivityStats)
		api.GET("/activities/:id/progress", handler.GetProgress)
	}

	return r, db
}

func TestCreateFlashcardActivity(t *testing.T) {
	r, db := setupFlashcardTestRouter(t)

	// Create test data
	_, err := db.Exec(`
		INSERT INTO groups (id, name) VALUES (1, 'Test Group');
		INSERT INTO words (id, japanese, english, romaji) VALUES 
		(1, 'こんにちは', 'hello', 'konnichiwa'),
		(2, 'さようなら', 'goodbye', 'sayounara'),
		(3, 'ありがとう', 'thank you', 'arigatou'),
		(4, 'おはよう', 'good morning', 'ohayou'),
		(5, 'こんばんは', 'good evening', 'konbanwa');
		INSERT INTO words_groups (word_id, group_id) VALUES 
		(1, 1), (2, 1), (3, 1), (4, 1), (5, 1);`)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Make request
	reqBody := CreateFlashcardRequest{
		GroupID:   1,
		Direction: "ja_to_en",
	}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/flashcards/activities", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusCreated {
		t.Errorf("expected status Created; got %v", w.Code)
	}

	var response models.FlashcardActivity
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response.GroupID != 1 {
		t.Errorf("expected group ID 1; got %d", response.GroupID)
	}

	if response.Direction != "ja_to_en" {
		t.Errorf("expected direction ja_to_en; got %s", response.Direction)
	}

	if len(response.Cards) < 5 {
		t.Errorf("expected at least 5 cards; got %d", len(response.Cards))
	}
}

func TestSaveFlashcardResult(t *testing.T) {
	r, db := setupFlashcardTestRouter(t)

	// Create test data
	_, err := db.Exec(`
		INSERT INTO study_activities (id, group_id, type, created_at) 
		VALUES (1, 1, 'flashcards', ?);`,
		time.Now())
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Make request
	result := models.FlashcardResult{
		ActivityID:    1,
		CardsReviewed: 8,
		TotalCards:    10,
		TimeTaken:     120,
		Confidence:    0.8,
	}
	body, _ := json.Marshal(result)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/flashcards/activities/1/result", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", w.Code)
	}

	// Verify result was saved
	var confidence float64
	err = db.QueryRow("SELECT confidence_score FROM study_activities WHERE id = 1").Scan(&confidence)
	if err != nil {
		t.Fatalf("failed to verify result: %v", err)
	}
	if confidence != 0.8 {
		t.Errorf("expected confidence 0.8; got %f", confidence)
	}
}

func TestGetFlashcardStats(t *testing.T) {
	r, db := setupFlashcardTestRouter(t)

	// Create test data
	_, err := db.Exec(`
		INSERT INTO study_activities (id, group_id, type, created_at, confidence_score) 
		VALUES (1, 1, 'flashcards', ?, 0.8);
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
	req, _ := http.NewRequest("GET", "/api/flashcards/activities/1/stats", nil)
	r.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", w.Code)
	}

	var stats models.FlashcardResult
	err = json.Unmarshal(w.Body.Bytes(), &stats)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if stats.CardsReviewed != 3 {
		t.Errorf("expected 3 cards reviewed; got %d", stats.CardsReviewed)
	}

	if stats.Confidence != 0.8 {
		t.Errorf("expected confidence 0.8; got %f", stats.Confidence)
	}
}

func TestGetFlashcardProgress(t *testing.T) {
	r, db := setupFlashcardTestRouter(t)

	// Create test data
	_, err := db.Exec(`
		INSERT INTO study_activities (id, group_id, type, created_at) 
		VALUES (1, 1, 'flashcards', ?);
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
	req, _ := http.NewRequest("GET", "/api/flashcards/activities/1/progress", nil)
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