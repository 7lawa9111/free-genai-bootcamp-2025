package services

import (
	"testing"

	"lang-portal/backend_go/internal/repository"
	"lang-portal/backend_go/internal/testutil"
)

func TestWordService_GetWords(t *testing.T) {
	// Setup test database
	db := testutil.SetupTestDB(t)
	
	// Create test data
	_, err := db.Exec(`
		INSERT INTO words (japanese, romaji, english) VALUES
		('こんにちは', 'konnichiwa', 'hello'),
		('さようなら', 'sayounara', 'goodbye')
	`)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Create repository and service
	repo := repository.NewWordRepository(db)
	service := NewWordService(repo)

	// Test GetWords
	words, total, err := service.GetWords(1, 10)
	if err != nil {
		t.Errorf("GetWords failed: %v", err)
	}

	if len(words) != 2 {
		t.Errorf("expected 2 words, got %d", len(words))
	}

	if total != 2 {
		t.Errorf("expected total of 2, got %d", total)
	}

	// Check first word
	if words[0].Japanese != "こんにちは" {
		t.Errorf("expected こんにちは, got %s", words[0].Japanese)
	}
}

func TestWordService_GetWordByID(t *testing.T) {
	// Setup test database
	db := testutil.SetupTestDB(t)
	
	// Create test data
	result, err := db.Exec(`
		INSERT INTO words (japanese, romaji, english) 
		VALUES ('こんにちは', 'konnichiwa', 'hello')
	`)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	wordID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("failed to get last insert id: %v", err)
	}

	// Create repository and service
	repo := repository.NewWordRepository(db)
	service := NewWordService(repo)

	// Test GetWordByID
	word, err := service.GetWordByID(wordID)
	if err != nil {
		t.Errorf("GetWordByID failed: %v", err)
	}

	if word.Japanese != "こんにちは" {
		t.Errorf("expected こんにちは, got %s", word.Japanese)
	}

	// Test non-existent word
	_, err = service.GetWordByID(999)
	if err == nil {
		t.Error("expected error for non-existent word")
	}
} 