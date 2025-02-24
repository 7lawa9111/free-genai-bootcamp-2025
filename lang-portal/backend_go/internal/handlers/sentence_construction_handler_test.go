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

func setupSentenceTestRouter(t *testing.T) (*gin.Engine, *sql.DB) {
	db := testutil.SetupTestDB(t)
	
	repo := repository.NewSentenceConstructionRepository(db)
	service := services.NewSentenceConstructionService(repo)
	handler := NewSentenceConstructionHandler(service)

	r := gin.Default()
	api := r.Group("/api/sentence-construction")
	{
		api.POST("/activities", handler.CreateActivity)
		api.POST("/activities/:id/result", handler.SaveResult)
		api.GET("/activities/:id/stats", handler.GetActivityStats)
		api.GET("/activities/:id/progress", handler.GetProgress)
	}

	return r, db
}

func TestCreateSentenceActivity(t *testing.T) {
	r, db := setupSentenceTestRouter(t)

	// Create test data
	_, err := db.Exec(`
		INSERT INTO groups (id, name) VALUES (1, 'Test Group');
		INSERT INTO sentences (id, japanese, english, words, hints) VALUES 
		(1, '私は学生です', 'I am a student', '["私は", "学生", "です"]', '["Subject marker", "Noun", "Copula"]'),
		(2, '本を読みます', 'I read a book', '["本を", "読みます"]', '["Object marker", "Verb"]'),
		(3, '日本語を勉強します', 'I study Japanese', '["日本語を", "勉強", "します"]', '["Object", "Noun", "Verb"]');
		INSERT INTO sentences_groups (sentence_id, group_id) VALUES 
		(1, 1), (2, 1), (3, 1);`)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Make request
	reqBody := CreateSentenceRequest{GroupID: 1}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/sentence-construction/activities", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusCreated {
		t.Errorf("expected status Created; got %v", w.Code)
	}

	var response models.SentenceConstruction
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response.GroupID != 1 {
		t.Errorf("expected group ID 1; got %d", response.GroupID)
	}

	if len(response.Sentences) < 3 {
		t.Errorf("expected at least 3 sentences; got %d", len(response.Sentences))
	}

	// Verify each sentence has required fields
	for i, s := range response.Sentences {
		if s.Japanese == "" || s.English == "" {
			t.Errorf("sentence %d: missing required fields", i+1)
		}
		if len(s.Words) < 2 {
			t.Errorf("sentence %d: expected at least 2 words; got %d", i+1, len(s.Words))
		}
	}
}

func TestSaveSentenceResult(t *testing.T) {
	r, db := setupSentenceTestRouter(t)

	// Create test data
	_, err := db.Exec(`
		INSERT INTO study_activities (id, group_id, type, created_at) 
		VALUES (1, 1, 'sentence_construction', ?);`,
		time.Now())
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Make request
	result := models.SentenceResult{
		ActivityID:      1,
		CompletedCount: 4,
		TotalSentences: 5,
		TimeTaken:     180,
		Accuracy:      80.0,
	}
	body, _ := json.Marshal(result)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/sentence-construction/activities/1/result", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", w.Code)
	}

	// Verify result was saved
	var accuracy float64
	err = db.QueryRow("SELECT accuracy_score FROM study_activities WHERE id = 1").Scan(&accuracy)
	if err != nil {
		t.Fatalf("failed to verify result: %v", err)
	}
	if accuracy != 80.0 {
		t.Errorf("expected accuracy 80.0; got %f", accuracy)
	}
}

func TestGetSentenceStats(t *testing.T) {
	r, db := setupSentenceTestRouter(t)

	// Create test data
	_, err := db.Exec(`
		INSERT INTO study_activities (id, group_id, type, created_at, accuracy_score) 
		VALUES (1, 1, 'sentence_construction', ?, 80.0);
		INSERT INTO study_sessions (id, group_id, study_activity_id, created_at)
		VALUES (1, 1, 1, ?);
		INSERT INTO sentence_review_items (sentence_id, study_session_id, correct, created_at)
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
	req, _ := http.NewRequest("GET", "/api/sentence-construction/activities/1/stats", nil)
	r.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", w.Code)
	}

	var stats models.SentenceResult
	err = json.Unmarshal(w.Body.Bytes(), &stats)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if stats.CompletedCount != 3 {
		t.Errorf("expected 3 completed sentences; got %d", stats.CompletedCount)
	}

	if stats.Accuracy != 80.0 {
		t.Errorf("expected accuracy 80.0; got %f", stats.Accuracy)
	}
}

func TestGetSentenceProgress(t *testing.T) {
	r, db := setupSentenceTestRouter(t)

	// Create test data
	_, err := db.Exec(`
		INSERT INTO study_activities (id, group_id, type, created_at) 
		VALUES (1, 1, 'sentence_construction', ?);
		INSERT INTO study_sessions (id, group_id, study_activity_id, created_at)
		VALUES (1, 1, 1, ?);
		INSERT INTO sentence_review_items (sentence_id, study_session_id, correct, created_at)
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
	req, _ := http.NewRequest("GET", "/api/sentence-construction/activities/1/progress", nil)
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

	expectedProgress := 0.8 // 4 completed out of 5 total
	if response.Progress != expectedProgress {
		t.Errorf("expected progress %.2f; got %.2f", expectedProgress, response.Progress)
	}

	if !response.IsComplete {
		t.Error("expected activity to be complete")
	}
} 