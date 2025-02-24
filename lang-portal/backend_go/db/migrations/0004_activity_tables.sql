-- Add tables for all study activity types
CREATE TABLE IF NOT EXISTS word_matching_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    study_session_id INTEGER NOT NULL,
    word_id INTEGER NOT NULL,
    matched_word_id INTEGER NOT NULL,
    correct BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (study_session_id) REFERENCES study_sessions(id),
    FOREIGN KEY (word_id) REFERENCES words(id),
    FOREIGN KEY (matched_word_id) REFERENCES words(id)
);

CREATE TABLE IF NOT EXISTS vocabulary_quiz_answers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    study_session_id INTEGER NOT NULL,
    word_id INTEGER NOT NULL,
    selected_answer TEXT NOT NULL,
    correct BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (study_session_id) REFERENCES study_sessions(id),
    FOREIGN KEY (word_id) REFERENCES words(id)
);

CREATE TABLE IF NOT EXISTS flashcard_reviews (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    study_session_id INTEGER NOT NULL,
    word_id INTEGER NOT NULL,
    revealed BOOLEAN NOT NULL,
    confidence INTEGER CHECK(confidence BETWEEN 1 AND 5),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (study_session_id) REFERENCES study_sessions(id),
    FOREIGN KEY (word_id) REFERENCES words(id)
);

-- Add indices for better query performance
CREATE INDEX idx_word_matching_session ON word_matching_items(study_session_id);
CREATE INDEX idx_vocabulary_quiz_session ON vocabulary_quiz_answers(study_session_id);
CREATE INDEX idx_flashcard_reviews_session ON flashcard_reviews(study_session_id); 