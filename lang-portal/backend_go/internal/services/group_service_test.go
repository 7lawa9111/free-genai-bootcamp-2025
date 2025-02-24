package services

import (
	"testing"

	"lang-portal/backend_go/internal/repository"
	"lang-portal/backend_go/internal/testutil"
)

func TestGroupService_GetGroups(t *testing.T) {
	// Setup test database
	db := testutil.SetupTestDB(t)
	
	// Create test data
	_, err := db.Exec(`
		INSERT INTO groups (name) VALUES
		('Basic Greetings'),
		('Numbers')
	`)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Create repositories and service
	groupRepo := repository.NewGroupRepository(db)
	wordRepo := repository.NewWordRepository(db)
	sessionRepo := repository.NewStudySessionRepository(db)
	service := NewGroupService(groupRepo, wordRepo, sessionRepo)

	// Test GetGroups
	groups, total, err := service.GetGroups(1, 10)
	if err != nil {
		t.Errorf("GetGroups failed: %v", err)
	}

	if len(groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(groups))
	}

	if total != 2 {
		t.Errorf("expected total of 2, got %d", total)
	}

	// Check first group
	if groups[0].Name != "Basic Greetings" {
		t.Errorf("expected Basic Greetings, got %s", groups[0].Name)
	}
}

func TestGroupService_GetGroupWords(t *testing.T) {
	// Setup test database
	db := testutil.SetupTestDB(t)
	
	// Create test data
	_, err := db.Exec(`
		INSERT INTO groups (id, name) VALUES (1, 'Basic Greetings');
		INSERT INTO words (id, japanese, romaji, english) VALUES 
		(1, 'こんにちは', 'konnichiwa', 'hello');
		INSERT INTO words_groups (word_id, group_id) VALUES (1, 1);
	`)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Create repositories and service
	groupRepo := repository.NewGroupRepository(db)
	wordRepo := repository.NewWordRepository(db)
	sessionRepo := repository.NewStudySessionRepository(db)
	service := NewGroupService(groupRepo, wordRepo, sessionRepo)

	// Test GetGroupWords
	words, total, err := service.GetGroupWords(1, 1, 10)
	if err != nil {
		t.Errorf("GetGroupWords failed: %v", err)
	}

	if len(words) != 1 {
		t.Errorf("expected 1 word, got %d", len(words))
	}

	if total != 1 {
		t.Errorf("expected total of 1, got %d", total)
	}

	// Check word details
	if words[0].Japanese != "こんにちは" {
		t.Errorf("expected こんにちは, got %s", words[0].Japanese)
	}
} 