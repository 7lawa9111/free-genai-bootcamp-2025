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

func setupVocabularyQuizTestRouter(t *testing.T) (*gin.Engine, *sql.DB) {
	db := testutil.SetupTestDB(t)
	
	repo := repository.NewVocabularyQuizRepository(db)
	service := services.NewVocabularyQuizService(repo)
	handler := NewVocabularyQuizHandler(service)

	r := gin.Default()
	api := r.Group("/api/vocabulary-quiz")
	{
		api.POST("/quizzes", handler.CreateQuiz)
		api.POST("/quizzes/:id/result", handler.SaveResult)
		api.GET("/quizzes/:id/stats", handler.GetQuizStats)
		api.GET("/quizzes/:id/progress", handler.GetProgress)
	}

	return r, db
}

func TestCreateQuiz(t *testing.T) {
	r, db := setupVocabularyQuizTestRouter(t)

	// Create test data
	_, err := db.Exec(`
		INSERT INTO groups (id, name) VALUES (1, 'Test Group');
		INSERT INTO words (id, japanese, english) VALUES 
		(1, 'こんにちは', 'hello'),
		(2, 'さようなら', 'goodbye'),
		(3, 'ありがとう', 'thank you'),
		(4, 'おはよう', 'good morning'),
		(5, 'こんばんは', 'good evening');
		INSERT INTO words_groups (word_id, group_id) VALUES 
		(1, 1), (2, 1), (3, 1), (4, 1), (5, 1);`)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Make request
	reqBody := CreateQuizRequest{GroupID: 1}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/vocabulary-quiz/quizzes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusCreated {
		t.Errorf("expected status Created; got %v", w.Code)
	}

	var response models.VocabularyQuiz
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response.GroupID != 1 {
		t.Errorf("expected group ID 1; got %d", response.GroupID)
	}

	if len(response.Questions) < 5 {
		t.Errorf("expected at least 5 questions; got %d", len(response.Questions))
	}

	// Verify each question has 4 options
	for i, q := range response.Questions {
		if len(q.Options) != 4 {
			t.Errorf("question %d: expected 4 options; got %d", i+1, len(q.Options))
		}
	}
}

func TestSaveResult(t *testing.T) {
	r, db := setupVocabularyQuizTestRouter(t)

	// Create test data
	_, err := db.Exec(`
		INSERT INTO study_activities (id, group_id, type, created_at) 
		VALUES (1, 1, 'vocabulary_quiz', ?);`,
		time.Now())
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Make request
	result := models.QuizResult{
		ActivityID:    1,
		CorrectCount:  8,
		TotalCount:    10,
		TimeTaken:     120,
		Score:         80.0,
	}
	body, _ := json.Marshal(result)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/vocabulary-quiz/quizzes/1/result", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", w.Code)
	}

	// Verify result was saved
	var savedScore float64
	err = db.QueryRow("SELECT score FROM study_activities WHERE id = 1").Scan(&savedScore)
	if err != nil {
		t.Fatalf("failed to verify result: %v", err)
	}
	if savedScore != 80.0 {
		t.Errorf("expected score 80.0; got %f", savedScore)
	}
}

func TestGetQuizStats(t *testing.T) {
	r, db := setupVocabularyQuizTestRouter(t)

	// Create test data
	_, err := db.Exec(`
		INSERT INTO study_activities (id, group_id, type, created_at, score) 
		VALUES (1, 1, 'vocabulary_quiz', ?, 80.0);
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
	req, _ := http.NewRequest("GET", "/api/vocabulary-quiz/quizzes/1/stats", nil)
	r.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", w.Code)
	}

	var stats models.QuizResult
	err = json.Unmarshal(w.Body.Bytes(), &stats)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if stats.CorrectCount != 2 {
		t.Errorf("expected 2 correct answers; got %d", stats.CorrectCount)
	}

	if stats.Score != 80.0 {
		t.Errorf("expected score 80.0; got %f", stats.Score)
	}
}

func TestGetProgress(t *testing.T) {
	r, db := setupVocabularyQuizTestRouter(t)

	// Create test data
	_, err := db.Exec(`
		INSERT INTO study_activities (id, group_id, type, created_at) 
		VALUES (1, 1, 'vocabulary_quiz', ?);
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
	req, _ := http.NewRequest("GET", "/api/vocabulary-quiz/quizzes/1/progress", nil)
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
		t.Error("expected quiz to be complete")
	}
} 