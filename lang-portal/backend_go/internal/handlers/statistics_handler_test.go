package handlers

import (
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

func setupStatisticsTestRouter(t *testing.T) (*gin.Engine, *sql.DB) {
	db := testutil.SetupTestDB(t)
	
	repo := repository.NewStatisticsRepository(db)
	service := services.NewStatisticsService(repo)
	handler := NewStatisticsHandler(service)

	r := gin.Default()
	api := r.Group("/api/statistics")
	{
		api.GET("/user", handler.GetUserStats)
		api.GET("/activities", handler.GetActivityStats)
		api.GET("/progress", handler.GetStudyProgress)
		api.GET("/metrics", handler.GetStudyMetrics)
	}

	return r, db
}

func TestGetUserStats(t *testing.T) {
	r, db := setupStatisticsTestRouter(t)

	// Create test data
	_, err := db.Exec(`
		INSERT INTO study_activities (id, group_id, type, created_at, completed_at, accuracy_score) VALUES 
		(1, 1, 'word_matching', ?, ?, 85.0),
		(2, 1, 'vocabulary_quiz', ?, ?, 90.0),
		(3, 1, 'flashcards', ?, NULL, NULL);`,
		time.Now().Add(-48*time.Hour), time.Now().Add(-47*time.Hour),
		time.Now().Add(-24*time.Hour), time.Now().Add(-23*time.Hour),
		time.Now(), nil)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Make request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/statistics/user", nil)
	r.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", w.Code)
	}

	var stats models.UserStats
	err = json.Unmarshal(w.Body.Bytes(), &stats)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if stats.TotalActivities != 3 {
		t.Errorf("expected 3 total activities; got %d", stats.TotalActivities)
	}

	if stats.CompletedActivities != 2 {
		t.Errorf("expected 2 completed activities; got %d", stats.CompletedActivities)
	}

	expectedAccuracy := 87.5 // Average of 85 and 90
	if stats.AverageAccuracy != expectedAccuracy {
		t.Errorf("expected accuracy %.2f; got %.2f", expectedAccuracy, stats.AverageAccuracy)
	}
}

func TestGetActivityStats(t *testing.T) {
	r, db := setupStatisticsTestRouter(t)

	// Create test data
	_, err := db.Exec(`
		INSERT INTO study_activities (id, group_id, type, created_at, completed_at, accuracy_score) VALUES 
		(1, 1, 'word_matching', ?, ?, 85.0),
		(2, 1, 'word_matching', ?, ?, 95.0),
		(3, 1, 'vocabulary_quiz', ?, ?, 90.0),
		(4, 1, 'vocabulary_quiz', ?, NULL, NULL);`,
		time.Now().Add(-72*time.Hour), time.Now().Add(-71*time.Hour),
		time.Now().Add(-48*time.Hour), time.Now().Add(-47*time.Hour),
		time.Now().Add(-24*time.Hour), time.Now().Add(-23*time.Hour),
		time.Now(), nil)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Make request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/statistics/activities", nil)
	r.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", w.Code)
	}

	var stats map[string]models.ActivityStats
	err = json.Unmarshal(w.Body.Bytes(), &stats)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	wordMatching := stats["word_matching"]
	if wordMatching.CompletionRate != 100.0 {
		t.Errorf("expected word matching completion rate 100.0; got %.2f", wordMatching.CompletionRate)
	}
	if wordMatching.AverageAccuracy != 90.0 {
		t.Errorf("expected word matching accuracy 90.0; got %.2f", wordMatching.AverageAccuracy)
	}

	vocabQuiz := stats["vocabulary_quiz"]
	if vocabQuiz.CompletionRate != 50.0 {
		t.Errorf("expected vocabulary quiz completion rate 50.0; got %.2f", vocabQuiz.CompletionRate)
	}
}

func TestGetStudyProgress(t *testing.T) {
	r, db := setupStatisticsTestRouter(t)

	// Create test data for the last 7 days
	now := time.Now()
	for i := 0; i < 7; i++ {
		date := now.AddDate(0, 0, -i)
		if i < 3 { // Add activities for last 3 days to create a streak
			_, err := db.Exec(`
				INSERT INTO study_activities (id, group_id, type, created_at, completed_at) 
				VALUES (?, 1, 'word_matching', ?, ?);`,
				i+1, date, date.Add(time.Hour))
			if err != nil {
				t.Fatalf("failed to insert test data: %v", err)
			}
		}
	}

	// Make request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/statistics/progress", nil)
	r.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", w.Code)
	}

	var progress models.StudyProgress
	err := json.Unmarshal(w.Body.Bytes(), &progress)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if progress.DailyStreak != 3 {
		t.Errorf("expected daily streak 3; got %d", progress.DailyStreak)
	}

	if len(progress.WeeklyProgress) != 7 {
		t.Errorf("expected 7 days of progress; got %d", len(progress.WeeklyProgress))
	}
}

func TestGetStudyMetrics(t *testing.T) {
	r, db := setupStatisticsTestRouter(t)

	// Create test data
	_, err := db.Exec(`
		INSERT INTO study_activities (id, group_id, type, created_at, completed_at) VALUES 
		(1, 1, 'word_matching', ?, ?),
		(2, 1, 'vocabulary_quiz', ?, ?),
		(3, 1, 'flashcards', ?, NULL);
		
		INSERT INTO study_sessions (id, study_activity_id, time_taken_seconds, created_at) VALUES
		(1, 1, 600, ?),  -- 10 minutes
		(2, 2, 900, ?);  -- 15 minutes`,
		time.Now().Add(-2*time.Hour), time.Now().Add(-1*time.Hour),
		time.Now().Add(-4*time.Hour), time.Now().Add(-3*time.Hour),
		time.Now(), 
		time.Now().Add(-2*time.Hour),
		time.Now().Add(-4*time.Hour))
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Make request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/statistics/metrics", nil)
	r.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", w.Code)
	}

	var metrics StudyMetricsResponse
	err = json.Unmarshal(w.Body.Bytes(), &metrics)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	expectedAvgTime := 12.5 // (10 + 15) / 2 minutes
	if metrics.AverageStudyTimeMinutes != expectedAvgTime {
		t.Errorf("expected average time %.2f; got %.2f", expectedAvgTime, metrics.AverageStudyTimeMinutes)
	}

	expectedCompletionRate := 66.67 // 2 completed out of 3 total
	if metrics.CompletionRatePercent < 66.0 || metrics.CompletionRatePercent > 67.0 {
		t.Errorf("expected completion rate around %.2f; got %.2f", expectedCompletionRate, metrics.CompletionRatePercent)
	}
} 